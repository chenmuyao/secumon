package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LogDAO interface {
	UpsertSecurityEvent(ctx context.Context, secuEvt SecurityEvent) error
	FindAlerts(
		ctx context.Context,
		alertType string,
		limit int,
	) ([]SecurityEvent, error)
}

type SecurityEvent struct {
	ID        uint `gorm:"primarykey,autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"index"`

	Type      string    `gorm:"uniqueIndex:ip_type_ts,length=128"`
	ClientIP  string    `gorm:"uniqueIndex:ip_type_ts,length=128"`
	Timestamp time.Time `gorm:"uniqueIndex:ip_type_ts"`
	Attacks   int
	Details   string
}

type GORMLogDAO struct {
	db *gorm.DB
}

// FindAlerts implements LogDAO.
func (g *GORMLogDAO) FindAlerts(
	ctx context.Context,
	alertType string,
	limit int,
) ([]SecurityEvent, error) {
	var res []SecurityEvent
	var err error
	if alertType == "" {
		// don't select on the alert type
		err = g.db.WithContext(ctx).Limit(limit).Order("updated_at DESC").Find(&res).Error
		return res, err
	}
	err = g.db.WithContext(ctx).
		Where("type = ?", alertType).
		Limit(limit).
		Order("updated_at DESC").
		Find(&res).
		Error
	return res, err
}

// UpsertSecurityEvent implements LogDAO.
func (g *GORMLogDAO) UpsertSecurityEvent(ctx context.Context, secuEvt SecurityEvent) error {
	now := time.Now()
	return g.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "client_ip"}, {Name: "type"}, {Name: "timestamp"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"updated_at": now,
			"attacks":    gorm.Expr("security_events.attacks + 1"),
		}),
	}).Create(&secuEvt).Error
}

func NewLogDAO(db *gorm.DB) LogDAO {
	return &GORMLogDAO{
		db: db,
	}
}
