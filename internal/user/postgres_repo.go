package user

import (
	"context"
	"errors"
	"fmt"
	"iotstarter/internal/model"
	"iotstarter/internal/security"
	"log"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(ctx context.Context, db *pgxpool.Pool) *postgresRepo {
	return &postgresRepo{db: db}
}

func (s *postgresRepo) Create(ctx context.Context, userName string, password string) (*model.User, error) {
	sql := `
		INSERT INTO users (username, hashed_password)
		VALUES ($1, $2)
		Returning id, username, created_at, active
	`

	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		log.Println("ERROR:", err.Error())
		return nil, security.ErrHashingError
	}

	user := model.User{
		Username:       userName,
		HashedPassword: hashedPassword,
	}

	row := s.db.QueryRow(ctx, sql, user.Username, user.HashedPassword)
	err = row.Scan(&user.ID, &user.Username, &user.CreatedAt, &user.Active)

	if isUniqueViolationError(err) {
		return nil, ErrUsernameTaken
	}
	if err != nil {
		return nil, fmt.Errorf("failed to register user %v: %w", user, err)
	}

	return &user, nil
}

func (s *postgresRepo) GetFromCreds(ctx context.Context, userName string, password string) (*model.User, error) {
	sql := `
		SELECT id, username, hashed_password, created_at, active
		FROM users 
		WHERE username = $1 AND active = TRUE
	`

	user := model.User{}

	row := s.db.QueryRow(ctx, sql, userName)
	if err := row.Scan(&user.ID, &user.Username, &user.HashedPassword, &user.CreatedAt, &user.Active); err != nil {
		if isNoRowsFoundError(err) {
			return nil, ErrUserNotExists
		}
		return nil, fmt.Errorf("failed to find user %v: %w", user, err)
	}

	if !security.CheckPasswordHash(password, user.HashedPassword) {
		return nil, errors.New("incorrect username or password")
	}

	return &user, nil
}

func (s *postgresRepo) GetById(ctx context.Context, userId model.UserId) (*model.User, error) {
	sql := `
		SELECT id, username, created_at, active
		FROM users 
		WHERE user_id = $1 AND active = TRUE
	`
	user := model.User{}

	row := s.db.QueryRow(ctx, sql, user.Username)
	if err := row.Scan(&user.ID, &user.Username, &user.CreatedAt, &user.Active); err != nil {
		if isNoRowsFoundError(err) {
			return nil, ErrUserNotExists
		}
		return nil, fmt.Errorf("failed to find user %v: %w", user, err)
	}

	return &user, nil

}

var ErrUsernameTaken = errors.New("username taken")
var ErrUserNotExists = errors.New("user does not exist")

func isUniqueViolationError(err error) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return true
	}
	return false
}

func isNoRowsFoundError(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
