package url

import (
	"context"
	"errors"
)

type URLRepository interface {
	Add(ctx context.Context, url *URL) (*URL, error)
	GetByShortKey(ctx context.Context, shortKey string) (*URL, error)
}

var (
	ErrNoURL = errors.New("no such url")
)
