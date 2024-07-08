package config

import (
	"flag"
	"fmt"
)

// Config содержит конфигурацию сервиса
type Config struct {
	Address string
	BaseURL string
}

// NewConfig инициализирует конфигурацию из аргументов командной строки
func NewConfig() *Config {
	address := flag.String("a", "localhost:8080", "Address to run HTTP server")
	baseURL := flag.String("b", "http://localhost:8080", "Base URL for shortened URLs")

	flag.Parse()

	return &Config{
		Address: *address,
		BaseURL: *baseURL,
	}
}

func (c *Config) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("address cannot be empty")
	}
	if c.BaseURL == "" {
		return fmt.Errorf("base URL cannot be empty")
	}
	return nil
}
