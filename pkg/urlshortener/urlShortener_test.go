package urlshortener

import (
	"testing"
)

func TestShorten(t *testing.T) {
	url := "https://example.com"
	shortURL := Shorten(url)

	if len(shortURL) != 8 {
		t.Errorf("expected shortened URL of length 8, got %d", len(shortURL))
	}
}
