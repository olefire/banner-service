package banner

import (
	"banner-service/internal/controller/http"
	"banner-service/internal/models"
	"banner-service/internal/repository"
	"context"
	"github.com/jellydator/ttlcache/v3"
	"log"
)

type Repository interface {
	GetBanner(ctx context.Context, tagId, featureId uint64, isAdmin bool) (models.BannerContent, error)
	GetListOfVersions(ctx context.Context, bannerId uint64) ([]models.Banner, error)
	ChooseBannerVersion(ctx context.Context, bannerId uint64, version uint64) error
	GetFilteredBanners(ctx context.Context, filter *models.FilterBanner) ([]models.Banner, error)
	CreateBanner(ctx context.Context, banner *models.Banner) (uint64, error)
	PartialUpdateBanner(ctx context.Context, bannerId uint64, bannerPartial *models.PatchBanner) error
	DeleteBanner(ctx context.Context, id uint64) error
	MarkBannersAsDeleted(ctx context.Context, featureId, tagId *uint64) error
}

type Deps struct {
	BannerRepo Repository
	Cache      *ttlcache.Cache[models.FeatureTag, models.BannerContent]
}

type Service struct {
	Deps
}

func NewService(d Deps) *Service {
	return &Service{
		Deps: d,
	}
}

var _ http.BannerManagement = (*Service)(nil)

func (s *Service) GetBanner(ctx context.Context, tagId uint64, featureId uint64, role models.UserRole, useLastRevision bool) (string, error) {
	if !useLastRevision {
		if banner := s.Cache.Get(models.FeatureTag{FeatureId: featureId, TagId: tagId}); banner != nil {
			if banner.Value().IsActive || role == models.Admin {
				log.Println("get banner from cache", banner.Value())
				return banner.Value().Content, nil
			} else {
				return "", repository.ErrBannerInactive
			}
		}
	}

	content, err := s.BannerRepo.GetBanner(ctx, tagId, featureId, role == models.Admin)
	if err != nil {
		return "", err
	}

	s.Cache.Set(models.FeatureTag{FeatureId: featureId, TagId: tagId}, content, ttlcache.DefaultTTL)

	return content.Content, nil
}

func (s *Service) GetListOfVersions(ctx context.Context, bannerId uint64) ([]models.Banner, error) {
	if banners, err := s.BannerRepo.GetListOfVersions(ctx, bannerId); err != nil {
		return nil, err
	} else {
		return banners, nil
	}
}

func (s *Service) ChooseBannerVersion(ctx context.Context, bannerId uint64, version uint64) error {
	if err := s.BannerRepo.ChooseBannerVersion(ctx, bannerId, version); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetFilteredBanners(ctx context.Context, filter *models.FilterBanner) ([]models.Banner, error) {
	if banners, err := s.BannerRepo.GetFilteredBanners(ctx, filter); err != nil {
		return nil, err
	} else {
		return banners, nil
	}
}

func (s *Service) CreateBanner(ctx context.Context, banner *models.Banner) (uint64, error) {
	bannerId, err := s.BannerRepo.CreateBanner(ctx, banner)
	if err != nil {
		return 0, err
	}
	return bannerId, nil
}

func (s *Service) PartialUpdateBanner(ctx context.Context, bannerId uint64, bannerPartial *models.PatchBanner) error {
	if err := s.BannerRepo.PartialUpdateBanner(ctx, bannerId, bannerPartial); err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteBanner(ctx context.Context, bannerId uint64) error {
	if err := s.BannerRepo.DeleteBanner(ctx, bannerId); err != nil {
		return err
	}
	return nil
}

func (s *Service) MarkBannerAsDeleted(ctx context.Context, featureId, tagId *uint64) error {
	if err := s.BannerRepo.MarkBannersAsDeleted(ctx, featureId, tagId); err != nil {
		return err
	}
	return nil
}
