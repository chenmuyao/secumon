package logmonitor

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogHandler struct{}

func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

func (l *LogHandler) RegisterHandlers(s *gin.Engine) {
	s.POST("/logs", l.AccessLog)
}

func (l *LogHandler) AccessLog(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
