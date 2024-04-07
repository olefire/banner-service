package service

import (
	"banner-service/internal/models"
	"context"
)

type Repository interface {
	GetBanner(ctx context.Context, tagId uint, featureId uint, useLastVersion bool) (models.Banner, error)
	GetFilteredBanners(ctx context.Context, filter models.FilterBanner) ([]models.Banner, error)
	CreateBanner(ctx context.Context, banner models.Banner) (uint, error)
	PartialUpdateBanner(ctx context.Context, bannerPartial models.Banner) error
	DeleteBanner(ctx context.Context, id uint) error
}

type Deps struct {
	BannerRepo Repository
}

type Service struct {
	Deps
}

func NewService(d Deps) *Service {
	return &Service{
		Deps: d,
	}
}
