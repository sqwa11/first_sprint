package post

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestHandleShorten(t *testing.T) {
	router := chi.NewRouter()
	SetBaseURL("http://localhost:8080")
	router.Post("/api/shorten", HandleShorten)

	longURL := "https://example.com"
	requestBody := map[string]string{"url": longURL}
	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(response["result"], "http://localhost:8080/") {
		t.Errorf("handler returned unexpected body: got %v want %v", response["result"], "http://localhost:8080/")
	}
}
