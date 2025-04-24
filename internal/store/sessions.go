package store

import (
	"context"
	"fmt"
	"iotstarter/internal/model"
	"time"
)

func (s *PostgresStore) CreateUserSession(ctx context.Context, userId int) (*model.Session, error) {
	sql := `
		INSERT INTO sessions user_id, token, expires_at
		VALUES ($1, $2, $3, $4)
		RETURNING (id, user_id, token, created_at, expires_at)
	`

	sesh := model.Session{
		UserId:    userId,
		Token:     "foo",
		ExpiresAt: time.Now().UTC().Add(3 * time.Hour),
	}

	row := s.db.QueryRow(ctx, sql, sesh.UserId, sesh.Token, sesh.ExpiresAt)
	if err := row.Scan(&sesh.ID, &sesh.UserId, &sesh.Token, &sesh.CreatedAt, &sesh.ExpiresAt); err != nil {
		return nil, fmt.Errorf("failed to create sesh %v: %w", sesh, err)
	}

	return &sesh, nil
}
