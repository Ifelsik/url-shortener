package handlers

type AddURLRequest struct {
	OriginalURL string `json:"original_url"`
}

type AddURLResponse struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

type GetURLByShortResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

const UserTokenCookie = "user_id"
