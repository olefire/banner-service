package repository

import (
	"banner-service/internal/models"
	"context"
	"errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jellydator/ttlcache/v3"
	"log"
)

type BannerRepository struct {
	pool  *pgxpool.Pool
	cache *ttlcache.Cache[models.FeatureTag, models.BannerContent]
}

func NewBannerRepository(p *pgxpool.Pool, c *ttlcache.Cache[models.FeatureTag, models.BannerContent]) *BannerRepository {
	return &BannerRepository{
		pool:  p,
		cache: c,
	}
}

func RunInTx(ctx context.Context, pool *pgxpool.Pool, f func(tx pgx.Tx) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	if err = f(tx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return tx.Commit(ctx)
}

func (b *BannerRepository) GetBanner(ctx context.Context, tagId uint64, featureId uint64, isAdmin bool, useLastRevision bool) (string, error) {
	if !useLastRevision {
		if banner := b.cache.Get(models.FeatureTag{FeatureId: featureId, TagId: tagId}); banner != nil {
			if banner.Value().IsActive || isAdmin {
				log.Println("get banner from cache", banner.Value())
				return banner.Value().Content, nil
			} else {
				return "", ErrBannerInactive
			}
		}
	}

	const (
		selectBannerQuery = `select bv.content, b.is_active
           from banner_feature_tag bft
           join banner b using (banner_id)
           join banner_version bv on b.banner_id = bv.banner_id and b.active_version = bv.version
           where bft.feature_id = $1 and bft.tag_id = $2 and not must_be_deleted`
	)

	var bannerContent models.BannerContent
	if err := pgxscan.Get(ctx, b.pool, &bannerContent, selectBannerQuery, featureId, tagId); errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	} else if err != nil {
		return "", err
	} else if !(bannerContent.IsActive || isAdmin) {
		return "", ErrBannerInactive
	}

	_ = b.cache.Set(models.FeatureTag{FeatureId: featureId, TagId: tagId}, bannerContent, ttlcache.DefaultTTL)

	return bannerContent.Content, nil
}

func (b *BannerRepository) GetListOfVersions(ctx context.Context, bannerId uint64) ([]models.Banner, error) {
	const (
		selectBannersQuery = `SELECT b.banner_id, bft.feature_id, array_agg(DISTINCT bft.tag_id) AS tag_ids,
              bv.content, b.is_active, bv.version, b.created_at, bv.updated_at
         FROM banner_version bv
         JOIN banner b USING (banner_id)
         JOIN banner_feature_tag bft USING (banner_id)
         WHERE banner_id = $1 and b.must_be_deleted = false
         GROUP BY b.banner_id, bft.feature_id, bv.content, b.is_active, bv.version, b.created_at, bv.updated_at`
	)
	var banners []models.Banner

	if err := pgxscan.Select(ctx, b.pool, &banners, selectBannersQuery, bannerId); errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return banners, nil
}

func (b *BannerRepository) ChooseBannerVersion(ctx context.Context, bannerId uint64, version uint64) error {
	const (
		chooseVersionQuery = `update banner
						 set active_version = $2
						 where banner_id = $1`
	)

	res, err := b.pool.Exec(ctx, chooseVersionQuery, bannerId, version)
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

func (b *BannerRepository) GetFilteredBanners(ctx context.Context, filter *models.FilterBanner) ([]models.Banner, error) {
	const (
		selectFilteredBannersQuery = `
   	        select b.banner_id, bft.feature_id, array_agg(distinct bft.tag_id) as tag_ids,
                      bv.content, b.is_active, bv.version, b.created_at, bv.updated_at
            from banner b
            join banner_version bv on b.banner_id = bv.banner_id and b.active_version = bv.version
            join banner_feature_tag bft on b.banner_id = bft.banner_id
            where bft.feature_id = $1 or bft.tag_id = $2 and b.must_be_deleted = false
            group by b.banner_id, bft.feature_id, bv.content, b.is_active, bv.version, b.created_at, bv.updated_at
            order by b.banner_id desc 
            limit $3 offset $4`
	)

	var banners []models.Banner
	if err := pgxscan.Select(ctx, b.pool, &banners, selectFilteredBannersQuery, filter.FeatureId, filter.TagId, filter.Limit, filter.Offset); errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return banners, nil
}

func (b *BannerRepository) CreateBanner(ctx context.Context, banner *models.Banner) (uint64, error) {
	const (
		createBannerQuery = `insert into banner (is_active) values ($1) returning banner_id`

		addFeatureAndTagsQuery = `
            insert into banner_feature_tag (banner_id, tag_id, feature_id)
            select $1, r.tag, $3    
            from unnest($2::int[]) as r(tag)`

		createVersionQuery = `insert into banner_version (banner_id, content) values ($1, $2)`
	)

	var bannerId uint64
	err := RunInTx(ctx, b.pool, func(tx pgx.Tx) error {
		if err := pgxscan.Get(ctx, tx, &bannerId, createBannerQuery, banner.IsActive); errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		} else if err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, addFeatureAndTagsQuery, bannerId, banner.TagIds, banner.FeatureId); err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, createVersionQuery, bannerId, string(banner.Content)); err != nil {
			return err
		}

		return nil
	})

	return bannerId, err
}

