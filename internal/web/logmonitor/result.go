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

var ResultErrBadRequest = Result{
	Status:  "fail",
	Message: "Bad request, wrong query",
}

var ResultInternal = Result{
	Status:  "fail",
	Message: "failed to get the results",
}
