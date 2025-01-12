package logmonitor

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/chenmuyao/secumon/internal/domain"
	"github.com/chenmuyao/secumon/internal/event/monitor"
	"github.com/chenmuyao/secumon/internal/service"
	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	timeFormat        string
	publisher         monitor.LogMonitorPublisher
	alertSvc          service.AlertService
	defaultQueryLimit int
}

func NewLogHandler(
	publisher monitor.LogMonitorPublisher,
	alertSvc service.AlertService,
) *LogHandler {
	return &LogHandler{
		timeFormat:        "2006-01-02T03:04:05Z",
		publisher:         publisher,
		alertSvc:          alertSvc,
		defaultQueryLimit: 10,
	}
}

func (l *LogHandler) RegisterHandlers(s *gin.Engine) {
	s.POST("/logs", l.AccessLog)
	s.GET("/alerts", l.Alerts)
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

	// TODO: To improve the performance, we could create a standalone goroutine
	// for the publisher. Here we simply push the log into a queue, and
	// the publisher process publishes a batch at a time, so that we can
	// reuse the channel and other connection ressources.
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

func (l *LogHandler) Alerts(ctx *gin.Context) {
	// Query accepted: type, limit
	// default page size (limit) : 10
	var err error

	limit := l.defaultQueryLimit
	limitStr := ctx.Query("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ResultErrBadRequest)
			return
		}
	}

	alertType := ctx.Query("type")

	alerts, err := l.alertSvc.GetAlerts(ctx, alertType, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResultInternal)
		return
	}

	ctx.JSON(http.StatusOK, alerts)
}
