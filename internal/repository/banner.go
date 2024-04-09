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

func (b BannerRepository) GetBannerIsActive(ctx context.Context, tagId uint64, featureId uint64) error {
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

func (b BannerRepository) GetBanner(ctx context.Context, tagId uint64, featureId uint64) (string, error) {
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

func (b BannerRepository) GetFilteredBanners(ctx context.Context, filter *models.FilterBanner) ([]models.Banner, error) {
	//TODO implement me
	panic("implement me")
}

func (b BannerRepository) CreateBanner(ctx context.Context, banner *models.Banner) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (b BannerRepository) PartialUpdateBanner(ctx context.Context, bannerPartial *models.Banner) error {
	//TODO implement me
	panic("implement me")
}

func (b BannerRepository) DeleteBanner(ctx context.Context, id uint64) error {
	//TODO implement me
	panic("implement me")
}

func NewBannerRepository(p *pgxpool.Pool) *BannerRepository {
	return &BannerRepository{
		pool: p,
	}
}
