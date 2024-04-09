package auth

import (
	"banner-service/internal/models"
	"context"
)

type AuthManagement interface {
	SignIn(ctx context.Context, signInInput *models.User) (string, error)
	SignUp(ctx context.Context, user *models.User) error
}
