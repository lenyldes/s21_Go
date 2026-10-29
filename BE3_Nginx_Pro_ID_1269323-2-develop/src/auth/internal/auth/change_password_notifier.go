package auth

import (
	"context"
	"fmt"
)

type ChangePasswordNotifier struct {
}

func NewChangePasswordNotifier() *ChangePasswordNotifier {
	return &ChangePasswordNotifier{}
}

func (n *ChangePasswordNotifier) Notify(ctx context.Context, login string, tmpPassword string, email string) error {
	fmt.Printf("new temporary password %s was sended to %s on email %s\n", tmpPassword, login, email)
	return nil
}
