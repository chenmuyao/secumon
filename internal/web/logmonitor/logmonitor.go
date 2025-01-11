package logmonitor

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/chenmuyao/secumon/internal/domain"
	"github.com/chenmuyao/secumon/internal/event/monitor"
	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	timeFormat string
	publisher  monitor.LogMonitorPublisher
}

func NewLogHandler(publisher monitor.LogMonitorPublisher) *LogHandler {
	return &LogHandler{
		timeFormat: "2006-01-02T03:04:05Z",
		publisher:  publisher,
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

	t, err := time.Parse(l.timeFormat, req.Timestamp)
	if err != nil {
		slog.Error("access log input time error", slog.Any("err", err), slog.Any("req", req))
		ctx.Status(http.StatusBadRequest)
		return
	}

	slog.Debug("access log", slog.Any("al", req))

	err = l.publisher.Publish(ctx, domain.AccessLog{
		Timestamp:  t,
		ClientIP:   req.ClientIP,
		Endpoint:   req.Endpoint,
		Method:     req.Method,
		StatusCode: req.StatusCode,
	})
	if err != nil {
		slog.Error("publish log error", slog.Any("err", err), slog.Any("req", req))
		ctx.JSON(http.StatusInternalServerError, ResultLogErrPublish)
		return
	}

	ctx.JSON(http.StatusOK, ResultLogOK)
}
