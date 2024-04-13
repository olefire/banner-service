package auth

import (
	"banner-service/internal/controller/http"
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
	GetUserResources(ctx context.Context, username string) (models.UserResources, error)
}

type Deps struct {
	AuthRepo   Repository
	PrivateKey string
	PublicKey  string
}

type Provider struct {
	Deps
}

func NewAuthProvider(d Deps) *Provider {
	return &Provider{
		Deps: d,
	}
}

var _ http.AuthManagement = (*Provider)(nil)

func (s *Provider) SignUp(ctx context.Context, user *models.User) error {
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

func (s *Provider) SignIn(ctx context.Context, signInInput *models.User) (string, error) {
	hashPassword, err := s.AuthRepo.GetHashPassword(ctx, signInInput.Username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(signInInput.Password))
	if err != nil {
		return "", fmt.Errorf("invalid login or password: user=%v", signInInput)
	}

	resources, err := s.AuthRepo.GetUserResources(ctx, signInInput.Username)
	if err != nil {
		return "", err
	}

	token, err := utils.CreateToken(time.Hour, resources, s.PrivateKey)
	if err != nil {
		return "", err
	}

	return token, nil
}
