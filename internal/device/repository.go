package device

import (
	"context"
	"iotstarter/internal/model"
)

type Repository interface {
	Create(ctx context.Context, device *model.Device) (*model.Device, error)
	Reauth(ctx context.Context, device *model.Device) (*model.Device, error)
	List(ctx context.Context, userId model.UserId) ([]*model.Device, error)
	GetByUserById(ctx context.Context, userId model.UserId, deviceId model.DeviceId) (*model.Device, error)
	GetById(ctx context.Context, deviceId model.DeviceId) (*model.Device, error)
	DeleteByUserIdById(ctx context.Context, userId model.UserId, deviceId model.DeviceId) error
}
