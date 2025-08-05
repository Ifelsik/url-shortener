package url

import (
	"context"
	"fmt"

	"github.com/Ifelsik/url-shortener/internal/domain/url"
)

type GetURLByShortRequest struct {
	ShortKey string
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
}

func NewGetURLByShortKey(urlRepo url.URLRepository) *getURLByShort {
	return &getURLByShort{urlRepo: urlRepo}
}

func (g *getURLByShort) Handle(ctx context.Context,
	request *GetURLByShortRequest) (*GetURLByShortResponse, error) {
	if request == nil {
		return nil, ErrEmptyRequest
	}

	var shortKeyURL string = request.ShortKey
	url, err := g.urlRepo.GetByShortKey(ctx, shortKeyURL)
	if err != nil {
		return nil, fmt.Errorf("get url by short key: %w", err)
	}

	// TODO: probably need to add a validation
	result := &GetURLByShortResponse{
		ShortURL:    url.ShortKey,
		OriginalURL: url.OriginalURL,
	}
	
	return result, nil
}
