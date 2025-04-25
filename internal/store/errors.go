package store

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrDeviceNotFound = errors.New("device not found")
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
