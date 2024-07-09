package post

import (
	"github.com/sqwa11/first_sprint/config"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleShorten(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(HandleShorten))
	defer ts.Close()

	cfg := &config.Config{
		ServerAddress: ts.URL,
		BaseURL:       ts.URL,
	}
	SetBaseURL(cfg.BaseURL)

	longURL := "https://google.com"
	body := strings.NewReader(longURL)
	resp, err := http.Post(ts.URL, "text/plain", body)
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
