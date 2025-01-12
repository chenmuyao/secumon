package logmonitor

import (
	"context"
	"log/slog"
	"time"

	"github.com/chenmuyao/secumon/internal/domain"
	"github.com/chenmuyao/secumon/internal/repository"
	"github.com/chenmuyao/secumon/internal/repository/cache"
)

type BruteForceDetector interface {
	Detect(ctx context.Context, log domain.AccessLog) error
}

type bruteForceDetector struct {
	repo              repository.LogRepo
	bruteForceChecker cache.BruteForceChecker
	alertCache        cache.AlertCache
}

// Detect implements BruteForceDetector.
func (b *bruteForceDetector) Detect(ctx context.Context, log domain.AccessLog) error {
	// Use the brute-force checker to check if it is an attack
	secuEvt, err := b.bruteForceChecker.Check(ctx, log)
	if err != nil {
		// Probably a redis error. Since all the service depends on redis,
		// if an error happens, just don't bother and return.
		return err
	}
	// Not considered as a brute force attack
	if secuEvt.Type == "" {
		return nil
	}
	// In case of an attack
	// Log
	slog.Info("[ALERT] Brute force attack detected", slog.Any("IP", secuEvt.ClientIP))
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

func NewBruteForceDetector(
	repo repository.LogRepo,
	bruteForce cache.BruteForceChecker,
	alertCache cache.AlertCache,
) BruteForceDetector {
	return &bruteForceDetector{
		repo:              repo,
		bruteForceChecker: bruteForce,
		alertCache:        alertCache,
	}
}
