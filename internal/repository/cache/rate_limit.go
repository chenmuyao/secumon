package cache

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"time"

	"github.com/chenmuyao/secumon/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RateLimitChecker interface {
	Check(ctx context.Context, log domain.AccessLog) (domain.SecurityEvent, error)
}

//go:embed lua/sliding_window_rate_limit.lua
var luaScript string

type RedisRateLimitChecker struct {
	cmd        redis.Cmdable
	limit      int
	windowSize int64
	name       string
}

// Check implements RateLimitChecker.
func (r *RedisRateLimitChecker) Check(
	ctx context.Context,
	log domain.AccessLog,
) (domain.SecurityEvent, error) {
	res, err := r.cmd.Eval(
		ctx,
		luaScript,
		[]string{r.Key(log.ClientIP)},
		r.limit,
		r.windowSize,
		time.Now().UnixNano(),
	).Int()
	if err != nil {
		slog.Error("redis error", "err", err)
		return domain.SecurityEvent{}, err
	}
	switch res {
	case -1: // reached the limit
		return domain.SecurityEvent{
			Type:      r.name,
			Timestamp: time.Now(),
			ClientIP:  log.ClientIP,
		}, nil
	default: // not yet
		return domain.SecurityEvent{}, nil
	}
}

func (r *RedisRateLimitChecker) Key(ip string) string {
	return fmt.Sprintf("%s:%s", r.name, ip)
}

func NewRateLimitChecker(
	cmd redis.Cmdable,
	name string,
	limit int,
	windowSize int64,
) RateLimitChecker {
	return &RedisRateLimitChecker{
		cmd:        cmd,
		limit:      limit,
		windowSize: windowSize,
		name:       name,
	}
}
