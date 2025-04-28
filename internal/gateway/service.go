package gateway

import (
	"context"
	"fmt"
	"iotstarter/internal/device"
	"iotstarter/internal/model"
	"strconv"
)

type Service struct {
	pub     Publisher
	devices device.Repository
}

func NewService(p Publisher, d device.Repository) *Service {
	return &Service{pub: p, devices: d}
}

func (s *Service) Publish(ctx context.Context, subject string, measurement *model.Measurement) error {
	return s.pub.Publish(ctx, subject, measurement)
}

func (s *Service) CheckDeviceIsAuthed(ctx context.Context, deviceId model.DeviceId, apiKey model.ApiKey) error {
	measurementDevice, err := s.devices.GetById(ctx, deviceId)
	if err != nil {
		return device.ErrDeviceNotFound
	}

	deviceIDstr := strconv.Itoa(int(deviceId))
	if measurementDevice.ApiKey != apiKey {
		return fmt.Errorf("invalid API key for device id: %s", deviceIDstr)
	}
	return nil
}
