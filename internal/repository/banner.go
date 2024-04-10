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

func (b *BannerRepository) GetBannerIsActive(ctx context.Context, tagId uint64, featureId uint64) error {
	query := `SELECT b.is_active
           FROM banner_feature_tag bft
           JOIN banner b ON bft.banner_id = b.banner_id
           WHERE bft.feature_id = $1 AND bft.tag_id = $2`

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
	query := `SELECT bv.content
           FROM banner_feature_tag bft
           JOIN banner b ON bft.banner_id = b.banner_id
           JOIN banner_version bv ON b.banner_id = bv.banner_id AND b.active_version = bv.version
           WHERE bft.feature_id = $1 AND bft.tag_id = $2`

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

// CreateBanner TODO Refactor this method
func (b *BannerRepository) CreateBanner(ctx context.Context, banner *models.Banner) (uint64, error) {
	tx, err := b.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	//ToDo: handle rollback error
	defer tx.Rollback(ctx)

	query := `INSERT INTO banner (is_active) VALUES ($1) RETURNING banner_id`
	var bannerId uint64
	err = pgxscan.Get(ctx, tx, &bannerId, query, banner.IsActive)
	if err != nil {
		return 0, fmt.Errorf("failed to insert banner: %w", err)
	}

	for _, tagId := range banner.TagIds {
		query = `INSERT INTO banner_feature_tag (banner_id, tag_id, feature_id) VALUES ($1, $2, $3)`
		_, err = tx.Exec(ctx, query, bannerId, tagId, banner.FeatureId)
		if err != nil {
			return 0, err
		}
	}

	var version uint64
	query = `INSERT INTO banner_version (banner_id, content, version) VALUES ($1, $2) RETURNING version`
	err = pgxscan.Get(ctx, tx, &version, query, bannerId, string(banner.Content))
	if err != nil {
		return 0, fmt.Errorf("failed to insert banner version: %w", err)
	}

	query = `UPDATE banner SET active_version = $1 WHERE banner_id = $2`
	_, err = tx.Exec(ctx, query, version, bannerId)
	if err != nil {
		return 0, fmt.Errorf("failed to set banner: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	return bannerId, nil
}

// PartialUpdateBanner TODO Refactor this method
func (b *BannerRepository) PartialUpdateBanner(ctx context.Context, bannerPartial *models.PatchBanner) error {
	//tx, err := b.pool.Begin(ctx)
	//if err != nil {
	//	return err
	//}
	////TODO: handle rollback error
	//defer tx.Rollback(ctx)
	//
	//if bannerPartial.IsActive.Valid {
	//	query := `UPDATE banner b
	//	          SET is_active = $1
	//	          FROM banner_feature_tag bft
	//	          WHERE b.banner_id = bft.banner_id AND bft.feature_id = $2 AND bft.tag_id = $3`
	//	_, err = tx.Exec(ctx, query, bannerPartial.IsActive.Bool, bannerPartial.FeatureId, tagId)
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//if bannerPartial.Content != "" {
	//	query := `UPDATE banner_version bv
	//	          SET content = $1
	//	          FROM banner b, banner_feature_tag bft
	//	          WHERE b.banner_id = bv.banner_id AND b.banner_id = bft.banner_id AND bft.feature_id = $2 AND bft.tag_id = $3 AND b.active_version = bv.version`
	//	_, err = tx.Exec(ctx, query, bannerPartial.Content, featureId, tagId)
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//err = tx.Commit(ctx)
	//if err != nil {
	//	return err
	//}

	return nil
}
func (b *BannerRepository) DeleteBanner(ctx context.Context, id uint64) error {
	//TODO implement me
	panic("implement me")
}
