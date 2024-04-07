package repository

import (
	"banner-service/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type AuthRepository struct {
	pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		pool: pool,
	}
}

func (au *AuthRepository) SignUp(ctx context.Context, user *models.User) error {
	query := "INSERT INTO user(username, hash_password, is_admin) VALUES($1, $2, $3)"
	_, err := au.pool.Exec(ctx, query, user.Username, user.Password, user.IsAdmin)
	if err != nil {
		return err
	}
	return nil
}

func (au *AuthRepository) GetPassword(ctx context.Context, signInInput *models.SignInInput) (string, error) {
	var hashPassword string
	query := "SELECT hash_password FROM user WHERE username = $1"
	if err := au.pool.QueryRow(ctx, query, signInInput.Username).Scan(&hashPassword); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("no such user")
		}
		return "", err
	}
	return hashPassword, nil
}
