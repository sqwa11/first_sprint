package get

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sqwa11/first_sprint/config"
	"github.com/sqwa11/first_sprint/internal/app/post"
)

func TestHandleRedirect(t *testing.T) {
	cfg := config.NewConfig()
	post.SetBaseURL(cfg.BaseURL)

	router := chi.NewRouter()
	router.Get("/{id}", HandleRedirect)

	shortURL := "abcd1234"
	longURL := "https://example.com"
	post.URLMap[shortURL] = longURL

	req, err := http.NewRequest(http.MethodGet, "/"+shortURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusTemporaryRedirect)
	}

	redirectURL := rr.Header().Get("Location")
	if redirectURL != longURL {
		t.Errorf("handler returned wrong redirect URL: got %v want %v", redirectURL, longURL)
	}
}
