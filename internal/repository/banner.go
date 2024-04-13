package repository

import (
	"banner-service/internal/models"
	"context"
	"errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BannerRepository struct {
	pool *pgxpool.Pool
}

func NewBannerRepository(p *pgxpool.Pool) *BannerRepository {
	return &BannerRepository{
		pool: p,
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

func (b *BannerRepository) GetBannerIsActive(ctx context.Context, tagId uint64, featureId uint64) error {
	const (
		selectBannerIsActiveQuery = `select b.is_active
           from banner_feature_tag bft
           join banner b on bft.banner_id = b.banner_id
           where bft.feature_id = $1 and bft.tag_id = $2`
	)

	var isActive bool
	if err := pgxscan.Get(ctx, b.pool, &isActive, selectBannerIsActiveQuery, featureId, tagId); errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	} else if err != nil {
		return err
	}

	if !isActive {
		return ErrBannerInactive
	}
	return nil
}

func (b *BannerRepository) GetBanner(ctx context.Context, tagId uint64, featureId uint64) (string, error) {
	const (
		selectBannerQuery = `select bv.content
           from banner_feature_tag bft
           join banner b using (banner_id)
           join banner_version bv on b.banner_id = bv.banner_id and b.active_version = bv.version
           where bft.feature_id = $1 and bft.tag_id = $2`
	)

	var content string
	if err := pgxscan.Get(ctx, b.pool, &content, selectBannerQuery, featureId, tagId); errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	} else if err != nil {
		return "", err
	}

	return content, nil
}

func (b *BannerRepository) GetListOfVersions(ctx context.Context, bannerId uint64) ([]models.Banner, error) {
	const (
		selectBannersQuery = `SELECT b.banner_id, bft.feature_id, array_agg(DISTINCT bft.tag_id) AS tag_ids,
              bv.content, b.is_active, bv.version, b.created_at, bv.updated_at
         FROM banner_version bv
         JOIN banner b USING (banner_id)
         JOIN banner_feature_tag bft USING (banner_id)
         WHERE banner_id = $1
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

func (b *BannerRepository) ChooseVersion(ctx context.Context, bannerId uint64, version uint64) error {
	const (
		chooseVersionQuery = `update banner
						 set active_version = $2
						 where banner_id = $1`
	)

	_, err := b.pool.Exec(ctx, chooseVersionQuery, bannerId, version)
	return err
}

func (b *BannerRepository) GetFilteredBanners(ctx context.Context, filter *models.FilterBanner) ([]models.Banner, error) {
	const (
		selectFilteredBannersQuery = `
   SELECT b.banner_id, bft.feature_id, array_agg(DISTINCT bft.tag_id) AS tag_ids,
              bv.content, b.is_active, bv.version, b.created_at, bv.updated_at
    FROM banner b
    JOIN banner_version bv ON b.banner_id = bv.banner_id AND b.active_version = bv.version
    JOIN banner_feature_tag bft ON b.banner_id = bft.banner_id
    WHERE ($1::int IS NULL OR bft.feature_id = $1)
    AND ($2::int IS NULL OR bft.tag_id = $2)
    GROUP BY b.banner_id, bft.feature_id, bv.content, b.is_active, bv.version, b.created_at, bv.updated_at
    LIMIT $3 OFFSET $4`
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
		createBannerQuery = `insert into banner default values returning banner_id`

		addFeatureAndTagsQuery = `
        insert into banner_feature_tag (banner_id, tag_id, feature_id)
        select $1, r.tag, $3
        from unnest($2::int[]) as r(tag)`

		createVersionQuery = `insert into banner_version (banner_id, content) values ($1, $2)`
	)

	var bannerId uint64
	err := RunInTx(ctx, b.pool, func(tx pgx.Tx) error {
		if err := pgxscan.Get(ctx, tx, &bannerId, createBannerQuery); errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		} else if err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, addFeatureAndTagsQuery, bannerId, banner.TagIds, banner.FeatureId); errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		} else if err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, createVersionQuery, bannerId, string(banner.Content)); errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		} else if err != nil {
			return err
		}

		return nil
	})

	return bannerId, err
}

// PartialUpdateBanner ToDo test this method
func (b *BannerRepository) PartialUpdateBanner(ctx context.Context, bannerId uint64, bannerPartial *models.PatchBanner) error {
	const (
		createNewVersionQuery = `insert into banner_version
    								(banner_id, version, content, updated_at)
									select $1, (select max(version) + 1 from banner_version), coalesce($2, content), now()
									from banner_version
									where banner_id = $1 and version =
									                         (select active_version from banner where banner_id=$1)
									returning version`

		updateActiveVersionQuery = `update banner set active_version = $1 where banner_id = $2`

		deleteQuery = `delete from banner_feature_tag where feature_id = $1 or tag_id = any($2)`

		addNewTagsQuery = `INSERT INTO banner_feature_tag (banner_id, tag_id, feature_id)
                     SELECT $1, unnest($2::int[]), $3`
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

		_, err := b.pool.Exec(ctx, updateActiveVersionQuery, version, bannerId)
		if err != nil {
			return err
		}

		_, err = b.pool.Exec(ctx, deleteQuery, bannerPartial.FeatureId, bannerPartial.TagIds)
		if err != nil {
			return err
		}

		_, err = b.pool.Exec(ctx, addNewTagsQuery, bannerId, bannerPartial.TagIds, bannerPartial.FeatureId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (b *BannerRepository) DeleteBanner(ctx context.Context, bannerId uint64) error {
	const (
		deleteBannerVersionQuery    = `DELETE FROM banner_version WHERE banner_id = $1`
		deleteBannerFeatureTagQuery = `DELETE FROM banner_feature_tag WHERE banner_id = $1`
		deleteBannerQuery           = `DELETE FROM banner WHERE banner_id = $1`
	)

	err := RunInTx(ctx, b.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, deleteBannerVersionQuery, bannerId); err != nil {
			return err
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

// MarkBannerAsDeleted ToDo test this method
func (b *BannerRepository) MarkBannerAsDeleted(ctx context.Context, tagId uint64, featureId uint64) error {
	const (
		markBannerAsDeletedQuery = `update banner
									set must_be_deleted = false
									where banner_id in
									      (select banner_id
									       from banner_feature_tag
									       where tag_id = $1 or feature_id = $2)`
	)

	_, err := b.pool.Exec(ctx, markBannerAsDeletedQuery, tagId, featureId)

	return err
}

// DeleteMarkedBanners todo test this method
func (b *BannerRepository) DeleteMarkedBanners(ctx context.Context) error {
	const (
		deleteBannerVersionQuery = `delete from banner where must_be_deleted`

		deleteBannerFeatureTagQuery = `delete from banner_version
										using banner
										where must_be_deleted and banner.banner_id = banner_version.banner_id`

		deleteBannerQuery = `delete from banner_feature_tag
								using banner
								where must_be_deleted and banner.banner_id = banner_feature_tag.banner_id`
	)

	err := RunInTx(ctx, b.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, deleteBannerFeatureTagQuery); err != nil {
			return err
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
