package measurement

import (
	"context"
	"iotstarter/internal/model"
)

type Repository interface {
	Create(ctx context.Context, measurement *model.Measurement) error
}
