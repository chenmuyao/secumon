package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/chenmuyao/generique/gslice"
	"github.com/chenmuyao/secumon/internal/domain"
	"github.com/chenmuyao/secumon/internal/repository/cache"
	"github.com/chenmuyao/secumon/internal/repository/dao"
)

type AlertRepo interface {
	GetAlerts(ctx context.Context, alertType string, limit int) ([]domain.Alert, error)
}

type CachedAlertRepo struct {
	logDAO          dao.LogDAO
	alertCache      cache.AlertCache
	timeFormat      string
	defaultPageSize int
}

// GetAlerts implements AlertRepo.
func (a *CachedAlertRepo) GetAlerts(
	ctx context.Context,
	alertType string,
	limit int,
) ([]domain.Alert, error) {
	var res []domain.Alert
	if limit < a.defaultPageSize {
		// check the cache
		res, err := a.alertCache.GetAlerts(ctx, alertType)
		if err == nil {
			return res[:limit], nil
		}
		// NOTE: check here if there is a redis error, and we can decide if we
		// still check the DB or just return with error to not stress the DB.
	}

	// NOTE: if the limit get over the defaultPageSize, it must be a special request

	// no cache, compute from the DB
	daoRes, err := a.logDAO.FindAlerts(ctx, alertType, limit)
	if err != nil {
		return []domain.Alert{}, err
	}

	res = gslice.Map(daoRes, func(id int, src dao.SecurityEvent) domain.Alert {
		return domain.Alert{
			Type:      src.Type,
			Timestamp: src.CreatedAt.Format(a.timeFormat),
			ClientIP:  src.ClientIP,
			Details:   "details todo", // TODO: switch on the type and get a proper detail message
		}
	})

	// async write the pageSize of data back to the cache
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		var toCache []domain.Alert
		if len(res) < a.defaultPageSize {
			toCache = make([]domain.Alert, a.defaultPageSize)
			copy(toCache, res)
		} else {
			toCache = res[:a.defaultPageSize]
		}

		er := a.alertCache.SetAlerts(ctx, alertType, toCache)
		if er != nil {
			slog.Error("redis write error", slog.Any("er", er))
		}
	}()

	return res, nil
}

func NewAlertRepo(logDAO dao.LogDAO, alertCache cache.AlertCache) AlertRepo {
	return &CachedAlertRepo{
		logDAO:          logDAO,
		alertCache:      alertCache,
		timeFormat:      "2006-01-02T03:04:05Z",
		defaultPageSize: 10,
	}
}
