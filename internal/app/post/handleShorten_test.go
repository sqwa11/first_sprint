package post

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestHandleShorten(t *testing.T) {
	// Инициализируем маршрутизатор и задаем базовый URL
	router := chi.NewRouter()
	SetBaseURL("http://localhost:8080")
	router.Post("/", HandleShorten)

	// Создаем новый запрос с методом POST и телом запроса
	longURL := "https://example.com"
	body := bytes.NewBufferString(longURL)
	req, err := http.NewRequest(http.MethodPost, "/", body)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем ResponseRecorder для получения ответа
	rr := httptest.NewRecorder()

	// Вызываем хендлер через маршрутизатор
	router.ServeHTTP(rr, req)

	// Проверяем статус код ответа
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Проверяем содержимое ответа
	responseBody, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	shortURL := strings.TrimSpace(string(responseBody))
	if !strings.HasPrefix(shortURL, "http://localhost:8080/") {
		t.Errorf("handler returned unexpected body: got %v", shortURL)
	}

	// Проверяем, что longURL сохранен в URLMap под сгенерированным shortURL
	id := strings.TrimPrefix(shortURL, "http://localhost:8080/")
	savedURL, exists := URLMap[id]
	if !exists {
		t.Errorf("short URL not saved in map")
	}
	if savedURL != longURL {
		t.Errorf("saved long URL does not match: got %v want %v", savedURL, longURL)
	}
}
