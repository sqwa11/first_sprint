package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sqwa11/first_sprint/internal/app/get"
	"github.com/sqwa11/first_sprint/internal/app/post"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/", post.HandleShorten)
	r.Get("/{id}", get.HandleRedirect)

	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
