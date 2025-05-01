package user

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Register(ctx context.Context, userName, password string) error {
	_, err := s.repo.Create(ctx, userName, password)
	return err
}
