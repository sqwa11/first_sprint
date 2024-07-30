package post

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
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

	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	longURL := strings.TrimSpace(string(body))

	if longURL == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL()
	URLMap[shortURL] = longURL

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(baseURL + "/" + shortURL))
}

func HandleAPIPostShorten(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	longURL := strings.TrimSpace(req.URL)
	if longURL == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL()
	URLMap[shortURL] = longURL

	resp := struct {
		Result string `json:"result"`
	}{
		Result: baseURL + "/" + shortURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func generateShortURL() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
