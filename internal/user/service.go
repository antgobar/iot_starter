package user

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

func (s *Service) Register(ctx context.Context, userName, password string) error {
	_, err := s.repo.Create(ctx, userName, password)
	return err
}

func (s *Service) GetById(ctx context.Context, userId model.UserId) (*model.User, error) {
	return s.repo.GetById(ctx, userId)
}
