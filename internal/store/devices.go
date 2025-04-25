package store

import (
	"context"
	"errors"
	"fmt"
	"iotstarter/internal/auth"
	"iotstarter/internal/model"

	"github.com/jackc/pgx/v5"
)

func (s *PostgresStore) RegisterDevice(ctx context.Context, userId int, location string) (*model.Device, error) {
	sql := `
        INSERT INTO devices (user_id, location, api_key)
        VALUES ($1, $2, $3)
        RETURNING id, user_id, location, created_at, api_key
    `

	device := model.Device{
		UserId:   userId,
		Location: location,
		ApiKey:   auth.GenerateUUID(),
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
	device := model.Device{
		ID:     deviceId,
		UserId: userId,
		ApiKey: auth.GenerateUUID(),
	}

	row := s.db.QueryRow(ctx, sql, device.ApiKey, device.ID, device.UserId)
	if err := row.Scan(&device.ID, &device.UserId, &device.Location, &device.CreatedAt, &device.ApiKey); err != nil {
		return nil, fmt.Errorf("failed to register device %v: %w", device, err)
	}
	return &device, nil

}

func (s *PostgresStore) GetDevices(ctx context.Context) ([]model.Device, error) {
	sql := `SELECT id, location, created_at, api_key FROM devices`
	rows, err := s.db.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to query devices: %w", err)
	}
	defer rows.Close()

	var devices []model.Device
	for rows.Next() {
		var d model.Device
		if err := rows.Scan(&d.ID, &d.Location, &d.CreatedAt, &d.ApiKey); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		devices = append(devices, d)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDeviceNotFound
		}
		return nil, fmt.Errorf("failed to retrieve device id %v: %w", deviceId, err)
	}
	return &device, nil
}
