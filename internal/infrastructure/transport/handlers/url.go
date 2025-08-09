package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Ifelsik/url-shortener/internal/app"
	"github.com/Ifelsik/url-shortener/internal/app/url"
	"github.com/Ifelsik/url-shortener/internal/pkg/logger"
	"github.com/gorilla/mux"
)

const (
	ShortURLSlug = "key"
)

type URLHandlers struct {
	logger     logger.Logger
	urlService *app.URLService
}

func NewURLHandlers(urlService *app.URLService, logger logger.Logger) *URLHandlers {
	return &URLHandlers{logger: logger, urlService: urlService}
}

func (h *URLHandlers) AddShortURL(w http.ResponseWriter, r *http.Request) {
	if log, err := logger.FromContext(r.Context()); err == nil {
		h.logger = log
	} else {
		h.logger.Warningf("AddShortURL http handler: %v", err)
	}

	h.logger.Debugf("Reading request body")
	var body bytes.Buffer
	_, err := body.ReadFrom(r.Body)

	defer func(){ 
			_ = r.Body.Close() 
		}()

	if err != nil {
		h.logger.Errorf("AddShortURL http handler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.logger.Debugf("Unmarshaling request body")
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
	if _, err = w.Write(responseJSON); err != nil {
		h.logger.Errorf("AddShortURL http handler: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *URLHandlers) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	if log, err := logger.FromContext(r.Context()); err == nil {
		h.logger = log
	} else {
		h.logger.Warningf("GetOriginalURL http handler: %v", err)
	}

	vars := mux.Vars(r)
	shortURL, ok := vars[ShortURLSlug]
	if !ok {
		h.logger.Errorf("GetOriginalURL http handler: %s var is empty", ShortURLSlug)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.logger.Debugf("GetOriginalURL http handler: got short url: %s", shortURL)

	getOriginalURLReq := &url.GetURLByShortRequest{
		ShortKey: shortURL,
	}
	getOriginalURLResp, err :=
		h.urlService.GetByShort.Handle(r.Context(), getOriginalURLReq)
	if err != nil {
		h.logger.Errorf("GetOriginalURL http handler: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Debugf("GetOriginalURL http handler: got data from app layer")
	response := GetURLByShortResponse{
		ShortURL:    getOriginalURLResp.ShortURL,
		OriginalURL: getOriginalURLResp.OriginalURL,
	}

	h.logger.Debugf("GetOriginalURL http handler: marshal response")
	responseJSON, err := json.Marshal(response)
	if err != nil {
		h.logger.Errorf("GetOriginalURL http handler: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Debugf("GetOriginalURL http handler: write response")
	if _, err = w.Write(responseJSON); err != nil {
		h.logger.Errorf("GetOriginalURL http handler: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
