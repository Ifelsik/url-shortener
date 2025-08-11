package handlers

type AddURLRequest struct {
	OriginalURL string `json:"originalUrl"`
}

type AddURLResponse struct {
	OriginalURL string `json:"originalUrl"`
	ShortURL    string `json:"shortUrl"`
}

type GetURLByShortResponse struct {
	ShortURL    string `json:"shortUrl"`
	OriginalURL string `json:"originalUrl"`
}

const UserTokenCookie = "user_id"

const QueryShortURL = "short"
