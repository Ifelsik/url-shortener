package user

import (
	"context"
	"fmt"

	"github.com/Ifelsik/url-shortener/internal/domain/user"
	"github.com/Ifelsik/url-shortener/internal/pkg/identifier"
	"github.com/Ifelsik/url-shortener/internal/pkg/timing"
)

type AddUserResponse struct {
	UserToken string
}

type AddUser interface {
	Handle(ctx context.Context) (*AddUserResponse, error)
}

type addUser struct {
	userRepo   user.UserRepository
	timing     timing.Timing
	identifier identifier.Identifier
}

func NewAddUser(userRepo user.UserRepository, timing timing.Timing) *addUser {
	return &addUser{userRepo: userRepo, timing: timing}
}

func (a *addUser) Handle(ctx context.Context) (*AddUserResponse, error) {
	user, err := a.userRepo.Add(
		ctx,
		&user.User{
			Token:     a.identifier.String(),
			CreatedAt: a.timing.Now(),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("add user: %w", err)
	}

	return &AddUserResponse{UserToken: user.Token}, nil
}
