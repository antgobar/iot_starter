package consumer

import (
	"context"
	"iotstarter/internal/typing"
)

type Subscriber interface {
	Subscribe(ctx context.Context, subject string, handler typing.MeasurementHandler) error
}
