package gateway

import (
	"context"
	"iotstarter/internal/model"
)

type Publisher interface {
	Publish(ctx context.Context, subject string, measurement *model.Measurement) error
}
