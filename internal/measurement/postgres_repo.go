package measurement

import (
	"context"
	"fmt"
	"iotstarter/internal/model"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(ctx context.Context, db *pgxpool.Pool) *postgresRepo {
	return &postgresRepo{db: db}
}

func (s *postgresRepo) Create(ctx context.Context, m *model.Measurement) error {
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
		return fmt.Errorf("failed to insert measurement %v: %w", m, err)
	}
	return nil
}

func (s *postgresRepo) GetDeviceMeasurements(ctx context.Context, userId model.UserId, deviceId model.DeviceId, start, end time.Time) ([]*model.Measurement, error) {
	sql := `
		SELECT measurements.id, measurements.device_id, measurements.name, measurements.value, measurements.unit, measurements.timestamp 
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
