package auth

import (
	"context"
	"errors"
	"iotstarter/internal/model"
)

type contextKey string

const userKey contextKey = "user"

var ErrNoUser = errors.New("no user in context")

func WithUser(ctx context.Context, u *model.User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func UserFromContext(ctx context.Context) (*model.User, error) {
	u, ok := ctx.Value(userKey).(*model.User)
	if !ok || u == nil {
		return nil, ErrNoUser
	}
	return u, nil
}
