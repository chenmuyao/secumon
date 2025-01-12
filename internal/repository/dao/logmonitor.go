package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LogDAO interface {
	UpsertSecurityEvent(ctx context.Context, secuEvt SecurityEvent) error
}

type SecurityEvent struct {
	ID        uint `gorm:"primarykey,autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"index"`

	Type      string    `gorm:"uniqueIndex:ip_type_ts,length=128"`
	Timestamp time.Time `gorm:"uniqueIndex:ip_type_ts"`
	ClientIP  string    `gorm:"uniqueIndex:ip_type_ts,length=128"`
	Attacks   int
	Details   string
}

type GORMLogDAO struct {
	db *gorm.DB
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
