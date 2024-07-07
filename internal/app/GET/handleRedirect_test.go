package GET

import (
	"github.com/sqwa11/first_sprint/internal/app/POST"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRedirect(t *testing.T) {
	// Создаем новый сокращенный URL и сохраняем его в URLMap
	shortURL := "abcd1234"
	longURL := "https://google.com"

	URLMap := POST.UrlMap
	URLMap[shortURL] = longURL

	// Создаем новый запрос с методом GET
	req, err := http.NewRequest(http.MethodGet, "/"+shortURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем ResponseRecorder для получения ответа
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleRedirect)

	// Вызываем хендлер
	handler.ServeHTTP(rr, req)

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
