package auth

import (
	"context"
	"errors"
	"iotstarter/internal/model"
)

type contextKey string

const userKey contextKey = "user"

var ErrNoUser = errors.New("no user in context")

func WithUserId(ctx context.Context, u model.UserId) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func UserIdFromContext(ctx context.Context) (model.UserId, bool) {
	u, ok := ctx.Value(userKey).(model.UserId)
	return u, ok
}
