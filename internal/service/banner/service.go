package banner

import (
	"banner-service/internal/models"
	"context"
)

type Repository interface {
	GetBanner(ctx context.Context, tagId uint64, featureId uint64) (*models.Banner, error)
	GetFilteredBanners(ctx context.Context, filter *models.FilterBanner) ([]models.Banner, error)
	CreateBanner(ctx context.Context, banner *models.Banner) (uint64, error)
	PartialUpdateBanner(ctx context.Context, bannerPartial *models.Banner) error
	DeleteBanner(ctx context.Context, id uint64) error
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

func (s *Service) GetBanner(ctx context.Context, tagId uint64, featureId uint64) (*models.Banner, error) {
	panic("implement me")
}

func (s *Service) GetFilteredBanners(ctx context.Context, filter *models.FilterBanner) ([]models.Banner, error) {
	panic("implement me")
}

func (s *Service) CreateBanner(ctx context.Context, banner *models.Banner) (uint64, error) {
	panic("implement me")
}

func (s *Service) PartialUpdateBanner(ctx context.Context, bannerPartial *models.Banner) error {
	panic("implement me")
}

func (s *Service) DeleteBanner(ctx context.Context, bannerId uint64) error {
	panic("implement me")
}
