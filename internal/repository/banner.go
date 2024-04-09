package repository

import (
	"banner-service/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BannerRepository struct {
	pool *pgxpool.Pool
}

func (b BannerRepository) GetBanner(ctx context.Context, tagId uint64, featureId uint64) (*models.Banner, error) {
	//TODO implement me
	panic("implement me")
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
