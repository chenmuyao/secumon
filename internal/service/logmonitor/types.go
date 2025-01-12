package logmonitor

import (
	"context"

	"github.com/chenmuyao/secumon/internal/domain"
)

type Detector interface {
	Detect(ctx context.Context, log domain.AccessLog) error
}
