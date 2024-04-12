package http

import (
	"banner-service/internal/models"
	"context"
)

type AuthManagement interface {
	SignIn(ctx context.Context, signInInput *models.User) (string, error)
	SignUp(ctx context.Context, user *models.User) error
}

type BannerManagement interface {
	GetBanner(ctx context.Context, tagId uint64, featureId uint64, role models.UserRole) (string, error)
	GetFilteredBanners(ctx context.Context, filter *models.FilterBanner) ([]models.Banner, error)
	CreateBanner(ctx context.Context, banner *models.Banner) (uint64, error)
	PartialUpdateBanner(ctx context.Context, bannerId uint64, bannerPartial *models.PatchBanner) error
	DeleteBanner(ctx context.Context, id uint64) error
}