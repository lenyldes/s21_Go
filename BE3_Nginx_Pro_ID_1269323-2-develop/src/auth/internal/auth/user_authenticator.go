package auth

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserProvider interface {
	GetUserById(ctx context.Context, userId uuid.UUID) (*User, error)
}

type UserAuthenticator struct{}

func NewUserAuthenticator(up UserProvider) *UserAuthenticator {
	return &UserAuthenticator{}
}

func (ua *UserAuthenticator) VerifyCredentials(ctx context.Context, user *User, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return err
	}
	return nil
}
