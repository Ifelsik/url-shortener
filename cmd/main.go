package main

import (
	"github.com/Ifelsik/url-shortener/internal/app"
	"github.com/Ifelsik/url-shortener/internal/infrastructure/storage/memory"
	"github.com/Ifelsik/url-shortener/internal/infrastructure/transport"
	"github.com/Ifelsik/url-shortener/internal/infrastructure/validator"
	"github.com/Ifelsik/url-shortener/internal/pkg/base62"
	"github.com/Ifelsik/url-shortener/internal/pkg/hasher"
	"github.com/Ifelsik/url-shortener/internal/pkg/identifier"
	"github.com/Ifelsik/url-shortener/internal/pkg/logger"
	"github.com/Ifelsik/url-shortener/internal/pkg/timing"

	appUrl "github.com/Ifelsik/url-shortener/internal/app/url"
	appUser "github.com/Ifelsik/url-shortener/internal/app/user"
)

func main() {
	userStorage := memory.NewUserStorage()
	urlStorage := memory.NewURLStorage()

	log := logger.NewLogrusLogWrap(nil)
	tp := timing.NewTimingProvider()
	b62 := base62.NewBase62Encoder()
	id := identifier.NewUUIDProvider()
	hasher := hasher.NewHasher32()
	valid := validator.NewValidator()

	App := app.Services{
		URLService: &app.URLService{
			AddURL: appUrl.NewAddURL(
				urlStorage,
				userStorage,
				tp,
				b62,
				hasher,
				valid,
			),
			GetByShort: appUrl.NewGetURLByShortKey(urlStorage, valid),
		},
		UserService: &app.UserService{
			AddUser: appUser.NewAddUser(userStorage, tp),
		},
	}

	httpServer := transport.NewHTTPServer(&App, log, id, tp)
	_ = httpServer.ListenAndServe()
}
