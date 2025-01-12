package cache

import (
	"context"
	"time"

	"github.com/chenmuyao/secumon/internal/domain"
	"github.com/redis/go-redis/v9"
)

type BruteForceChecker interface {
	Check(ctx context.Context, log domain.AccessLog) (domain.SecurityEvent, error)
}

type RedisBruteForceChecker struct {
	cmd redis.Cmdable
}

// Check implements BruteForceChecker.
func (r *RedisBruteForceChecker) Check(
	ctx context.Context,
	log domain.AccessLog,
) (domain.SecurityEvent, error) {
	return domain.SecurityEvent{
		Type:      "bruteforce",
		Timestamp: time.Now(),
		ClientIP:  "192.268.1.1",
		Details:   "details",
	}, nil
}

func NewBruteForceChecker(cmd redis.Cmdable) BruteForceChecker {
	return &RedisBruteForceChecker{}
}
