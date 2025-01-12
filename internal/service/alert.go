package service

import (
	"context"

	"github.com/chenmuyao/secumon/internal/domain"
	"github.com/chenmuyao/secumon/internal/repository"
)

type AlertService interface {
	GetAlerts(ctx context.Context, alertType string, limit int) ([]domain.Alert, error)
}

type alertService struct {
	repo repository.AlertRepo
}

// GetAlerts implements AlertService.
func (a *alertService) GetAlerts(
	ctx context.Context,
	alertType string,
	limit int,
) ([]domain.Alert, error) {
	return a.repo.GetAlerts(ctx, alertType, limit)
}

func NewAlertService(repo repository.AlertRepo) AlertService {
	return &alertService{
		repo: repo,
	}
}
