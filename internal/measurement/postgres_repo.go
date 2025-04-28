package measurement

import (
	"context"
	"fmt"
	"iotstarter/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(ctx context.Context, db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (s *PostgresRepo) Create(ctx context.Context, m *model.Measurement) error {
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
