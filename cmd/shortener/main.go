package main

import (
	"compress/gzip"
	"go.uber.org/zap"
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
	r.Use(WithLogging(logger))
	r.Use(DecompressMiddleware) // Добавляем middleware для распаковки запросов
	r.Use(CompressMiddleware)   // Добавляем middleware для сжатия ответов

	post.SetBaseURL(cfg.BaseURL)

	r.Post("/", post.HandleShorten)
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

// Middleware для распаковки запросов с Content-Encoding: gzip
func DecompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress request body", http.StatusInternalServerError)
				return
			}
			defer gz.Close()
			r.Body = gz
		}
		next.ServeHTTP(w, r)
	})
}

// Middleware для сжатия ответов с Accept-Encoding: gzip
func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acceptEncoding := r.Header.Get("Accept-Encoding")
		if acceptEncoding != "" && containsGzip(acceptEncoding) {
			gzw := gzip.NewWriter(w)
			defer gzw.Close()
			w.Header().Set("Content-Encoding", "gzip")
			w = gzipResponseWriter{Writer: gzw, ResponseWriter: w}
		}
		next.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func containsGzip(acceptEncoding string) bool {
	for _, encoding := range []string{"gzip", "x-gzip"} {
		if strings.Contains(acceptEncoding, encoding) {
			return true
		}
	}
	return false
}
