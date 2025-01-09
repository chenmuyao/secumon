package logmonitor

type AccessLog struct {
	Timestamp  string `json:"timestamp"   binding:"required"`
	ClientIP   string `json:"client_ip"   binding:"ip"`
	Endpoint   string `json:"endpoint"    binding:"required"`
	Method     string `json:"method"      binding:"required"`
	StatusCode int    `json:"status_code" binding:"required"`
}
