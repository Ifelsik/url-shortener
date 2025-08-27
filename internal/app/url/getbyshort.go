package url

import (
	"context"
	"errors"
	"fmt"

	"github.com/Ifelsik/url-shortener/internal/domain/url"
	"github.com/Ifelsik/url-shortener/internal/pkg/validator"
)

type GetURLByShortRequest struct {
	ShortKey string `validate:"required"`
}

type GetURLByShortResponse struct {
	ShortURL    string
	OriginalURL string
}

type GetURLByShort interface {
	Handle(ctx context.Context,
		request *GetURLByShortRequest) (*GetURLByShortResponse, error)
}

type getURLByShort struct {
	urlRepo url.URLRepository
	val     validator.Validator
}

func NewGetURLByShortKey(urlRepo url.URLRepository, val validator.Validator) *getURLByShort {
	return &getURLByShort{urlRepo: urlRepo, val: val}
}

func (g *getURLByShort) Handle(ctx context.Context,
	request *GetURLByShortRequest) (*GetURLByShortResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("get url by short key: %w", ErrEmptyRequest)
	}

	err := g.val.ValidateStruct(request)
	if err != nil {
		return nil, fmt.Errorf("get url by short key: %w", err)
	}

	var shortKeyURL = request.ShortKey
	originalURL, err := g.urlRepo.GetByShortKey(ctx, shortKeyURL)
	if errors.Is(err, url.ErrNoURL) {
		return nil, fmt.Errorf("get url by short key: %w", ErrNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("get url by short key: %w", err)
	}

	result := &GetURLByShortResponse{
		ShortURL:    originalURL.ShortKey,
		OriginalURL: originalURL.OriginalURL,
	}

	return result, nil
}
