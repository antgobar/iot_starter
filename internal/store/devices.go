package store

import (
	"context"
	"fmt"
	"iotstarter/internal/auth"
	"iotstarter/internal/model"
	"log"
)

func (s *PostgresStore) RegisterDevice(ctx context.Context, userId int, location string) (*model.Device, error) {
	sql := `
        INSERT INTO devices (user_id, location, api_key)
        VALUES ($1, $2, $3)
        RETURNING id, user_id, location, created_at, api_key
    `
	apiKey := model.ApiKey(auth.GenerateUUID())
	userIdTyped := model.UserId(userId)

	device := model.Device{
		UserId:   userIdTyped,
		Location: location,
		ApiKey:   apiKey,
	}

	row := s.db.QueryRow(ctx, sql, device.UserId, device.Location, device.ApiKey)
	if err := row.Scan(&device.ID, &device.UserId, &device.Location, &device.CreatedAt, &device.ApiKey); err != nil {
		return nil, fmt.Errorf("failed to register device %v: %w", device, err)
	}
	return &device, nil
}

func (s *PostgresStore) ReauthDevice(ctx context.Context, userId int, deviceId int) (*model.Device, error) {
	sql := `
		UPDATE devices
		SET api_key = $1
		WHERE id = $2 AND user_id = $3
		RETURNING id, user_id, location, created_at, api_key
	`
	deviceIdTyped := model.DeviceId(deviceId)
	userIdTyped := model.UserId(userId)
	apiKey := model.ApiKey(auth.GenerateUUID())
	device := model.Device{
		ID:     deviceIdTyped,
		UserId: userIdTyped,
		ApiKey: apiKey,
	}

	row := s.db.QueryRow(ctx, sql, device.ApiKey, device.ID, device.UserId)
	if err := row.Scan(&device.ID, &device.UserId, &device.Location, &device.CreatedAt, &device.ApiKey); err != nil {
		return nil, fmt.Errorf("failed to register device %v: %w", device, err)
	}
	return &device, nil

}

func (s *PostgresStore) GetDevices(ctx context.Context, userId int) ([]*model.Device, error) {
	sql := `
		SELECT id, user_id, location, created_at, api_key 
		FROM devices
		WHERE user_id = $1
	`
	rows, err := s.db.Query(ctx, sql, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query devices: %w", err)
	}
	defer rows.Close()

	devices := make([]*model.Device, 0)
	for rows.Next() {
		var d model.Device
		if err := rows.Scan(&d.ID, &d.UserId, &d.Location, &d.CreatedAt, &d.ApiKey); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		devices = append(devices, &d)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	log.Println("GOT devices for user", userId, "devices", len(devices))
	return devices, nil
}

func (s *PostgresStore) GetDeviceById(ctx context.Context, deviceId int) (*model.Device, error) {
	sql := `
		SELECT id, location, created_at, api_key
		FROM devices 
		WHERE id = $1
		`
	var device model.Device
	row := s.db.QueryRow(ctx, sql, deviceId)
	if err := row.Scan(&device.ID, &device.Location, &device.CreatedAt, &device.ApiKey); err != nil {
		if isNoRowsFoundError(err) {
			return nil, ErrDeviceNotFound
		}
		return nil, fmt.Errorf("failed to retrieve device id %v: %w", deviceId, err)
	}
	return &device, nil
}
