package device

import (
	"context"
	"errors"
	"fmt"
	"iotstarter/internal/model"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(ctx context.Context, db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (s *PostgresRepo) Create(ctx context.Context, device *model.Device) (*model.Device, error) {
	sql := `
        INSERT INTO devices (user_id, location, api_key)
        VALUES ($1, $2, $3)
        RETURNING id, user_id, location, created_at, api_key
    `

	row := s.db.QueryRow(ctx, sql, device.UserId, device.Location, device.ApiKey)
	if err := row.Scan(&device.ID, &device.UserId, &device.Location, &device.CreatedAt, &device.ApiKey); err != nil {
		return nil, fmt.Errorf("failed to register device %v: %w", device, err)
	}
	return device, nil
}

func (s *PostgresRepo) Reauth(ctx context.Context, device *model.Device) (*model.Device, error) {
	sql := `
		UPDATE devices
		SET api_key = $1
		WHERE id = $2 AND user_id = $3
		RETURNING id, user_id, location, created_at, api_key
	`

	row := s.db.QueryRow(ctx, sql, device.ApiKey, device.ID, device.ApiKey)
	storedDevice := model.Device{}
	if err := row.Scan(&storedDevice.ID, &storedDevice.UserId, &storedDevice.Location, &storedDevice.CreatedAt, &storedDevice.ApiKey); err != nil {
		return nil, fmt.Errorf("failed to register device %v: %w", device, err)
	}
	return &storedDevice, nil

}

func (s *PostgresRepo) List(ctx context.Context, userId model.UserId) ([]*model.Device, error) {
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

	return devices, nil
}

func (s *PostgresRepo) GetUserDeviceById(ctx context.Context, userId model.UserId, deviceId model.DeviceId) (*model.Device, error) {
	sql := `
		SELECT id, location, created_at, api_key
		FROM devices 
		WHERE id = $1 AND user_id = $2
		`
	device := model.Device{}

	row := s.db.QueryRow(ctx, sql, deviceId, userId)
	if err := row.Scan(&device.ID, &device.Location, &device.CreatedAt, &device.ApiKey); err != nil {
		if isNoRowsFoundError(err) {
			return nil, ErrDeviceNotFound
		}
		return nil, fmt.Errorf("failed to retrieve device id %v: %w", deviceId, err)
	}
	return &device, nil
}

func (s *PostgresRepo) GetById(ctx context.Context, deviceId model.DeviceId) (*model.Device, error) {
	sql := `
		SELECT id, location, created_at, api_key
		FROM devices 
		WHERE id = $1
		`
	device := model.Device{}

	row := s.db.QueryRow(ctx, sql, deviceId)
	if err := row.Scan(&device.ID, &device.Location, &device.CreatedAt, &device.ApiKey); err != nil {
		if isNoRowsFoundError(err) {
			return nil, ErrDeviceNotFound
		}
		return nil, fmt.Errorf("failed to retrieve device id %v: %w", deviceId, err)
	}
	return &device, nil
}

var ErrDeviceNotFound = errors.New("device not found")

func isNoRowsFoundError(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
