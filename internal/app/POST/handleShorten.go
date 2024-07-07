package POST

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

var (
	UrlMap             = make(map[string]string) // Хранение сокращенных URL
	BaseURL            = "http://localhost:8080" // Базовый URL вашего сервиса
	ShortenedURLLength = 8                       // Длина сокращенного URL
)

func HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	longURL := strings.TrimSpace(string(body))
	shortURL := generateShortURL()

	UrlMap[shortURL] = longURL

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", BaseURL, shortURL)
}

func generateShortURL() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // Доступные символы для генерации

	b := make([]byte, ShortenedURLLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
