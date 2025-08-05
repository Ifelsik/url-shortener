package url

type URLRepository interface {
	Add(url *URL) error
	GetByShortKey(shortKey string) (*URL, error)
}

