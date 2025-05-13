package session

import (
	"context"
	"iotstarter/internal/model"
	"iotstarter/internal/security"
	"sync"
	"time"
)

type memoryRepository struct {
	mu       sync.RWMutex
	sessions []*model.Session
}

func NewMemoryRepository() *memoryRepository {
	memrep := &memoryRepository{
		sessions: make([]*model.Session, 0),
	}
	return memrep
}

func (m *memoryRepository) Create(ctx context.Context, userId model.UserId) (*model.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	
	seshToken := model.SessionToken(security.GenerateUUID())
	sesh := model.Session{
		UserId:    userId,
		Token:     seshToken,
		ExpiresAt: time.Now().UTC().Add(3 * time.Hour),
	}

	m.sessions = append(m.sessions, &sesh)

	return &sesh, nil
}

func (m *memoryRepository) GetUserFromToken(ctx context.Context, token model.SessionToken) (*model.User, error) {
	return nil, nil
}