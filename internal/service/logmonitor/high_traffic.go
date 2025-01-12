package logmonitor

import (
	"context"
	"log/slog"
	"time"

	"github.com/chenmuyao/secumon/internal/domain"
	"github.com/chenmuyao/secumon/internal/repository"
	"github.com/chenmuyao/secumon/internal/repository/cache"
)

type HighTrafficDetector struct {
	repo               repository.LogRepo
	highTrafficChecker *cache.HighTrafficChecker
	alertCache         cache.AlertCache
}

// Detect implements HighTrafficDetector.
func (b *HighTrafficDetector) Detect(ctx context.Context, log domain.AccessLog) error {
	// Use the brute-force checker to check if it is an attack
	secuEvt, err := b.highTrafficChecker.Check(ctx, log)
	if err != nil {
		// Probably a redis error. Since all the service depends on redis,
		// if an error happens, just don't bother and return.
		return err
	}
	// Not considered as a high traffic attack
	if secuEvt.Type == "" {
		return nil
	}
	// In case of an attack
	// Log
	slog.Info("[ALERT] high traffic detected", slog.Any("IP", secuEvt.ClientIP))
	// write the DB

	err = b.repo.UpsertSecurityEvent(ctx, secuEvt)
	if err != nil {
		return err
	}

	// async delete the cache
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		er := b.alertCache.DeleteAlerts(ctx, secuEvt.Type)
		if er != nil {
			slog.Info("Failed to delete cache", slog.Any("err", er))
		}
		er = b.alertCache.DeleteAlerts(ctx, "all")
		if er != nil {
			slog.Info("Failed to delete cache", slog.Any("err", er))
		}
	}()
	return nil
}

func NewHighTrafficDetector(
	repo repository.LogRepo,
	highTraffic *cache.HighTrafficChecker,
	alertCache cache.AlertCache,
) Detector {
	return &HighTrafficDetector{
		repo:               repo,
		highTrafficChecker: highTraffic,
		alertCache:         alertCache,
	}
}
