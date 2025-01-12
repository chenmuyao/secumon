package domain

type Alert struct {
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	ClientIP  string `json:"client_ip"`
	Details   string `json:"details"`
}
