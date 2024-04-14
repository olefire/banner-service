package worker

import (
	"context"
	"log"
	"time"
)

type repository interface {
	DeleteMarkedBanners(ctx context.Context) error
}

type BannerCollector struct {
	repository repository
}

func NewBannerCollector(bannerRepo repository) *BannerCollector {
	return &BannerCollector{
		repository: bannerRepo,
	}
}

func (w *BannerCollector) Start(ctx context.Context) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := w.repository.DeleteMarkedBanners(ctx)
			if err != nil {
				log.Printf("worker: %v", err)
			}
		case <-ctx.Done():
			return

		}
	}
}
