package url

import (
	"context"
	"fmt"

	neturl "net/url"

	"github.com/Ifelsik/url-shortener/internal/domain/url"
	"github.com/Ifelsik/url-shortener/internal/domain/user"
	"github.com/Ifelsik/url-shortener/internal/pkg/base62"
	"github.com/Ifelsik/url-shortener/internal/pkg/hasher"
	"github.com/Ifelsik/url-shortener/internal/pkg/timing"
	"github.com/Ifelsik/url-shortener/internal/pkg/validator"

)

type AddURLRequest struct {
	OriginalURL string `validate:"required,url_without_scheme"`
	UserToken   string `validate:"required"`
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
	validator      validator.Validator
}

func NewAddURL(
	addURLRepo url.URLRepository,
	addUserRepo user.UserRepository,
	timingProvider timing.Timing,
	base62Provider base62.Base62Provider,
	hasher hasher.Hasher,
	validationProvider validator.Validator,
) *addURL {
	return &addURL{
		urlRepo:        addURLRepo,
		userRepo:       addUserRepo,
		timingProvider: timingProvider,
		base62Provider: base62Provider,
		hasher:         hasher,
		validator:      validationProvider,
	}
}

func (a *addURL) Handle(ctx context.Context,
	request *AddURLRequest) (*AddURLResponse, error) {
	if request == nil {
		return nil, ErrEmptyRequest
	}

	err := a.validator.ValidateStruct(request)
	if err != nil {
		return nil, fmt.Errorf("add url: %w", err)
	}

	request.OriginalURL = prepareURL(request.OriginalURL)

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

	savedURL, err := a.urlRepo.Add(ctx, &url)
	if err != nil {
		return nil, fmt.Errorf("add url: %w", err)
	}

	result := &AddURLResponse{
		OriginalURL: savedURL.OriginalURL,
		ShortURL:    savedURL.ShortKey,
	}

	return result, nil
}

func prepareURL(url string) string {
	u, _ := neturl.Parse(url)
	if u.Scheme == "" {
		// 'https' is considered as default scheme
		return "https://" + url
	}
	return url
}
