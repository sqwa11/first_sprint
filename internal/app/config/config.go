package config

import "errors"

type Config struct {
	BaseURL string
	Address string
}

func NewConfig() *Config {
	return &Config{
		BaseURL: "http://localhost:8080",
		Address: ":8080",
	}
}

func (c *Config) Validate() error {
	if c.BaseURL == "" || c.Address == "" {
		return errors.New("missing configuration")
	}
	return nil
}
