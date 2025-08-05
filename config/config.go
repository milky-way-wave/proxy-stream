package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	app_secret    string
	app_port      uint16
	stream_url    string
	cookie_secure bool
}

func New() (Config, error) {
	_ = godotenv.Load()

	app_secret := os.Getenv("APP_SECRET")
	if app_secret == "" {
		return Config{}, fmt.Errorf("Required environment variable APP_SECRET is not set")
	}

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
		app_port:      uint16(portInt),
		stream_url:    streamUrl,
		cookie_secure: cookieSecure,
		app_secret:    app_secret,
	}, nil
}

func (c Config) GetPort() uint16 {
	return c.app_port
}

func (c Config) GetStreamUrl() string {
	return c.stream_url
}

func (c Config) GetCookieSecure() bool {
	return c.cookie_secure
}

func (c Config) GetAppSecret() string {
	return c.app_secret
}
