package url

import (
	"context"
	"fmt"

	"github.com/Ifelsik/url-shortener/internal/domain/url"
	"github.com/Ifelsik/url-shortener/internal/domain/user"
	"github.com/Ifelsik/url-shortener/internal/pkg/timing"
)

type AddURLRequest struct {
	OriginalURL string
	ShortKey    string
	UserToken   string
}

type AddURL interface {
	Handle(ctx context.Context, request *AddURLRequest) error
}

type addURL struct {
	urlRepo        url.URLRepository
	userRepo       user.UserRepository
	timingProvider timing.Timing
}

func NewAddURL(
	addURLRepo url.URLRepository,
	addUserRepo user.UserRepository,
	timingProvider timing.Timing) *addURL {
	return &addURL{
		urlRepo:        addURLRepo,
		userRepo:       addUserRepo,
		timingProvider: timingProvider,
	}
}

func (a *addURL) Handle(ctx context.Context, request *AddURLRequest) error {
	if request == nil {
		return ErrEmptyRequest
	}
	// TODO: add validation
	user, err := a.userRepo.GetByToken(ctx, request.UserToken)
	if err != nil {
		return fmt.Errorf("add url: %w", err)
	}

	url := url.URL{
		OriginalURL: request.OriginalURL,
		ShortKey:    request.ShortKey,
		User:        user.ID,
		CreatedAt:   a.timingProvider.Now(),
	}

	// Realization doesn't need to know about URL in repository
	_, err = a.urlRepo.Add(ctx, &url)
	if err != nil {
		return fmt.Errorf("add url: %w", err)
	}

	return nil
}
