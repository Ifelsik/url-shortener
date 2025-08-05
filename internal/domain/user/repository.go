package user

import "context"

type UserRepository interface {
	Add(ctx context.Context, user *User) (*User, error)
	GetByToken(ctx context.Context, token string) (*User, error)
}
