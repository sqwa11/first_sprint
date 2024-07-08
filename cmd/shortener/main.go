package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sqwa11/first_sprint/config"
	"github.com/sqwa11/first_sprint/internal/app/get"
	"github.com/sqwa11/first_sprint/internal/app/post"
)

func main() {
	cfg := config.NewConfig()
	post.SetBaseURL(cfg.BaseURL)

	router := chi.NewRouter()
	router.Post("/", post.HandleShorten)
	router.Get("/{id}", get.HandleRedirect)

	log.Printf("Starting server at %s\n", cfg.ServerAddress)
	err := http.ListenAndServe(cfg.ServerAddress, router)
	if err != nil {
		log.Fatal(err)
	}
}
