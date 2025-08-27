package user

import (
	"context"
	"errors"
)

type UserRepository interface {
	Add(ctx context.Context, user *User) (*User, error)
	GetByToken(ctx context.Context, token string) (*User, error)
}

var (
	ErrNoUser = errors.New("no such user")
)
