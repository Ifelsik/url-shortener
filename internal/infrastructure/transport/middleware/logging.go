package middleware

import (
	"net/http"

	"github.com/Ifelsik/url-shortener/internal/pkg/identifier"
	"github.com/Ifelsik/url-shortener/internal/pkg/logger"
	"github.com/Ifelsik/url-shortener/internal/pkg/timing"
)

type LoggingMiddleware struct {
	log          logger.Logger
	idgen        identifier.Identifier
	timeProvider timing.Timing
}

func NewLoggingMiddleware(
	log logger.Logger,
	idgen identifier.Identifier,
	timeProvider timing.Timing,
) *LoggingMiddleware {
	return &LoggingMiddleware{
		log:          log.WithFields(logger.LoggerFields{}),
		idgen:        idgen,
		timeProvider: timeProvider,
	}
}

func (m *LoggingMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.log = m.log.WithFields(logger.LoggerFields{
			"method":    r.Method,
			"URL":       r.URL.Path,
			"requestID": m.idgen.String(),
		})
		m.log.Debugf("got request")
		ctxWithLogger := logger.ToContext(r.Context(), m.log)
		startTime := m.timeProvider.Now()
		next.ServeHTTP(w, r.WithContext(ctxWithLogger))
		elapsedTime := m.timeProvider.Since(startTime)
		m.log = m.log.WithFields(logger.LoggerFields{
			"elapsed": elapsedTime,
		})
		m.log.Debugf("request served")
	})
}
