package user

import (
	"context"
	"iotstarter/internal/model"
)

type Repository interface {
	Create(ctx context.Context, userName string, password string) (*model.User, error)
	GetFromCreds(ctx context.Context, userName string, password string) (*model.User, error)
}
