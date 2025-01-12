package cache

import (
	"context"
	"errors"

	"github.com/chenmuyao/secumon/internal/domain"
	"github.com/redis/go-redis/v9"
)

type AlertCache interface {
	GetAlerts(ctx context.Context, alertType string) ([]domain.Alert, error)
	SetAlerts(ctx context.Context, alertType string, alerts []domain.Alert) error
	DeleteAlerts(ctx context.Context, alertType string) error
}

type RedisAlertCache struct {
	cmd redis.Cmdable
}

// DeleteAlerts implements AlertCache.
func (r *RedisAlertCache) DeleteAlerts(ctx context.Context, alertType string) error {
	panic("unimplemented")
}

// SetAlerts implements AlertCache.
func (r *RedisAlertCache) SetAlerts(
	ctx context.Context,
	alertType string,
	alerts []domain.Alert,
) error {
	return errors.New("not implemented")
}

// GetAlerts implements AlertCache.
func (r *RedisAlertCache) GetAlerts(
	ctx context.Context,
	alertType string,
) ([]domain.Alert, error) {
	return []domain.Alert{}, errors.New("not implemented")
}

func NewRedisAlertCache(cmd redis.Cmdable) AlertCache {
	return &RedisAlertCache{
		cmd: cmd,
	}
}
