package session

import (
	"iotstarter/internal/model"
	"sync"
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

func (m *memoryRepository) foo() {}
