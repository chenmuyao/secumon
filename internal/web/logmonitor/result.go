package logmonitor

type Result struct {
	status  string `json:"status"`
	message string `json:"message"`
}

var ResultLogOK = Result{
	status:  "sucess",
	message: "Log received and queued",
}
