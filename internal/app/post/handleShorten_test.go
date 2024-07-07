package post

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleShorten(t *testing.T) {
	// Создаем новый запрос с методом post и телом запроса
	longURL := "https://oogle.com"
	body := bytes.NewBufferString(longURL)
	req, err := http.NewRequest(http.MethodPost, "/", body)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем ResponseRecorder для получения ответа
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleShorten)

	// Вызываем хендлер
	handler.ServeHTTP(rr, req)

	// Проверяем статус код ответа
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Проверяем содержимое ответа
	responseBody, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	shortURL := strings.TrimSpace(string(responseBody))
	BaseURL := BaseURL
	if !strings.HasPrefix(shortURL, BaseURL+"/") {
		t.Errorf("handler returned unexpected body: got %v", shortURL)
	}

	// Проверяем, что longURL сохранен в urlMap под сгенерированным shortURL
	id := strings.TrimPrefix(shortURL, BaseURL+"/")
	URLMap := UrlMap
	savedURL, exists := URLMap[id]
	if !exists {
		t.Errorf("short URL not saved in map")
	}
	if savedURL != longURL {
		t.Errorf("saved long URL does not match: got %v want %v", savedURL, longURL)
	}
}
