package post

import (
	"encoding/json"
	"github.com/sqwa11/first_sprint/pkg/urlshortener"
	"net/http"
)

var URLMap = map[string]string{}
var baseURL string

func SetBaseURL(url string) {
	baseURL = url
}

func HandleShorten(w http.ResponseWriter, r *http.Request) {
	var longURL string
	if err := json.NewDecoder(r.Body).Decode(&longURL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL := urlshortener.Shorten(longURL)
	URLMap[shortURL] = longURL

	response := map[string]string{"result": baseURL + "/" + shortURL}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
