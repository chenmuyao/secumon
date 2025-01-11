package logmonitor

type Result struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var ResultLogOK = Result{
	Status:  "success",
	Message: "Log received and queued",
}

var ResultLogErrPublish = Result{
	Status:  "fail",
	Message: "Log received but not queued",
}
