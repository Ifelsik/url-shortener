package memstore

import (
	"context"
	"sync"

	"github.com/Ifelsik/url-shortener/internal/domain/url"
)

type urlStorage struct {
	id   uint64
	urls map[string]*url.URL
	mu   *sync.RWMutex
}

func NewURLStorage() *urlStorage {
	return &urlStorage{
		id:   0,
		urls: make(map[string]*url.URL),
		mu:   &sync.RWMutex{},
	}
}

func (s *urlStorage) Add(ctx context.Context, url *url.URL) (*url.URL, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	url.ID = s.id
	s.id++
	s.urls[url.ShortKey] = url
	return url, nil
}

func (s *urlStorage) GetByShortKey(ctx context.Context, 
		shortKey string) (*url.URL, error) {
	url, ok := s.urls[shortKey]
	if !ok {
		return nil, ErrNoURL
	}
	return url, nil
}
