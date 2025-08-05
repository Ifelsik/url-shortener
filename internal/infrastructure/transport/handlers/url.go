package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Ifelsik/url-shortener/internal/app"
	"github.com/Ifelsik/url-shortener/internal/app/url"
	"github.com/Ifelsik/url-shortener/internal/pkg/logger"
)

type HTTPHandlers struct {
	logger     logger.Logger
	urlService *app.URLService
}

func NewHTTPHandlers(urlService *app.URLService, logger logger.Logger) *HTTPHandlers {
	return &HTTPHandlers{logger: logger, urlService: urlService}
}

func (h *HTTPHandlers) AddShortURL(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("Reading request body")
	var body bytes.Buffer
	_, err := body.ReadFrom(r.Body)
	defer r.Body.Close()
	if err != nil {
		h.logger.Errorf("AddShortURL http handler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.logger.Debugf("Unmarshalling request body")
	jsonBody := new(AddURLRequest)
	err = json.Unmarshal(body.Bytes(), jsonBody)
	if err != nil {
		h.logger.Errorf("AddShortURL http handler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.logger.Debugf("Reading cookie %s", UserTokenCookie)
	userToken, err := r.Cookie(UserTokenCookie)
	if err != nil {
		h.logger.Errorf("AddShortURL http handler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if userToken == nil || userToken.Value == "" {
		h.logger.Errorf("AddShortURL http handler: %s cookie is empty or nil",
			UserTokenCookie)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.logger.Debugf("Calling application layer method to add short url")
	addURLRequest := &url.AddURLRequest{
		OriginalURL: jsonBody.OriginalURL,
		UserToken:   userToken.Value,
	}
	addURLResponse, err := h.urlService.AddURL.Handle(r.Context(), addURLRequest)
	if err != nil {
		h.logger.Errorf("AddShortURL http handler: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := &AddURLResponse{
		OriginalURL: addURLResponse.OriginalURL,
		ShortURL:    addURLResponse.ShortURL,
	}
	h.logger.Debugf("Marshalling response")
	responseJSON, err := json.Marshal(response)
	if err != nil {
		h.logger.Errorf("AddShortURL http handler: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Debugf("Writing response")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func (h *HTTPHandlers) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}
