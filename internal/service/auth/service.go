package auth

import (
	"banner-service/internal/models"
	utils "banner-service/pkg/utils/auth"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Repository interface {
	SignUp(ctx context.Context, user *models.User) error
	GetHashPassword(ctx context.Context, username string) (string, error)
	GetRole(ctx context.Context, username string) (string, error)
}

type Deps struct {
	AuthRepo   Repository
	PrivateKey string
	PublicKey  string
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
	hashPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashPassword

	if err := s.AuthRepo.SignUp(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *Service) SignIn(ctx context.Context, signInInput *models.User) (string, error) {
	hashPassword, err := s.AuthRepo.GetHashPassword(ctx, signInInput.Username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(signInInput.Password))
	if err != nil {
		return "", fmt.Errorf("invalid password: user=%v", signInInput)
	}

	role, err := s.AuthRepo.GetRole(ctx, signInInput.Username)
	if err != nil {
		return "", fmt.Errorf("failed to get role: user=%v", signInInput)
	}

	token, err := utils.CreateToken(time.Hour, role, s.PrivateKey)
	if err != nil {
		return "", err
	}

	return token, nil
}
