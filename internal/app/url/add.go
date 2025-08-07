package url

import (
	"context"
	"fmt"

	"github.com/Ifelsik/url-shortener/internal/domain/url"
	"github.com/Ifelsik/url-shortener/internal/domain/user"
	"github.com/Ifelsik/url-shortener/internal/pkg/base62"
	"github.com/Ifelsik/url-shortener/internal/pkg/hasher"
	"github.com/Ifelsik/url-shortener/internal/pkg/timing"
)

type AddURLRequest struct {
	OriginalURL string
	UserToken   string
}

type AddURLResponse struct {
	OriginalURL string
	ShortURL    string
}

type AddURL interface {
	Handle(ctx context.Context, request *AddURLRequest) (*AddURLResponse, error)
}

type addURL struct {
	urlRepo        url.URLRepository
	userRepo       user.UserRepository
	timingProvider timing.Timing
	base62Provider base62.Base62Provider
	hasher         hasher.Hasher
}

func NewAddURL(
	addURLRepo url.URLRepository,
	addUserRepo user.UserRepository,
	timingProvider timing.Timing,
	base62Provider base62.Base62Provider,
	hasher hasher.Hasher) *addURL {
		return &addURL{
			urlRepo:        addURLRepo,
			userRepo:       addUserRepo,
			timingProvider: timingProvider,
			base62Provider: base62Provider,
			hasher:         hasher,
		}
}

func (a *addURL) Handle(ctx context.Context,
	request *AddURLRequest) (*AddURLResponse, error) {
	if request == nil {
		return nil, ErrEmptyRequest
	}
	// TODO: add validation
	user, err := a.userRepo.GetByToken(ctx, request.UserToken)
	if err != nil {
		return nil, fmt.Errorf("add url: %w", err)
	}

	urlHash := a.hasher.String(request.OriginalURL)
	shortKey := a.base62Provider.EncodeToString([]byte(urlHash))

	url := url.URL{
		OriginalURL: request.OriginalURL,
		User:        user.ID,
		CreatedAt:   a.timingProvider.Now(),
		ShortKey:    shortKey,
	}

	// Realization doesn't need to know about URL in repository
	savedURL, err := a.urlRepo.Add(ctx, &url)
	if err != nil {
		return nil, fmt.Errorf("add url: %w", err)
	}

	result := &AddURLResponse{
		OriginalURL: savedURL.OriginalURL,
		ShortURL:    savedURL.ShortKey,
	}

	// TODO: validate result
	return result, nil
}
