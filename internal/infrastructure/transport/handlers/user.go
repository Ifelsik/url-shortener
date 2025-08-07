package handlers

import (
	"net/http"
	"time"

	"github.com/Ifelsik/url-shortener/internal/app"
	"github.com/Ifelsik/url-shortener/internal/pkg/logger"
	"github.com/Ifelsik/url-shortener/internal/pkg/timing"
)

type UserHandlers struct {
	log          logger.Logger
	userService  *app.UserService
	timeProvider timing.Timing
}

func NewUserHandlers(
	userService *app.UserService,
	logger logger.Logger,
	timeProvider timing.Timing,
) *UserHandlers {
	return &UserHandlers{
		log:         logger,
		userService: userService,
		timeProvider: timeProvider,
	}
}

func (h *UserHandlers) AddUser(w http.ResponseWriter, r *http.Request) {
	if log, err := logger.FromContext(r.Context()); err == nil {
		h.log = log
	} else {
		h.log.Warningf("AddUser http handler: %v", err)
	}

	tokenCookie, err := r.Cookie(UserTokenCookie)
	if err == nil {
		h.log.Debugf("AddUser http handler: got user token: %s", tokenCookie.Value)
		w.WriteHeader(http.StatusOK)
		return
	}
	// TODO: refresh user token

	h.log.Debugf("AddUser http handler: creating new user")
	newUser, err := h.userService.AddUser.Handle(r.Context())
	if err != nil {
		h.log.Errorf("AddUser http handler: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.log.Debugf("AddUser http handler: forming cookie")
	userTokenCookie := &http.Cookie{
		Name:     UserTokenCookie,
		Value:    newUser.UserToken,
		HttpOnly: true,
		Expires: h.timeProvider.AfterNow(30 * 24 * time.Hour),
	}

	h.log.Debugf("AddUser http handler: setting cookie")
	http.SetCookie(w, userTokenCookie)
	w.WriteHeader(http.StatusOK)
}
