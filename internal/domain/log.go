package domain

import "time"

type AccessLog struct {
	Timestamp  time.Time `json:"timestamp"`
	ClientIP   string    `json:"client_ip"`
	Endpoint   string    `json:"endpoint"`
	Method     string    `json:"method"`
	StatusCode int       `json:"status_code"`
}
