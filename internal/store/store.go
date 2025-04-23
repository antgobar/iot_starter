package store

import (
	"context"
	"iotstarter/internal/model"
	"time"
)

type Store interface {
	RegisterDevice(ctx context.Context, location string) (*model.Device, error)
	GetDevices(ctx context.Context) ([]model.Device, error)
	GetDeviceById(ctx context.Context, deviceId int) (*model.Device, error)
	SaveMeasurement(ctx context.Context, m *model.Measurement) error
	GetDeviceMeasurements(ctx context.Context, deviceId int, start, end time.Time) ([]model.Measurement, error)
}
