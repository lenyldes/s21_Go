package auth

import (
	"context"

	"github.com/google/uuid"
)

type AuthRepository interface {
	SaveUser(ctx context.Context, user *User) error
	GetUserById(ctx context.Context, userId uuid.UUID) (*User, error)
	GetUserByLogin(ctx context.Context, login string) (*User, error)
	ChangeUserPassword(ctx context.Context, userId uuid.UUID, newPassword []byte) error
}
