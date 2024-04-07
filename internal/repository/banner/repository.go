package repository

import "github.com/jackc/pgx/v5"

type BannerRepository struct {
	collection *pgx.Conn
}

func NewBannerRepository(col *pgx.Conn) *BannerRepository {
	return &BannerRepository{
		collection: col,
	}
}
