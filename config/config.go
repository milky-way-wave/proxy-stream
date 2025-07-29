package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT          uint16
	STREAM_URL    string
	COOKIE_SECURE bool
}

func New() (*Config, error) {
	_ = godotenv.Load()

	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}

	portInt, err := strconv.Atoi(portStr)
	if err != nil || portInt < 0 || portInt > 65535 {
		return nil, fmt.Errorf("Invalid PORT value: %s", portStr)
	}

	streamUrl := os.Getenv("STREAM_URL")
	if streamUrl == "" {
		return nil, fmt.Errorf("Required environment variable STREAM_URL is not set")
	}

	cookieSecure := false
	cookieSecureStr := os.Getenv("COOKIE_SECURE")
	if cookieSecureStr != "" {
		cookieSecure, err = strconv.ParseBool(cookieSecureStr)
		if err != nil {
			return nil, fmt.Errorf("Invalid COOKIE_SECURE value: %s", cookieSecureStr)
		}
	}

	return &Config{
		STREAM_URL:    streamUrl,
		PORT:          uint16(portInt),
		COOKIE_SECURE: cookieSecure,
	}, nil
}
