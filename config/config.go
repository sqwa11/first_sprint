package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress string
	BaseURL       string
}

func NewConfig() *Config {
	config := &Config{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080",
	}

	flag.StringVar(&config.ServerAddress, "a", "", "HTTP server address")
	flag.StringVar(&config.BaseURL, "b", "", "Base URL for shortened links")
	flag.Parse()

	if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
		config.ServerAddress = envServerAddress
	} else if flag.Lookup("a").Value.String() != "" {
		config.ServerAddress = flag.Lookup("a").Value.String()
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		config.BaseURL = envBaseURL
	} else if flag.Lookup("b").Value.String() != "" {
		config.BaseURL = flag.Lookup("b").Value.String()
	}

	return config
}
