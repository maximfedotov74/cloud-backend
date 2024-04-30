package mw

import (
	"log"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.status = statusCode
}

func ApplyLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)
		log.Printf("Method: %s; URL: %s; StatusCode: %d; StatusText: %s", r.Method, r.RequestURI, wrapped.status,
			http.StatusText(wrapped.status))
	})
}
