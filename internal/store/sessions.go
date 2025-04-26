package store

import (
	"context"
	"fmt"
	"iotstarter/internal/auth"
	"iotstarter/internal/model"
	"time"
)

func (s *PostgresStore) CreateUserSession(ctx context.Context, userId int) (*model.Session, error) {
	deleteSessionSql := `
		DELETE FROM sessions
		WHERE user_id = $1
	`
	_, err := s.db.Exec(ctx, deleteSessionSql, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete existing sessions for user %d: %w", userId, err)
	}

	sql := `
		INSERT INTO sessions (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, token, created_at, expires_at
	`

	sesh := model.Session{
		UserId:    userId,
		Token:     auth.GenerateUUID(),
		ExpiresAt: time.Now().UTC().Add(3 * time.Hour),
	}

	row := s.db.QueryRow(ctx, sql, sesh.UserId, sesh.Token, sesh.ExpiresAt)
	if err := row.Scan(&sesh.ID, &sesh.UserId, &sesh.Token, &sesh.CreatedAt, &sesh.ExpiresAt); err != nil {
		return nil, fmt.Errorf("failed to create sesh %v: %w", sesh, err)
	}

	return &sesh, nil
}

func (s *PostgresStore) GetUserFromToken(ctx context.Context, token string) (*model.User, error) {
	sql := `
		SELECT users.id, users.username, users.created_at, users.active
		FROM users
		INNER JOIN sessions ON users.id = sessions.user_id
		WHERE sessions.token = $1 AND sessions.expires_at > NOW()
	`

	user := model.User{}

	row := s.db.QueryRow(ctx, sql, token)
	if err := row.Scan(&user.ID, &user.Username, &user.CreatedAt, &user.Active); err != nil {
		return nil, fmt.Errorf("failed to retrieve user from token %v: %w", user, err)
	}
	return &model.User{}, nil
}
