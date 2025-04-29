package session

import (
	"context"
	"iotstarter/internal/model"
)

type Repository interface {
	Create(ctx context.Context, userId model.UserId) (*model.Session, error)
	Get(ctx context.Context, token model.SessionToken) (*model.Session, error)
	GetUserIdFromToken(ctx context.Context, token model.SessionToken) (model.UserId, error)
	Clear(ctx context.Context, userId model.UserId) error
}
