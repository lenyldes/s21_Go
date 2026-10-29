package auth

import (
	"errors"
)

var ErrUserNotFound = errors.New("User not found")
var ErrLoginTaken = errors.New("Login already exists")

type ErrInvalidUserData struct {
	message string
}

func (e ErrInvalidUserData) Error() string {
	return e.message
}