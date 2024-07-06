package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

var (
	urlMap             = make(map[string]string) // Хранение сокращенных URL
	baseURL            = "http://localhost:8080" // Базовый URL вашего сервиса
	shortenedURLLength = 8                       // Длина сокращенного URL
)

func main() {
	http.HandleFunc("/", handleShorten)
	http.HandleFunc("/{id}", handleRedirect)

	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body := make([]byte, r.ContentLength)
	_, err := r.Body.Read(body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	longURL := strings.TrimSpace(string(body))
	shortURL := generateShortURL()

	urlMap[shortURL] = longURL

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s\n", baseURL, shortURL)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")
	longURL, exists := urlMap[id]
	if !exists {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}

func generateShortURL() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // Доступные символы для генерации

	b := make([]byte, shortenedURLLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
