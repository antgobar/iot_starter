package session

import (
	"context"
	"iotstarter/internal/model"
	"testing"
)

func TestMemoryRepoCreate(t *testing.T) {
	ctx := context.TODO()
	userId := model.UserId(123)
	memrep := NewMemoryRepository()

	memrep.Create(ctx, userId)

	session := memrep.sessions[0]

	if session.UserId != userId {
		t.Errorf("want %v, got %v", userId, session.UserId)
	}
}

func TestMemoryRepoGetUserFromToken(t *testing.T) {
	ctx := context.TODO()
	userId := model.UserId(123)
	memrep := NewMemoryRepository()

	memrep.Create(ctx, userId)
	token := memrep.sessions[0].Token

	user, err := memrep.GetUserFromToken(ctx, token)
	if user.ID != userId || err != nil {
		t.Errorf("want %v, got %v, err: %v", userId, user.ID, err)
	}

}
