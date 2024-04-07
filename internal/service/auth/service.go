package auth

import (
	"banner-service/internal/models"
	utils "banner-service/pkg/utils/auth"
	"context"
)

type Repository interface {
	SignUp(ctx context.Context, user *models.User) error
	SignIn(ctx context.Context, signInInput *models.SignInInput) (string, error)
}

type Deps struct {
	AuthRepo Repository
}

type Service struct {
	Deps
}

func NewService(d Deps) *Service {
	return &Service{
		Deps: d,
	}
}

func (s *Service) SignUp(ctx context.Context, user *models.User) error {
	if hashPassword, err := utils.HashPassword(user.Password); err != nil {
		return err
	} else {
		user.Password = hashPassword
	}

	if err := s.AuthRepo.SignUp(ctx, user); err != nil {
		return err
	}

	return nil
}
