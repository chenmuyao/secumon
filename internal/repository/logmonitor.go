package repository

import (
	"context"

	"github.com/chenmuyao/secumon/internal/domain"
)

type LogRepo interface {
	UpsertSecurityEvent(ctx context.Context, secuEvt domain.SecurityEvent) error
}

type logRepo struct{}

// UpsertSecurityEvent implements LogRepo.
func (l *logRepo) UpsertSecurityEvent(ctx context.Context, secuEvt domain.SecurityEvent) error {
	return nil
}

func NewLogRepo() LogRepo {
	return &logRepo{}
}
