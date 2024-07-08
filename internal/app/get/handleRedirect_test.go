package get

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sqwa11/first_sprint/internal/app/post"
)

func TestHandleRedirect(t *testing.T) {
	// Инициализируем маршрутизатор и задаем базовый URL
	router := chi.NewRouter()
	post.SetBaseURL("http://localhost:8080")
	router.Get("/{id}", HandleRedirect)

	// Создаем новый сокращенный URL и сохраняем его в URLMap
	shortURL := "abcd1234"
	longURL := "https://example.com"
	post.URLMap[shortURL] = longURL

	// Создаем новый запрос с методом GET
	req, err := http.NewRequest(http.MethodGet, "/"+shortURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем ResponseRecorder для получения ответа
	rr := httptest.NewRecorder()

	// Вызываем хендлер через маршрутизатор
	router.ServeHTTP(rr, req)

	// Проверяем статус код ответа
	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusTemporaryRedirect)
	}

	// Проверяем заголовок Location
	redirectURL := rr.Header().Get("Location")
	if redirectURL != longURL {
		t.Errorf("handler returned wrong redirect URL: got %v want %v", redirectURL, longURL)
	}
}
