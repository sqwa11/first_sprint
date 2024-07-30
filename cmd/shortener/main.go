package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sqwa11/first_sprint/internal/app/config"
	"github.com/sqwa11/first_sprint/internal/app/get"
	"github.com/sqwa11/first_sprint/internal/app/post"
)

func main() {
	cfg := config.NewConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	r := chi.NewRouter()
	post.SetBaseURL(cfg.BaseURL)

	r.Post("/api/shorten", post.HandleShorten)
	r.Get("/{id}", get.HandleRedirect)

	log.Printf("Server listening on address %s...\n", cfg.Address)
	if err := http.ListenAndServe(cfg.Address, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
