//nolint:dupl
package memory

import (
	"context"
	"sync"

	"github.com/Ifelsik/url-shortener/internal/domain/user"
)

type userStorage struct {
	id    uint64
	users map[string]*user.User
	mu    *sync.RWMutex
}

func NewUserStorage() *userStorage {
	return &userStorage{
		id:    0,
		users: make(map[string]*user.User),
		mu:    &sync.RWMutex{},
	}
}

func (s *userStorage) Add(ctx context.Context,
	user *user.User) (*user.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user.ID = s.id
	s.id++
	s.users[user.Token] = user

	return user, nil
}

func (s *userStorage) GetByToken(ctx context.Context,
	token string) (*user.User, error) {
	user, ok := s.users[token]
	if !ok {
		return nil, ErrNoUser
	}

	return user, nil
}
