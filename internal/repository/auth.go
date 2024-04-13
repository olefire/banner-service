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
	const (
		insertUserQuery = `INSERT INTO users(username, hash_password, role) VALUES($1, $2, $3);`
	)

	if _, err := au.pool.Exec(ctx, insertUserQuery, user.Username, user.Password, user.Role); errors.Is(err, pgx.ErrNoRows) {
		return ErrAlreadyExists
	} else if err != nil {
		return err
	}

	return nil
}

func (au *AuthRepository) GetHashPassword(ctx context.Context, username string) (string, error) {
	const (
		getHashPasswordQuery = `select hash_password from users where username = $1;`
	)

	var passwordHash string
	if err := pgxscan.Get(ctx, au.pool, &passwordHash, getHashPasswordQuery, username); errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	} else if err != nil {
		return "", err
	}

	return passwordHash, nil
}

func (au *AuthRepository) GetUserResources(ctx context.Context, username string) (models.UserResources, error) {
	const (
		getResourcesQuery = `SELECT u.role, array_agg(re.resource) AS resources
        FROM users u
        JOIN role_endpoints re USING (role)
        WHERE u.username = $1
        GROUP BY u.role;`
	)

	var resources models.UserResources
	if err := pgxscan.Get(ctx, au.pool, &resources, getResourcesQuery, username); errors.Is(err, pgx.ErrNoRows) {
		return models.UserResources{}, ErrNotFound
	} else if err != nil {
		return models.UserResources{}, err
	}

	return resources, nil
}
