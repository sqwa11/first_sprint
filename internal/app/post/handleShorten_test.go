package post

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestHandleShortenWithGzip(t *testing.T) {
	router := chi.NewRouter()
	SetBaseURL("http://localhost:8080")
	router.Post("/", HandleShorten)

	longURL := "https://example.com"
	body := bytes.NewBufferString(longURL)

	// Сжатие тела запроса
	var compressedBody bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedBody)
	_, err := io.Copy(gzipWriter, body)
	require.NoError(t, err)
	gzipWriter.Close()

	req, err := http.NewRequest(http.MethodPost, "/", &compressedBody)
	require.NoError(t, err)
	req.Header.Set("Content-Encoding", "gzip")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Проверка сжатого ответа (ожидание несжатого ответа)
	responseBody := rr.Body.String()
	if !strings.HasPrefix(responseBody, "http://localhost:8080/") {
		t.Errorf("handler returned unexpected body: got %v", responseBody)
	}

	shortURL := strings.TrimSpace(responseBody)
	id := strings.TrimPrefix(shortURL, "http://localhost:8080/")
	savedURL, exists := URLMap[id]
	if !exists {
		t.Errorf("short URL not saved in map")
	}
	if savedURL != longURL {
		t.Errorf("saved long URL does not match: got %v want %v", savedURL, longURL)
	}
}

func TestHandleAPIPostShortenWithGzip(t *testing.T) {
	router := chi.NewRouter()
	SetBaseURL("http://localhost:8080")
	router.Post("/api/shorten", HandleAPIPostShorten)

	reqBody := struct {
		URL string `json:"url"`
	}{
		URL: "https://practicum.yandex.ru",
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	// Сжатие тела запроса
	var compressedBody bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedBody)
	_, err = gzipWriter.Write(body)
	require.NoError(t, err)
	gzipWriter.Close()

	req, err := http.NewRequest(http.MethodPost, "/api/shorten", &compressedBody)
	require.NoError(t, err)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Проверка сжатого ответа (ожидание несжатого ответа)
	var respBody struct {
		Result string `json:"result"`
	}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	shortURL := respBody.Result
	if !strings.HasPrefix(shortURL, "http://localhost:8080/") {
		t.Errorf("handler returned unexpected body: got %v", shortURL)
	}

	id := strings.TrimPrefix(shortURL, "http://localhost:8080/")
	savedURL, exists := URLMap[id]
	if !exists {
		t.Errorf("short URL not saved in map")
	}
	if savedURL != reqBody.URL {
		t.Errorf("saved long URL does not match: got %v want %v", savedURL, reqBody.URL)
	}
}
