package url

import "context"

type URLRepository interface {
	Add(ctx context.Context, url *URL) (*URL, error)
	GetByShortKey(ctx context.Context, shortKey string) (*URL, error)
}
