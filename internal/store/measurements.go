package store

import (
	"context"
	"fmt"
	"iotstarter/internal/model"
	"log"
	"time"
)

func (s *PostgresStore) SaveMeasurement(ctx context.Context, m *model.Measurement) (*model.Measurement, error) {
	log.Println("reached saved measurement", time.Now(), m)
	sql := `
		INSERT INTO measurements (device_id, name, value, unit, timestamp)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, device_id, name, value, unit, timestamp 
		`

	row := s.db.QueryRow(ctx, sql, m.DeviceId, m.Name, m.Value, m.Unit, m.Timestamp)
	err := row.Scan(
		&m.ID, &m.DeviceId, &m.Name, &m.Value, &m.Unit, &m.Timestamp,
	)
	if err != nil {
		return m, fmt.Errorf("failed to insert measurement %v: %w", m, err)
	}
	return m, nil
}

func (s *PostgresStore) GetDeviceMeasurements(ctx context.Context, deviceId int, start, end time.Time) ([]model.Measurement, error) {
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

	var measurements = make([]model.Measurement, 0)
	for rows.Next() {
		var m model.Measurement
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
