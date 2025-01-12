package domain

import "time"

type SecurityEvent struct {
	Type      string
	Timestamp time.Time
	ClientIP  string
	Details   string
}
