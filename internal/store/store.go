package store

import (
	"context"
	"fmt"
	"iotstarter/internal/measurement"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(ctx context.Context, url string) (*Store, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	store := Store{db: pool}
	store.setUpTables(ctx)
	return &store, nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) RegisterDevice(ctx context.Context, location string) error {
	newDevice := measurement.Device{
		Location:  location,
		CreatedAt: time.Now().UTC(),
	}
	sql := `
        INSERT INTO devices (location, created_at)
        VALUES ($1, $2)
        RETURNING id
    `
	var deviceId int
	if err := s.db.QueryRow(ctx, sql, newDevice.Location, newDevice.CreatedAt).Scan(&deviceId); err != nil {
		return fmt.Errorf("failed to insert device %v: %w", newDevice, err)
	}
	return nil
}

func (s *Store) GetDevices(ctx context.Context) ([]measurement.Device, error) {
	sql := `SELECT id, location, created_at FROM devices`
	rows, err := s.db.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to query devices: %w", err)
	}
	defer rows.Close()

	var devices []measurement.Device
	for rows.Next() {
		var d measurement.Device
		if err := rows.Scan(&d.ID, &d.Location, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		devices = append(devices, d)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return devices, nil
}

func (s *Store) SaveMeasurement(ctx context.Context, m measurement.Measurement) error {
	sql := `
		INSERT INTO measurements (device_id, name, value, unit, timestamp)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, device_id, name, value, unit, timestamp 
		`
	var storedM measurement.Measurement
	row := s.db.QueryRow(ctx, sql, m.DeviceId, m.Name, m.Value, m.Unit, m.Timestamp)
	err := row.Scan(
		&storedM.ID, &storedM.DeviceId, &storedM.Name, &storedM.Value, &storedM.Unit, &storedM.Timestamp,
	)
	if err != nil {
		return fmt.Errorf("failed to insert measurement %v: %w", storedM, err)
	}
	return nil
}

func (s *Store) GetDeviceMeasurements(ctx context.Context, deviceId int, start, end time.Time) ([]measurement.Measurement, error) {
	sql := `
		SELECT * FROM measurements
		WHERE device_id = $1
		AND timestamp BETWEEN $2 AND $3
	`
	rows, err := s.db.Query(ctx, sql, deviceId, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve measurement: %w", err)
	}
	defer rows.Close()

	var measurements = make([]measurement.Measurement, 0)
	for rows.Next() {
		var m measurement.Measurement
		if err := rows.Scan(&m.ID, &m.DeviceId, &m.Name, &m.Value, &m.Unit, &m.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		measurements = append(measurements, m)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return measurements, nil

}

func (s *Store) setUpTables(ctx context.Context) error {
	err := s.execSQLFile(ctx, "sql/schema.sql")
	if err != nil {
		return err
	}
	log.Println("Tables created")
	return nil
}

func (s *Store) execSQLFile(ctx context.Context, path string) error {
	sqlBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read SQL file: %w", err)
	}

	sql := string(sqlBytes)
	_, err = s.db.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("failed to execute SQL: %w", err)
	}

	return nil
}
