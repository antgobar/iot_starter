package session

import (
	"context"
	"iotstarter/internal/model"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) GetUserFromToken(ctx context.Context, token model.SessionToken) (*model.User, error) {
	return s.repo.GetUserFromToken(ctx, token)
}
