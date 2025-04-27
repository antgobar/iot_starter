package device

import (
	"context"
	"iotstarter/internal/auth"
	"iotstarter/internal/model"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func generateApiKey() model.ApiKey {
	return model.ApiKey(auth.GenerateUUID())
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

func (s *Service) GetById(ctx context.Context, u model.UserId, d model.DeviceId) (*model.Device, error) {
	return s.repo.GetById(ctx, u, d)
}

func (s *Service) GetMeasurements(ctx context.Context, u model.UserId, d model.DeviceId, start, end time.Time) ([]*model.Measurement, error) {
	return s.repo.GetMeasurements(ctx, u, d, start, end)
}
