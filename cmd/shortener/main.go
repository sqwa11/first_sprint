package main

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sqwa11/first_sprint/internal/app/config"
	"github.com/sqwa11/first_sprint/internal/app/get"
	"github.com/sqwa11/first_sprint/internal/app/post"
	"go.uber.org/zap"
)

func main() {
	cfg := config.NewConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	logger := NewLogger()
	defer logger.Sync()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(WithLogging(logger)) // Добавляем middleware логирования
	r.Use(WithGzipCompression) // Middleware для gzip

	post.SetBaseURL(cfg.BaseURL)

	r.Post("/api/shorten", post.HandleShorten)
	r.Get("/{id}", get.HandleRedirect)

	log.Printf("Server listening on address %s...\n", cfg.Address)
	if err := http.ListenAndServe(cfg.Address, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func NewLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}

func WithLogging(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			responseData := &responseData{
				status: 0,
				size:   0,
			}
			lw := loggingResponseWriter{
				ResponseWriter: w,
				responseData:   responseData,
			}
			h.ServeHTTP(&lw, r)

			duration := time.Since(start)

			logger.Infow("Handled request",
				"uri", r.RequestURI,
				"method", r.Method,
				"status", responseData.status,
				"duration", duration,
				"size", responseData.size,
			)
		})
	}
}

func WithGzipCompression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acceptEncoding := r.Header.Get("Accept-Encoding")
		if strings.Contains(acceptEncoding, "gzip") {
			gzipWriter := gzip.NewWriter(w)
			defer gzipWriter.Close()
			w.Header().Set("Content-Encoding", "gzip")
			w = &gzipResponseWriter{ResponseWriter: w, Writer: gzipWriter}
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			r.Body = &gzipReadCloser{ReadCloser: r.Body}
		}

		next.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

type gzipReadCloser struct {
	io.ReadCloser
}

func (rc *gzipReadCloser) Read(p []byte) (int, error) {
	return rc.ReadCloser.Read(p)
}
