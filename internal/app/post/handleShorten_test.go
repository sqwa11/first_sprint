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
	body := newGzipBuffer(t, longURL)
	req, err := http.NewRequest(http.MethodPost, "/", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Encoding", "gzip")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	responseBody, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

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

func newGzipBuffer(t *testing.T, data string) *bytes.Buffer {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write([]byte(data))
	require.NoError(t, err)
	require.NoError(t, gz.Close())
	return &buf
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
	if err != nil {
		t.Fatal(err)
	}
	gzipBody := newGzipBuffer(t, string(body))

	req, err := http.NewRequest(http.MethodPost, "/api/shorten", gzipBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var respBody struct {
		Result string `json:"result"`
	}
	zr, err := gzip.NewReader(rr.Body)
	require.NoError(t, err)
	err = json.NewDecoder(zr).Decode(&respBody)
	require.NoError(t, err)
	require.NoError(t, zr.Close())

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
