package device

import (
	"context"
	"fmt"
	"iotstarter/internal/model"
	"iotstarter/internal/security"
	"strconv"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func generateApiKey() model.ApiKey {
	return model.ApiKey(security.GenerateUUID())
}

func (s *Service) Register(ctx context.Context, userId model.UserId, location string) (*model.Device, error) {
	device := model.Device{
		UserId:   userId,
		Location: location,
		ApiKey:   generateApiKey(),
	}
	return s.repo.Create(ctx, &device)
}

func (s *Service) Reauth(ctx context.Context, userId model.UserId, deviceId model.DeviceId) (*model.Device, error) {
	device := model.Device{
		ID:     deviceId,
		UserId: userId,
		ApiKey: generateApiKey(),
	}
	return s.repo.Reauth(ctx, &device)
}

func (s *Service) List(ctx context.Context, u model.UserId) ([]*model.Device, error) {
	return s.repo.List(ctx, u)
}

func (s *Service) GetUserDeviceById(ctx context.Context, u model.UserId, d model.DeviceId) (*model.Device, error) {
	return s.repo.GetUserDeviceById(ctx, u, d)
}

func (s *Service) CheckDeviceToken(ctx context.Context, deviceId model.DeviceId, apiKey model.ApiKey) error {
	measurementDevice, err := s.repo.GetById(ctx, deviceId)
	if err != nil {
		return ErrDeviceNotFound
	}

	deviceIDstr := strconv.Itoa(int(deviceId))
	if measurementDevice.ApiKey != apiKey {
		return fmt.Errorf("invalid API key for device id: %s", deviceIDstr)
	}
	return nil
}
