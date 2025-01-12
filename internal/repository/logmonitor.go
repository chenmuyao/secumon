package repository

import (
	"context"
	"time"

	"github.com/chenmuyao/secumon/internal/domain"
	"github.com/chenmuyao/secumon/internal/repository/dao"
)

type LogRepo interface {
	UpsertSecurityEvent(ctx context.Context, secuEvt domain.SecurityEvent) error
}

type logRepo struct {
	logDAO dao.LogDAO
}

// UpsertSecurityEvent implements LogRepo.
func (l *logRepo) UpsertSecurityEvent(ctx context.Context, secuEvt domain.SecurityEvent) error {
	// NOTE: in case attacks arrives in the same minute, instead of log every attack,
	// we log only one per minute.
	now := time.Now()

	daoEvt := dao.SecurityEvent{
		CreatedAt: now,
		UpdatedAt: now,
		Type:      secuEvt.Type,
		Timestamp: secuEvt.Timestamp.Truncate(time.Minute),
		ClientIP:  secuEvt.ClientIP,
		Attacks:   1,
		Details:   secuEvt.Details,
	}

	return l.logDAO.UpsertSecurityEvent(ctx, daoEvt)
}

func NewLogRepo(dao dao.LogDAO) LogRepo {
	return &logRepo{
		logDAO: dao,
	}
}
