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

func UserFromContext(ctx context.Context) (*model.User, bool) {
	u, ok := ctx.Value(userKey).(*model.User)
	return u, ok
}
