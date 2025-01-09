package logmonitor

type AccessLog struct {
	Timestamp  string `json:"timestamp"`
	ClientIP   string `json:"client_ip"`
	Endpoint   string `json:"endpoint"`
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
}
