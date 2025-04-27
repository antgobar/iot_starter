package device

import (
	"context"
	"errors"
	"fmt"
	"iotstarter/internal/model"
	"log"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(ctx context.Context, url string) *PostgresRepo {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		log.Fatalln("Error connecting to PostgresDB", err.Error())
	}
	store := PostgresRepo{db: pool}
	return &store
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

func (s *PostgresRepo) GetById(ctx context.Context, userId model.UserId, deviceId model.DeviceId) (*model.Device, error) {
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

var ErrDeviceNotFound = errors.New("device not found")

func isNoRowsFoundError(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func (s *PostgresRepo) GetMeasurements(ctx context.Context, userId model.UserId, deviceId model.DeviceId, start, end time.Time) ([]*model.Measurement, error) {
	sql := `
		SELECT measurements.id, measurements.user_id, measurements.name, measurements.value, measurements.unit, measurements.timestamp 
		FROM measurements
		INNER JOIN devices ON devices.id = measurements.device_id
		WHERE measurements.device_id = $1
		AND devices.user_id = $2
		AND timestamp BETWEEN $3 AND $4
	`
	rows, err := s.db.Query(ctx, sql, deviceId, userId, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve measurement: %w", err)
	}
	defer rows.Close()

	var measurements = make([]*model.Measurement, 0)
	for rows.Next() {
		var m model.Measurement
		if err := rows.Scan(&m.ID, &m.DeviceId, &m.Name, &m.Value, &m.Unit, &m.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		measurements = append(measurements, &m)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return measurements, nil
}
