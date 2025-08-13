package transport

import (
	"context"
	"net/http"
	"time"

	"github.com/Ifelsik/url-shortener/internal/app"
	"github.com/Ifelsik/url-shortener/internal/infrastructure/config"
	"github.com/Ifelsik/url-shortener/internal/pkg/identifier"
	"github.com/Ifelsik/url-shortener/internal/pkg/logger"
	"github.com/Ifelsik/url-shortener/internal/pkg/timing"
)

type HTTPServer struct {
	srv    *http.Server
	conf   config.Server
	logger logger.Logger
}

func NewHTTPServer(
	conf config.Server,
	app *app.Services,
	log logger.Logger,
	ip identifier.Identifier,
	tp timing.Timing,
) *HTTPServer {
	mux := Router(app, log, ip, tp)

	if conf.Host == "" || conf.Port == "" {
		log.Warningf("http server: host or port is empty %s:%s",
			conf.Host, conf.Port)
	}

	return &HTTPServer{
		srv: &http.Server{
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		},
		logger: log,
		conf:   conf,
	}
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *HTTPServer) ListenAndServe() error {
	s.logger.Infof("Server started on %s:%s", s.conf.Host, s.conf.Port)
	return s.srv.ListenAndServe()
}
