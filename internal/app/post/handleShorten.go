package post

import (
	"compress/gzip"
	"io"
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

	body, err := readBody(r)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	longURL := strings.TrimSpace(string(body))
	if longURL == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL()
	URLMap[shortURL] = longURL

	response := baseURL + "/" + shortURL
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Возвращаем статус 201 Created
	w.Write([]byte(`{"result":"` + response + `"}`))
}

func readBody(r *http.Request) ([]byte, error) {
	var reader io.Reader = r.Body
	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		gzipReader, err := gzip.NewReader(r.Body)
		if err != nil {
			return nil, err
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	return io.ReadAll(reader)
}

func generateShortURL() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
