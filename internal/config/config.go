package config

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	DefaultPort       = "8080"
	DefaultIterations = 75000
)

type Config struct {
	GatewayBaseURL string
	Username       string
	Password       string
	Passphrase     string
	ServerAuthKey  string
	Port           string
	Iterations     int
}

func Load() (Config, error) {
	iterations := DefaultIterations
	if value := strings.TrimSpace(os.Getenv("PBKDF2_ITERATIONS")); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed <= 0 {
			return Config{}, errors.New("PBKDF2_ITERATIONS must be a positive integer")
		}
		iterations = parsed
	}

	cfg := Config{
		GatewayBaseURL: strings.TrimSpace(os.Getenv("GATEWAY_BASE_URL")),
		Username:       strings.TrimSpace(os.Getenv("USERNAME")),
		Password:       os.Getenv("PASSWORD"),
		Passphrase:     os.Getenv("PRIVATE_KEY"),
		ServerAuthKey:  strings.TrimSpace(os.Getenv("SERVER_AUTH_KEY")),
		Port:           strings.TrimSpace(os.Getenv("PORT")),
		Iterations:     iterations,
	}

	if cfg.Port == "" {
		cfg.Port = DefaultPort
	}

	missing := make([]string, 0, 5)
	if cfg.GatewayBaseURL == "" {
		missing = append(missing, "GATEWAY_BASE_URL")
	}
	if cfg.Username == "" {
		missing = append(missing, "USERNAME")
	}
	if cfg.Password == "" {
		missing = append(missing, "PASSWORD")
	}
	if cfg.Passphrase == "" {
		missing = append(missing, "PRIVATE_KEY")
	}
	if cfg.ServerAuthKey == "" {
		missing = append(missing, "SERVER_AUTH_KEY")
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}

func (c Config) HasServerAuthorization(value string) bool {
	if len(value) != len(c.ServerAuthKey) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(value), []byte(c.ServerAuthKey)) == 1
}
