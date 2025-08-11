package middleware

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
)

// Custom ResponseWriter with support of status code.
// Used for logging purposes.
type StatusResponseWriter struct {
	http.ResponseWriter

	Status int
}

var ErrNoHijacker = errors.New(" Hijacker is not supported")

func NewStatusResponseWriter(w http.ResponseWriter) *StatusResponseWriter {
	return &StatusResponseWriter{ResponseWriter: w, Status: 0}
}

// Overrides Write method.
// According to https://pkg.go.dev/net/http#ResponseWriter
// if WriteHeader isn't called explicitly, status code will be set to 200
// so we need to save this http status code.
func (srw *StatusResponseWriter) Write(b []byte) (int, error) {
	if srw.Status == 0 {
		srw.Status = http.StatusOK
	}

	return srw.ResponseWriter.Write(b)
}

// Overrides WriteHeader method.
// Saves http status code.
func (srw *StatusResponseWriter) WriteHeader(status int) {
	srw.Status = status
	srw.ResponseWriter.WriteHeader(status)
}

func (srw *StatusResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := srw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil,
			nil,
			fmt.Errorf("StatusResponseWriter: %w", ErrNoHijacker)
	}

	return h.Hijack()
}

func (srw *StatusResponseWriter) Flush() {
	fl, ok := srw.ResponseWriter.(http.Flusher)
	if ok {
		fl.Flush()
	}
}
