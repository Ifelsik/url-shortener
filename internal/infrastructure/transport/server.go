package transport

import (
	"net/http"

	"github.com/Ifelsik/url-shortener/internal/app"
	"github.com/Ifelsik/url-shortener/internal/pkg/identifier"
	"github.com/Ifelsik/url-shortener/internal/pkg/logger"
	"github.com/Ifelsik/url-shortener/internal/pkg/timing"
)

type HTTPServer struct {
	Host   string
	Port   string
	srv    *http.Server
	logger logger.Logger
}

func NewHTTPServer(
	app *app.Services,
	log logger.Logger,
	ip identifier.Identifier,
	tp timing.Timing,
) *HTTPServer {
	mux := Router(app, log, ip, tp)
	return &HTTPServer{
		srv: &http.Server{
			Handler: mux,
		},
		logger: log,
	}
}

func (s *HTTPServer) ListenAndServe() error {
	s.logger.Infof("Server started on %s:%s", s.Host, s.Port)
	return s.srv.ListenAndServe()
}
