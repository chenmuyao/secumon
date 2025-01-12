package logmonitor

import (
	"context"

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
}

// Detect implements BruteForceDetector.
func (b *bruteForceDetector) Detect(ctx context.Context, log domain.AccessLog) error {
	panic("unimplemented")
}

func NewBruteForceDetector(
	repo repository.LogRepo,
	bruteForce cache.BruteForceChecker,
) BruteForceDetector {
	return &bruteForceDetector{
		repo:              repo,
		bruteForceChecker: bruteForce,
	}
}
