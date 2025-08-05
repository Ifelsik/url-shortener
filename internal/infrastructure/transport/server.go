package transport

import (
	"net/http"

	"github.com/Ifelsik/url-shortener/internal/pkg/logger"
)

type HTTPServer struct {
	srv *http.Server
	logger logger.Logger
}

func NewHTTPServer(log logger.Logger) *HTTPServer {
	return &HTTPServer{
		srv: &http.Server{
			Handler: http.DefaultServeMux,
		},
		logger: log,
	}
}

func (s *HTTPServer) ListenAndServe() error {
	return s.srv.ListenAndServe()
}
