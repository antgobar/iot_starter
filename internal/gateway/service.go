package gateway

import (
	"context"
	"iotstarter/internal/model"
)

type Service struct {
	pub Publisher
}

func NewService(p Publisher) *Service {
	return &Service{pub: p}
}

func (s *Service) Publish(ctx context.Context, subject string, measurement *model.Measurement) error {
	return s.pub.Publish(ctx, subject, measurement)
}
