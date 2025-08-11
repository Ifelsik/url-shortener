package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Ifelsik/url-shortener/internal/domain/user"
	"github.com/Ifelsik/url-shortener/internal/pkg/identifier"
	"github.com/Ifelsik/url-shortener/internal/pkg/timing"
)

type AddUserRequest struct {
	UserToken string
}

type AddUserResponse struct {
	UserToken string
	ExpiresAt time.Time
}

type AddUser interface {
	Handle(ctx context.Context, request *AddUserRequest) (*AddUserResponse, error)
}

type addUser struct {
	userRepo   user.UserRepository
	timing     timing.Timing
	identifier identifier.Identifier
}

func NewAddUser(userRepo user.UserRepository, timing timing.Timing) *addUser {
	return &addUser{userRepo: userRepo, timing: timing, identifier: identifier.NewUUIDProvider()}
	return &addUser{userRepo: userRepo, timing: timing, identifier: identifier.NewUUIDProvider()}
}

func (a *addUser) Handle(
	ctx context.Context,
	request *AddUserRequest,
	) (*AddUserResponse, error) {
	existentUser, err := a.userRepo.GetByToken(ctx, request.UserToken)
	if err == nil {
		return &AddUserResponse{UserToken: existentUser.Token}, nil
	}

	newUser, err := a.userRepo.Add(
		ctx,
		&user.User{
			Token:     a.identifier.String(),
			CreatedAt: a.timing.Now(),
			ExpiresAt: a.timing.Now().Add(30 * 24 * time.Hour),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("add user: %w", err)
	}

	return &AddUserResponse{UserToken: newUser.Token}, nil
}
