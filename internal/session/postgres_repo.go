package session

import (
	"context"
	"fmt"
	"iotstarter/internal/model"
	"iotstarter/internal/security"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(ctx context.Context, db *pgxpool.Pool) *postgresRepo {
	return &postgresRepo{db: db}
}

func (s *postgresRepo) Create(ctx context.Context, userId model.UserId) (*model.Session, error) {
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

func (s *postgresRepo) GetUserFromToken(ctx context.Context, token model.SessionToken) (*model.User, error) {
	sql := `
		SELECT users.id, users.username, users.created_at, users.active
		FROM users
		INNER JOIN sessions ON users.id = sessions.user_id
		WHERE sessions.token = $1
	`

	var user model.User

	row := s.db.QueryRow(ctx, sql, token)
	if err := row.Scan(&user.ID, &user.Username, &user.CreatedAt, &user.Active); err != nil {
		return &user, noUserSessionErr(err)
	}
	return &user, nil
}

func (s *postgresRepo) Clear(ctx context.Context, userId model.UserId) error {
	sql := `
		DELETE FROM sessions
		WHERE user_id = $1
	`
	_, err := s.db.Exec(ctx, sql, userId)
	if err != nil {
		return failedToDeleteUserSessionErr(userId, err)
	}
	return nil
}
