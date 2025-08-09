package transport

import (
	"context"
	"net/http"
	"time"

	"github.com/Ifelsik/url-shortener/internal/app"
	"github.com/Ifelsik/url-shortener/internal/pkg/identifier"
	"github.com/Ifelsik/url-shortener/internal/pkg/logger"
	"github.com/Ifelsik/url-shortener/internal/pkg/timing"
)

type HTTPServer struct {
	host   string
	port   string
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
			ReadHeaderTimeout: 5 * time.Second,
		},
		logger: log,
	}
}

func (s *HTTPServer) Shutdown() error {
	return s.srv.Shutdown(context.Background())
}

func (s *HTTPServer) ListenAndServe() error {
	s.logger.Infof("Server started on %s:%s", s.host, s.port)

	return s.srv.ListenAndServe()
}
