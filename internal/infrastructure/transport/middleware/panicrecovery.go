package middleware

import (
	"net/http"

	"github.com/Ifelsik/url-shortener/internal/pkg/logger"
)

type PanicRecoveryMiddleware struct {
	log logger.Logger
}

func NewPanicRecoveryMiddleware(log logger.Logger) *PanicRecoveryMiddleware {
	return &PanicRecoveryMiddleware{log: log}
}

func (p *PanicRecoveryMiddleware) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				p.log.Errorf("Panic recovered: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		handler.ServeHTTP(w, r)
	})
}
