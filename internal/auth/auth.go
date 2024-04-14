package auth

import (
	"banner-service/internal/models"
	"banner-service/internal/repository"
	"context"
	"errors"
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

func (s *Provider) SignUp(ctx context.Context, user *models.User) (string, error) {
	hashPassword, err := HashPassword(user.Password)
	if err != nil {
		return "", err
	}

	user.Password = hashPassword

	err = s.AuthRepo.SignUp(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to sign-in: %w", err)
	}

	resources, err := s.AuthRepo.GetUserResources(ctx, user.Username)
	if err != nil {
		return "", err
	}

	token, err := CreateToken(time.Hour, user.Username, resources, s.PrivateKey)
	if err != nil {
		return "", err
	}

	return token, nil
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
	if errors.Is(err, repository.ErrNotFound) {
		return "", fmt.Errorf("user with name `%s` does not exist", signInInput.Username)
	} else if err != nil {
		return "", err
	}

	token, err := CreateToken(time.Hour, signInInput.Username, resources, s.PrivateKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

type roleKey struct{}

func SetRole(ctx context.Context, role models.UserRole) context.Context {
	return context.WithValue(ctx, roleKey{}, role)
}

func GetRole(ctx context.Context) models.UserRole {
	return ctx.Value(roleKey{}).(models.UserRole)
}
