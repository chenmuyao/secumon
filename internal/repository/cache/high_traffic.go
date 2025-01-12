package cache

import (
	"context"
	"time"

	"github.com/chenmuyao/secumon/internal/domain"
	"github.com/redis/go-redis/v9"
)

type HighTrafficChecker struct {
	checker RateLimitChecker
}

// Check implements HighTrafficChecker.
func (r *HighTrafficChecker) Check(
	ctx context.Context,
	log domain.AccessLog,
) (domain.SecurityEvent, error) {
	res, err := r.checker.Check(ctx, log)
	if err != nil {
		return domain.SecurityEvent{}, err
	}
	res.Details = "10 consecutive requests within 1 minute"
	return res, nil
}

func NewHighTrafficChecker(cmd redis.Cmdable) *HighTrafficChecker {
	checker := NewRateLimitChecker(cmd, "hightraffic", 10, int64(time.Minute))
	return &HighTrafficChecker{
		checker: checker,
	}
}
