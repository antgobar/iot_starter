package session

import (
	"context"
	"iotstarter/internal/model"
)

type Repository interface {
	Create(ctx context.Context, userId model.UserId) (*model.Session, error)
	GetUserFromToken(ctx context.Context, token model.SessionToken) (*model.User, error)
	Clear(ctx context.Context, userId model.UserId) error
}
