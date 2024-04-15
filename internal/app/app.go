package app

import (
	AuthProvider "banner-service/internal/auth"
	"banner-service/internal/config"
	controllerhttp "banner-service/internal/controller/http"
	"banner-service/internal/middleware"
	"banner-service/internal/models"
	"banner-service/internal/repository"
	BannerService "banner-service/internal/service/banner"
	"banner-service/internal/worker"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jellydator/ttlcache/v3"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func Start() {
	ctx := context.Background()
	cfg := config.NewConfig()
	pool, err := pgxpool.New(ctx, cfg.PgDSN)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		fmt.Print(err)
	}

	authRepo := repository.NewAuthRepository(pool)
	tokenCache := ttlcache.New[string, models.UserResources](ttlcache.
		WithTTL[string, models.UserResources](cfg.TokenTTL),
		ttlcache.WithCapacity[string, models.UserResources](cfg.TokenCacheCapacity),
	)
	bannerCache := ttlcache.New[models.FeatureTag, models.BannerContent](ttlcache.
		WithTTL[models.FeatureTag, models.BannerContent](cfg.BannerTTL),
		ttlcache.WithCapacity[models.FeatureTag, models.BannerContent](cfg.TokenCacheCapacity),
	)

	tokenProvider := AuthProvider.NewTokenProvider(tokenCache, cfg.PrivateKey, cfg.PublicKey, cfg.TokenTTL)

	bannerRepo := repository.NewBannerRepository(pool)

	authService := AuthProvider.NewAuthProvider(AuthProvider.Deps{
		AuthRepo:      authRepo,
		TokenProvider: tokenProvider,
	})

	bannerService := BannerService.NewService(BannerService.Deps{BannerRepo: bannerRepo, Cache: bannerCache})
	bannerTicker := worker.NewBannerCollector(bannerRepo)
	go bannerTicker.Start(ctx)

	ctr := controllerhttp.NewController(
		controllerhttp.AuthProvider{AuthManagement: authService, TokenProvider: tokenProvider},
		controllerhttp.BannerService{BannerManagement: bannerService},
	)

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
