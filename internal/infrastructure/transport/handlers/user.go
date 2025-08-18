package handlers

import (
	"net/http"

	"github.com/Ifelsik/url-shortener/internal/app"
	"github.com/Ifelsik/url-shortener/internal/app/user"
	"github.com/Ifelsik/url-shortener/internal/pkg/logger"
)

type UserHandlers struct {
	log         logger.Logger
	userService *app.UserService
}

func NewUserHandlers(
	userService *app.UserService,
	logger logger.Logger,
) *UserHandlers {
	return &UserHandlers{
		log:         logger,
		userService: userService,
	}
}

func (h *UserHandlers) AddUser(w http.ResponseWriter, r *http.Request) {
	log, err := logger.FromContext(r.Context())
	if err == nil {
		h.log = log
	} else {
		h.log.Warningf("AddUser http handler: %v", err)
	}

	tokenCookie, err := r.Cookie(UserTokenCookie)
	if err != nil {
		h.log.Debugf("AddUser http handler: %w", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	var tokenValue = tokenCookie.Value
	h.log.Debugf("AddUser http handler: creating new user")
	user, err := h.userService.AddUser.Handle(
		r.Context(),
		&user.AddUserRequest{
			UserToken: tokenValue,
		},
	)
	if err != nil {
		h.log.Errorf("AddUser http handler: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	h.log.Debugf("AddUser http handler: forming cookie")
	userCookie := &http.Cookie{
		Name:     UserTokenCookie,
		Value:    user.UserToken,
		HttpOnly: true,
		Expires:  user.ExpiresAt,
	}

	h.log.Debugf("AddUser http handler: setting cookie")
	http.SetCookie(w, userCookie)
	w.WriteHeader(http.StatusCreated)
}
