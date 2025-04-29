package session

import (
	"context"
	"fmt"
	"iotstarter/internal/model"
	"iotstarter/internal/security"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(ctx context.Context, db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (s *PostgresRepo) Create(ctx context.Context, userId model.UserId) (*model.Session, error) {
	err := s.Clear(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to clear sessions for user %d: %w", userId, err)
	}

	sql := `
		INSERT INTO sessions (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, token, created_at, expires_at
	`

	seshToken := model.SessionToken(security.GenerateUUID())
	sesh := model.Session{
		UserId:    userId,
		Token:     seshToken,
		ExpiresAt: time.Now().UTC().Add(3 * time.Hour),
	}

	row := s.db.QueryRow(ctx, sql, sesh.UserId, sesh.Token, sesh.ExpiresAt)
	if err := row.Scan(&sesh.ID, &sesh.UserId, &sesh.Token, &sesh.CreatedAt, &sesh.ExpiresAt); err != nil {
		return nil, fmt.Errorf("failed to create sesh %v: %w", sesh, err)
	}

	return &sesh, nil
}

func (s *PostgresRepo) Get(ctx context.Context, token model.SessionToken) (*model.Session, error) {
	sql := `
		SELECT * 
		FROM sessions
		WHERE token = $1
	`

	session := model.Session{}
	row := s.db.QueryRow(ctx, sql, token)
	if err := row.Scan(&session.ID, &session.UserId, &session.Token, &session.CreatedAt, &session.ExpiresAt); err != nil {
		return nil, fmt.Errorf("failed to retrieve session %v: %w", session, err)
	}
	return &session, nil
}

func (s *PostgresRepo) GetUserIdFromToken(ctx context.Context, token model.SessionToken) (model.UserId, error) {
	sql := `
		SELECT users.id
		FROM users
		INNER JOIN sessions ON users.id = sessions.user_id
		WHERE sessions.token = $1
	`

	var userId model.UserId

	row := s.db.QueryRow(ctx, sql, token)
	if err := row.Scan(&userId); err != nil {
		return userId, fmt.Errorf("failed to retrieve user id from token %v: %w", userId, err)
	}
	return userId, nil
}

func (s *PostgresRepo) Clear(ctx context.Context, userId model.UserId) error {
	sql := `
		DELETE FROM sessions
		WHERE user_id = $1
	`
	_, err := s.db.Exec(ctx, sql, userId)
	if err != nil {
		return fmt.Errorf("failed to delete existing sessions for user %d: %w", userId, err)
	}
	return nil
}
