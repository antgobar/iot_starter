package measurement

import (
	"context"
	"iotstarter/internal/model"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) GetMeasurements(ctx context.Context, u model.UserId, d model.DeviceId, start, end time.Time) ([]*model.Measurement, error) {
	return s.repo.GetDeviceMeasurements(ctx, u, d, start, end)
}
