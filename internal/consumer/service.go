package consumer

import (
	"context"
	"iotstarter/internal/measurement"
	"iotstarter/internal/model"
)

type Service struct {
	sub          Subscriber
	measurements measurement.Repository
	subject      string
}

func NewService(sub Subscriber, measurements measurement.Repository, subject string) *Service {
	return &Service{
		sub:          sub,
		measurements: measurements,
		subject:      subject,
	}
}

func (s *Service) StoreMeasurement(ctx context.Context, measurement *model.Measurement) error {
	return s.measurements.Create(ctx, measurement)
}
