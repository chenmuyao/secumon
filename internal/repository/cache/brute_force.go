package cache

import (
	"context"
	"net/http"
	"time"

	"github.com/chenmuyao/secumon/internal/domain"
	"github.com/redis/go-redis/v9"
)

type BruteForceChecker struct {
	checker RateLimitChecker
}

// Check implements BruteForceChecker.
func (r *BruteForceChecker) Check(
	ctx context.Context,
	log domain.AccessLog,
) (domain.SecurityEvent, error) {
	// only check if the status is unauthorized
	if log.StatusCode == http.StatusUnauthorized {
		res, err := r.checker.Check(ctx, log)
		if err != nil {
			return domain.SecurityEvent{}, err
		}
		res.Details = "5 consecutive 401 errors within 1 minute"
		return res, nil
	}
	return domain.SecurityEvent{}, nil
}

func NewBruteForceChecker(cmd redis.Cmdable) *BruteForceChecker {
	checker := NewRateLimitChecker(cmd, "bruteforce", 5, int64(time.Minute))
	return &BruteForceChecker{
		checker: checker,
	}
}
