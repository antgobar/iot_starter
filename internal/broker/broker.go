package broker

import (
	"context"
	"iotstarter/internal/model"
	"iotstarter/internal/typing"
)

type Broker interface {
	Publish(ctx context.Context, subject string, msg *model.Measurement) error
	Subscribe(ctx context.Context, subject string, handler typing.MeasurementHandler) error
}
