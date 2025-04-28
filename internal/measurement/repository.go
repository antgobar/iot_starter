package measurement

import (
	"context"
	"iotstarter/internal/model"
	"time"
)

type Repository interface {
	Create(ctx context.Context, measurement *model.Measurement) error
	GetDeviceMeasurements(ctx context.Context, userId model.UserId, deviceId model.DeviceId, start, end time.Time) ([]*model.Measurement, error)
}
