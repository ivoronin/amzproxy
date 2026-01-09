package main

import (
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		log.Printf("→ %s %s?%s from %s", r.Method, r.URL.Path, r.URL.RawQuery, r.RemoteAddr)

		next(lrw, r)

		duration := time.Since(start)
		log.Printf("← %d %s %s [%s]", lrw.statusCode, r.Method, r.URL.Path, duration)
	}
}
