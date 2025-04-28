package gateway

import (
	"context"
	"errors"
	"iotstarter/internal/device"
	"iotstarter/internal/model"
)

type Service struct {
	pub     Publisher
	devices device.Repository
}

func NewService(p Publisher) *Service {
	return &Service{pub: p}
}

func (s *Service) Publish(ctx context.Context, subject string, measurement *model.Measurement) error {
	return s.pub.Publish(ctx, subject, measurement)
}

func (s *Service) CheckDeviceIsAuthed(ctx context.Context, deviceId model.DeviceId, apiKey model.ApiKey) error {
	measurementDevice, err := s.devices.GetById(ctx, deviceId)
	if err != nil {
		return device.ErrDeviceNotFound
	}
	if measurementDevice.ApiKey != apiKey {
		return errors.New("invalid API key for device id: " + string(deviceId))
	}
	return nil
}
