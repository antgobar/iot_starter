package consumer

import (
	"context"
	"iotstarter/internal/model"
)

type MeasurementHandler func(msg *model.Measurement)

type Subscriber interface {
	Subscribe(ctx context.Context, subject string, handler MeasurementHandler) error
}
