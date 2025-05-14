package session

import (
	"fmt"
	"iotstarter/internal/model"
)

func noUserSessionErr(err error) error {
	return fmt.Errorf("failed to retrieve user id from token: %w", err)
}

func failedToDeleteUserSessionErr(userId model.UserId, err error) error {
	return fmt.Errorf("failed to delete existing sessions for user id %d: %w", userId, err)
}
