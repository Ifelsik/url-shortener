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

	url := url.URL{
		OriginalURL: request.OriginalURL,
		User:        user.ID,
		CreatedAt:   a.timingProvider.Now(),
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
