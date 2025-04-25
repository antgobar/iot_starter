package store

import (
	"context"
	"errors"
	"fmt"
	"iotstarter/internal/auth"
	"iotstarter/internal/model"
)

func (s *PostgresStore) RegisterUser(ctx context.Context, userName string, password string) (*model.User, error) {
	sql := `
		INSERT INTO users (username, hashed_password)
		VALUES ($1, $2)
		Returning id, username, created_at, active
	`

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Username:       userName,
		HashedPassword: hashedPassword,
	}

	row := s.db.QueryRow(ctx, sql, user.Username, user.HashedPassword)
	if err := row.Scan(&user.ID, &user.Username, &user.CreatedAt, &user.Active); err != nil {
		return nil, fmt.Errorf("failed to register user %v: %w", user, err)
	}

	return &user, nil
}

func (s *PostgresStore) GetUserFromCreds(ctx context.Context, userName string, password string) (*model.User, error) {
	sql := `
		SELECT id, username, hashed_password, created_at, active
		FROM users 
		WHERE username = $1 AND active = TRUE
	`

	user := model.User{
		Username: userName,
	}

	row := s.db.QueryRow(ctx, sql, user.Username)
	if err := row.Scan(&user.ID, &user.Username, &user.HashedPassword, &user.CreatedAt, &user.Active); err != nil {
		return nil, fmt.Errorf("failed to find user %v: %w", user, err)
	}

	if !auth.CheckPasswordHash(password, user.HashedPassword) {
		return nil, errors.New("incorrect username or password")
	}

	return &user, nil
}
