package auth

import (
	"context"
	"iotstarter/internal/model"
	"iotstarter/internal/session"
	"iotstarter/internal/user"
)

type Service struct {
	users    user.Repository
	sessions session.Repository
}

func NewService(u user.Repository, s session.Repository) *Service {
	return &Service{users: u, sessions: s}
}

func (s *Service) LogIn(ctx context.Context, username, password string) (*model.Session, error) {
	user, err := s.users.GetFromCreds(ctx, username, password)
	if err != nil {
		return nil, err
	}

	sesh, err := s.sessions.Create(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return sesh, nil
}

func (s *Service) LogOut(ctx context.Context, userId model.UserId) error {
	return s.sessions.Clear(ctx, userId)
}
