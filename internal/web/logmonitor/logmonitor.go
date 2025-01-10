package logmonitor

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	timeFormat string
}

func NewLogHandler() *LogHandler {
	return &LogHandler{
		timeFormat: "2006-01-02T03:04:05Z",
	}
}

func (l *LogHandler) RegisterHandlers(s *gin.Engine) {
	s.POST("/logs", l.AccessLog)
}

func (l *LogHandler) AccessLog(ctx *gin.Context) {
	var req AccessLog
	err := ctx.Bind(&req)
	if err != nil {
		slog.Error("access log input error", slog.Any("err", err))
		return
	}

	err = l.checkDateTime(req.Timestamp)
	if err != nil {
		slog.Error("access log input time error", slog.Any("err", err), slog.Any("req", req))
		ctx.Status(http.StatusBadRequest)
		return
	}

	slog.Debug("access log", slog.Any("al", req))

	ctx.JSON(http.StatusOK, ResultLogOK)
}

func (l *LogHandler) checkDateTime(dt string) error {
	_, err := time.Parse(l.timeFormat, dt)
	return err
}
