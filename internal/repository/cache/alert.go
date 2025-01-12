package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/chenmuyao/secumon/internal/domain"
	"github.com/redis/go-redis/v9"
)

type AlertCache interface {
	GetAlerts(ctx context.Context, alertType string) ([]domain.Alert, error)
	SetAlerts(ctx context.Context, alertType string, alerts []domain.Alert) error
	DeleteAlerts(ctx context.Context, alertType string) error
}

type RedisAlertCache struct {
	cmd        redis.Cmdable
	expiryTime time.Duration
}

// DeleteAlerts implements AlertCache.
func (r *RedisAlertCache) DeleteAlerts(ctx context.Context, alertType string) error {
	return r.cmd.Del(ctx, r.Key(alertType)).Err()
}

// SetAlerts implements AlertCache.
func (r *RedisAlertCache) SetAlerts(
	ctx context.Context,
	alertType string,
	alerts []domain.Alert,
) error {
	val, err := json.Marshal(alerts)
	if err != nil {
		slog.Error("json marshall error", slog.Any("err", err), slog.Any("alerts", alerts))
		return err
	}
	return r.cmd.Set(ctx, r.Key(alertType), val, r.expiryTime).Err()
}

// GetAlerts implements AlertCache.
func (r *RedisAlertCache) GetAlerts(
	ctx context.Context,
	alertType string,
) ([]domain.Alert, error) {
	resStr, err := r.cmd.Get(ctx, r.Key(alertType)).Result()
	if err != nil {
		return []domain.Alert{}, err
	}
	var res []domain.Alert
	err = json.Unmarshal([]byte(resStr), &res)
	if err != nil {
		slog.Error("json unmarshall error", slog.Any("err", err), slog.Any("resStr", resStr))
		return []domain.Alert{}, err
	}
	return res, err
}

func (r *RedisAlertCache) Key(alertType string) string {
	if alertType == "" {
		alertType = "all"
	}
	return fmt.Sprintf("alert:%s", alertType)
}

func NewRedisAlertCache(cmd redis.Cmdable) AlertCache {
	return &RedisAlertCache{
		cmd:        cmd,
		expiryTime: 15 * time.Minute,
	}
}
