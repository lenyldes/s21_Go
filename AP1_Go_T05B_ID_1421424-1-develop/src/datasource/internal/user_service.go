package internal

import (
	"context"
	"tictactoe/domain"
)

type UserService struct {
	repository domain.UserRepository
}

func NewUserService(r domain.UserRepository) *UserService {
	return &UserService{repository: r}
}

func (s *UserService) Create(ctx context.Context, user domain.User) (string, error) {
	return s.repository.Create(ctx, user)
}

func (s *UserService) GetByUUID(ctx context.Context, UUID string) (domain.User, error) {
	return s.repository.GetByUUID(ctx, UUID)
}

func (s *UserService) GetByLogin(ctx context.Context, login string) (domain.User, error) {
	return s.repository.GetByLogin(ctx, login)
}
