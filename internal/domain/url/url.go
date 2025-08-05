package url

import "time"

// Struct URL represents entity of short url
type URL struct {
	ID          uint64
	OriginalURL string
	ShortKey    string
	User        uint64
	CreatedAt   time.Time
}
