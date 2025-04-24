package store

import (
	"context"
	"iotstarter/internal/model"
	"time"
)

type Store interface {
	UserStore
	DeviceStore
	MeasurementStore
	SessionStore
}

type UserStore interface {
	RegisterUser(ctx context.Context, userName string, password string) (*model.User, error)
}

type DeviceStore interface {
	RegisterDevice(ctx context.Context, userId int, location string) (*model.Device, error)
	ReauthDevice(ctx context.Context, userId int, deviceId int) (*model.Device, error)
	GetDevices(ctx context.Context) ([]model.Device, error)
	GetDeviceById(ctx context.Context, deviceId int) (*model.Device, error)
}

type MeasurementStore interface {
	SaveMeasurement(ctx context.Context, m *model.Measurement) (*model.Measurement, error)
	GetDeviceMeasurements(ctx context.Context, deviceId int, start, end time.Time) ([]model.Measurement, error)
}

type SessionStore interface {
	CreateUserSession(ctx context.Context, userId int) (*model.Session, error)
}
