package session

import (
	"context"
	"iotstarter/internal/model"
	"iotstarter/internal/security"
	"slices"
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
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, session := range m.sessions {
		if session.Token == token {
			return &model.User{ID: session.UserId}, nil
		}
	}

	return nil, noUserSessionErr(nil)
}

func (m *memoryRepository) Clear(ctx context.Context, userId model.UserId) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, session := range m.sessions {
		if session.UserId == userId {
			m.sessions = slices.Delete(m.sessions, i, i+1)
			return nil
		}
	}
	return failedToDeleteUserSessionErr(userId, nil)
}
