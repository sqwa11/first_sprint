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
	// Создаем временный HTTP сервер для тестов
	ts := httptest.NewServer(http.HandlerFunc(HandleShorten))
	defer ts.Close()

	// Создаем фиктивную конфигурацию для теста
	cfg := &config.Config{
		ServerAddress: ts.URL,
		BaseURL:       ts.URL,
	}
	SetBaseURL(cfg.BaseURL)

	// Запрос к временному серверу
	longURL := "https://example.com"
	body := strings.NewReader(longURL)
	resp, err := http.Post(ts.URL, "text/plain", body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Проверяем корректность ответа
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", resp.StatusCode, http.StatusCreated)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	shortURL := strings.TrimSpace(string(responseBody))
	// Проверяем, что короткий URL начинается с ожидаемого базового URL
	if !strings.HasPrefix(shortURL, cfg.BaseURL+"/") {
		t.Errorf("handler returned unexpected body: got %v", shortURL)
	}

	// Извлекаем идентификатор из короткого URL
	id := strings.TrimPrefix(shortURL, cfg.BaseURL+"/")

	// Проверяем, что короткий URL сохранен в URLMap
	savedURL, exists := URLMap[id]
	if !exists {
		t.Errorf("short URL not saved in map")
	}
	if savedURL != longURL {
		t.Errorf("saved long URL does not match: got %v want %v", savedURL, longURL)
	}
}
