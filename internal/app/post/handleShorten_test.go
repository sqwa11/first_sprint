package post

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sqwa11/first_sprint/config"
)

func TestHandleShorten(t *testing.T) {
	cfg := config.NewConfig()
	SetBaseURL(cfg.BaseURL)

	router := chi.NewRouter()
	router.Post("/", HandleShorten)

	ts := httptest.NewServer(router)
	defer ts.Close()

	longURL := "https://example.com"
	body := strings.NewReader(longURL)
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/", body)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", resp.StatusCode, http.StatusCreated)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	shortURL := strings.TrimSpace(string(responseBody))
	if !strings.HasPrefix(shortURL, cfg.BaseURL+"/") {
		t.Errorf("handler returned unexpected body: got %v", shortURL)
	}

	id := strings.TrimPrefix(shortURL, cfg.BaseURL+"/")
	savedURL, exists := URLMap[id]
	if !exists {
		t.Errorf("short URL not saved in map")
	}
	if savedURL != longURL {
		t.Errorf("saved long URL does not match: got %v want %v", savedURL, longURL)
	}
}
