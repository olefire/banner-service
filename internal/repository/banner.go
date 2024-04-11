package repository

import (
	"banner-service/internal/models"
	"context"
	"errors"
	"fmt"
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
	query := `select b.is_active
           from banner_feature_tag bft
           join banner b on bft.banner_id = b.banner_id
           where bft.feature_id = $1 and bft.tag_id = $2`

	var isActive bool
	if err := pgxscan.Get(ctx, b.pool, &isActive, query, featureId, tagId); errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	} else if err != nil {
		return err
	}

	if !isActive {
		return ErrAccessDenied
	}
	return nil
}

func (b *BannerRepository) GetBanner(ctx context.Context, tagId uint64, featureId uint64) (string, error) {
	query := `select bv.content
           from banner_feature_tag bft
           join banner b on bft.banner_id = b.banner_id
           join banner_version bv on b.banner_id = bv.banner_id and b.active_version = bv.version
           where bft.feature_id = $1 and bft.tag_id = $2`

	var content string
	if err := pgxscan.Get(ctx, b.pool, &content, query, featureId, tagId); errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	} else if err != nil {
		return "", err
	}

	return content, nil
}

func (b *BannerRepository) GetFilteredBanners(ctx context.Context, filter *models.FilterBanner) ([]models.Banner, error) {
	//TODO implement me
	panic("implement me")
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
		if err := pgxscan.Get(ctx, tx, &bannerId, createBannerQuery); err != nil {
			return fmt.Errorf("failed to insert banner: %w", err)
		}

		if _, err := tx.Exec(ctx, addFeatureAndTagsQuery, bannerId, banner.TagIds, banner.FeatureId); err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, createVersionQuery, bannerId, string(banner.Content)); err != nil {
			return fmt.Errorf("failed to insert banner version: %w", err)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return bannerId, nil
}

// PartialUpdateBanner TODO Refactor this method
func (b *BannerRepository) PartialUpdateBanner(ctx context.Context, bannerId uint64, bannerPartial *models.PatchBanner) error {
	const (
		getActiveVersionQuery    = `select active_version from banner where banner_id=$1`
		createNewVersionQuery    = `insert into banner_version (banner_id, version, content, updated_at) select $1, (select max(version) from banner_version), coalesce($3, content), now() from banner_version where banner_id = $1 and version = $2`
		updateActiveVersionQuery = `update banner set active_version = $1 where banner_id = $2`
	)

	err := RunInTx(ctx, b.pool, func(tx pgx.Tx) error {
		var activeVersion int
		if err := pgxscan.Get(ctx, b.pool, &activeVersion, getActiveVersionQuery, bannerId); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrNotFound
			}
			return err
		}

		_, err := b.pool.Exec(ctx, createNewVersionQuery, bannerId, activeVersion, string(bannerPartial.Content))
		if err != nil {
			return err
		}

		_, err = b.pool.Exec(ctx, updateActiveVersionQuery, activeVersion+1, bannerId)
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

func (b *BannerRepository) MarkBannerAsDeleted(ctx context.Context, tagId uint64, featureId uint64) error {
	const (
		markBannerAsDeletedQuery = `update banner set is_active = false where banner_id in (select banner_id from banner_feature_tag where tag_id = $1 or feature_id = $2)`
	)

	_, err := b.pool.Exec(ctx, markBannerAsDeletedQuery, tagId, featureId)

	return err
}

func (b *BannerRepository) DeleteBannerByTagOrFeature(ctx context.Context, tagId uint64, featureId uint64) error {
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
