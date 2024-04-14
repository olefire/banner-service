package app

import (
	AuthProvider "banner-service/internal/auth"
	"banner-service/internal/config"
	controllerhttp "banner-service/internal/controller/http"
	"banner-service/internal/middleware"
	"banner-service/internal/repository"
	BannerService "banner-service/internal/service/banner"
	"banner-service/internal/worker"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func Start() {
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

	authRepo := repository.NewAuthRepository(pool)
	bannerRepo := repository.NewBannerRepository(pool)

	authService := AuthProvider.NewAuthProvider(AuthProvider.Deps{
		AuthRepo:   authRepo,
		PrivateKey: cfg.PrivateKey,
		PublicKey:  cfg.PublicKey,
	})

	bannerService := BannerService.NewService(BannerService.Deps{BannerRepo: bannerRepo})
	bannerTicker := worker.NewBannerCollector(bannerRepo)
	go bannerTicker.Start(ctx)

	ctr := controllerhttp.NewController(
		controllerhttp.AuthProvider{AuthManagement: authService},
		controllerhttp.BannerService{BannerManagement: bannerService},
		cfg.PublicKey)

	router := ctr.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
	})

	err = http.ListenAndServe(cfg.Port, middleware.PanicRecovery(middleware.LogRequest(c.Handler(router))))
	if err != nil {
		log.Fatalf("Unable to start server: %v\n", err)
	}

}