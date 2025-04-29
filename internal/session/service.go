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

func (s *Service) GetUserIdFromToken(ctx context.Context, token model.SessionToken) (model.UserId, error) {
	return s.repo.GetUserIdFromToken(ctx, token)
}
