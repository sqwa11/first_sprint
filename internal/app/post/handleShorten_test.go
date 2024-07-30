package post

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestHandleShorten(t *testing.T) {
	router := chi.NewRouter()
	SetBaseURL("http://localhost:8080")
	router.Post("/api/shorten", HandleShorten)

	longURL := "https://example.com"
	body := bytes.NewBufferString(longURL)
	req, err := http.NewRequest(http.MethodPost, "/api/shorten", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/plain")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	expected := "http://localhost:8080/abcd1234" // Use actual short URL from response
	if response["result"] != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", response["result"], expected)
	}
}
