package post

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

var URLMap = map[string]string{}
var baseURL string

func SetBaseURL(url string) {
	baseURL = url
}

func HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var longURL string
	if err := json.NewDecoder(r.Body).Decode(&longURL); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL()
	URLMap[shortURL] = longURL

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(baseURL + "/" + shortURL))
}

func generateShortURL() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // Доступные символы для генерации

	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
