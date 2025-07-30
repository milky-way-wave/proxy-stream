package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	port          uint16
	stream_url    string
	cookie_secure bool
}

func New() (Config, error) {
	_ = godotenv.Load()

	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}

	portInt, err := strconv.Atoi(portStr)
	if err != nil || portInt < 0 || portInt > 65535 {
		return Config{}, fmt.Errorf("Invalid PORT value: %s", portStr)
	}

	streamUrl := os.Getenv("STREAM_URL")
	if streamUrl == "" {
		return Config{}, fmt.Errorf("Required environment variable STREAM_URL is not set")
	}

	cookieSecure := false
	cookieSecureStr := os.Getenv("COOKIE_SECURE")
	if cookieSecureStr != "" {
		cookieSecure, err = strconv.ParseBool(cookieSecureStr)
		if err != nil {
			return Config{}, fmt.Errorf("Invalid COOKIE_SECURE value: %s", cookieSecureStr)
		}
	}

	return Config{
		port:          uint16(portInt),
		stream_url:    streamUrl,
		cookie_secure: cookieSecure,
	}, nil
}

func (c Config) GetPort() uint16 {
	return c.port
}

func (c Config) GetStreamUrl() string {
	return c.stream_url
}

func (c Config) GetCookieSecure() bool {
	return c.cookie_secure
}
