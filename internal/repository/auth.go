package repository

import (
	"banner-service/internal/models"
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
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
	query := "INSERT INTO users(username, hash_password, role) VALUES($1, $2, $3)"

	if _, err := au.pool.Exec(ctx, query, user.Username, user.Password, user.Role); err != nil {
		return err
	}

	return nil
}

func (au *AuthRepository) GetHashPassword(ctx context.Context, username string) (string, error) {
	query := "SELECT hash_password FROM users WHERE username = $1"

	var passwordHash string
	if err := pgxscan.Get(ctx, au.pool, &passwordHash, query, username); errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	} else if err != nil {
		return "", err
	}

	return passwordHash, nil
}

func (au *AuthRepository) GetRole(ctx context.Context, username string) (string, error) {
	query := "SELECT role FROM users WHERE username = $1"

	var role string
	if err := pgxscan.Get(ctx, au.pool, &role, query, username); errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	} else if err != nil {
		return "", err
	}

	return role, nil
}