// PartialUpdateBanner ToDo test this method
func (b *BannerRepository) PartialUpdateBanner(ctx context.Context, bannerId uint64, bannerPartial *models.PatchBanner) error {
	const (
		createNewVersionQuery = `
		    insert into banner_version (banner_id, version, content, updated_at)
		    select $1, (select max(version) + 1 from banner_version), coalesce($2, content), now()
		    from banner_version
		    where banner_id = $1 and version = (select active_version from banner where banner_id=$1)
		    returning version`

		updateActiveVersionQuery = `
            update banner set active_version = $2, is_active = $3
            where banner_id = $1`

		deleteQuery = `
		    delete from banner_feature_tag
            where banner_id = $1`

		addNewTagsQuery = `
		    insert into banner_feature_tag (banner_id, tag_id, feature_id)
		    select $1, unnest($2::int[]), $3`
	)

	err := RunInTx(ctx, b.pool, func(tx pgx.Tx) error {
		var version uint64
		var content *string
		if bannerPartial.Content != nil {
			str := string(bannerPartial.Content)
			content = &str
		}
		if err := pgxscan.Get(ctx, tx, &version, createNewVersionQuery, bannerId, content); errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		} else if err != nil {
			return err
		}

		_, err := b.pool.Exec(ctx, updateActiveVersionQuery, bannerId, version, bannerPartial.IsActive)
		if err != nil {
			return err
		}

		if bannerPartial.TagIds != nil && bannerPartial.FeatureId != nil {
			_, err = b.pool.Exec(ctx, deleteQuery, bannerId)
			if err != nil {
				return err
			}

			_, err = b.pool.Exec(ctx, addNewTagsQuery, bannerId, bannerPartial.TagIds, bannerPartial.FeatureId)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func (b *BannerRepository) DeleteBanner(ctx context.Context, bannerId uint64) error {
	const (
		deleteBannerVersionQuery = `
		    delete from banner_version
       	    where banner_id = $1`

		deleteBannerFeatureTagQuery = `
		    delete from banner_feature_tag
       	    where banner_id = $1`

		deleteBannerQuery = `
		    delete from banner
       	    where banner_id = $1`
	)

	err := RunInTx(ctx, b.pool, func(tx pgx.Tx) error {
		if res, err := tx.Exec(ctx, deleteBannerVersionQuery, bannerId); err != nil {
			return err
		} else if res.RowsAffected() == 0 {
			return ErrNotFound
		}

		if _, err := tx.Exec(ctx, deleteBannerFeatureTagQuery, bannerId); err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, deleteBannerQuery, bannerId); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (b *BannerRepository) MarkBannersAsDeleted(ctx context.Context, featureId, tagId *uint64) error {
	const (
		markBannersAsDeletedQuery = `
		update banner
		set must_be_deleted = true
		where banner_id in
			  (select banner_id
			   from banner_feature_tag
			   where ($1::int is null or feature_id = $1) or ($2::int is null or tag_id = $2))`
	)

	_, err := b.pool.Exec(ctx, markBannersAsDeletedQuery, featureId, tagId)

	return err
}

// DeleteMarkedBanners todo test this method
func (b *BannerRepository) DeleteMarkedBanners(ctx context.Context) error {
	const (
		deleteBannerVersionQuery = `delete from banner where must_be_deleted`

		deleteBannerFeatureTagQuery = `
		delete from banner_version
		using banner
		where must_be_deleted and banner.banner_id = banner_version.banner_id`

		deleteBannerQuery = `
delete from banner_feature_tag
								using banner
								where must_be_deleted and banner.banner_id = banner_feature_tag.banner_id`
	)

	err := RunInTx(ctx, b.pool, func(tx pgx.Tx) error {
		if res, err := tx.Exec(ctx, deleteBannerFeatureTagQuery); err != nil {
			return err
		} else if res.RowsAffected() == 0 {
			return ErrNotFound
		}

		if _, err := tx.Exec(ctx, deleteBannerQuery); err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, deleteBannerVersionQuery); err != nil {
			return err
		}

		return nil
	})

	return err
}
