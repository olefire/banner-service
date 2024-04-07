package main

import (
	"banner-service/internal/config"
	authRepo "banner-service/internal/repository/auth"
	bannerRepo "banner-service/internal/repository/banner"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func main() {
	ctx := context.Background()
	cfg := config.NewConfig()
	pool, err := pgxpool.New(ctx, cfg.PgURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		fmt.Print(err)
	}
	authRepo := authRepo.NewAuthRepository(conn)
	bannerRepo := bannerRepo.NewBannerRepository(conn)

}
