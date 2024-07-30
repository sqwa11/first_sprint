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

	// Создание сжатого тела запроса
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

	// Распаковка сжатого ответа
	gzipReader, err := gzip.NewReader(rr.Body)
	require.NoError(t, err)
	defer gzipReader.Close()

	responseBody, err := io.ReadAll(gzipReader)
	require.NoError(t, err)

	shortURL := strings.TrimSpace(string(responseBody))
	if !strings.HasPrefix(shortURL, "http://localhost:8080/") {
		t.Errorf("handler returned unexpected body: got %v", shortURL)
	}

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

	// Создание сжатого тела запроса
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

	// Распаковка сжатого ответа
	gzipReader, err := gzip.NewReader(rr.Body)
	require.NoError(t, err)
	defer gzipReader.Close()

	var respBody struct {
		Result string `json:"result"`
	}
	if err := json.NewDecoder(gzipReader).Decode(&respBody); err != nil {
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
