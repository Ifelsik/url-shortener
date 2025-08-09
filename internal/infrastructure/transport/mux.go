package transport

import (
	"fmt"

	"github.com/Ifelsik/url-shortener/internal/app"
	"github.com/Ifelsik/url-shortener/internal/infrastructure/transport/handlers"
	"github.com/Ifelsik/url-shortener/internal/infrastructure/transport/middleware"
	"github.com/Ifelsik/url-shortener/internal/pkg/identifier"
	"github.com/Ifelsik/url-shortener/internal/pkg/logger"
	"github.com/Ifelsik/url-shortener/internal/pkg/timing"

	"github.com/gorilla/mux"
)

func Router(
	appSrv *app.Services,
	log logger.Logger,
	id identifier.Identifier,
	timing timing.Timing,
) *mux.Router {
	r := mux.NewRouter()

	userHandlers := handlers.NewUserHandlers(appSrv.UserService, log, timing)
	urlHandlers := handlers.NewURLHandlers(appSrv.URLService, log)

	r.HandleFunc("/user", userHandlers.AddUser).Methods("POST")
	r.HandleFunc("/url", urlHandlers.AddShortURL).Methods("POST")
	r.HandleFunc(
		fmt.Sprintf("/{%s}", handlers.ShortURLSlug),
		urlHandlers.GetOriginalURL).Methods("GET")

	// set up middleware
	panicMiddleware := middleware.NewPanicRecoveryMiddleware(log)
	r.Use(panicMiddleware.Middleware)

	logMiddleware := middleware.NewLoggingMiddleware(log, id, timing)
	r.Use(logMiddleware.Middleware)

	return r
}
