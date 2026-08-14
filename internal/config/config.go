package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Application ApplicationConfig
	Database    DatabaseConfig
	Server      ServerConfig
	Auth        AuthConfig
}

type ApplicationConfig struct {
	Name        string
	Environment string
}

type DatabaseConfig struct {
	URL      string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

type ServerConfig struct {
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type AuthConfig struct {
	JWT           JWTConfig
	RefreshToken  RefreshTokenConfig
	RefreshCookie RefreshCookieConfig
}

type JWTConfig struct {
	Secret   string
	Issuer   string
	Audience string
	TTL      time.Duration
}

type RefreshTokenConfig struct {
	TTL time.Duration
}

type RefreshCookieConfig struct {
	Name     string
	Path     string
	Domain   string
	Secure   bool
	HTTPOnly bool
	SameSite string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Application: ApplicationConfig{
			Name:        getenvOrDefault("APP_NAME", "neurosis-go-api"),
			Environment: getenvOrDefault("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			URL:      os.Getenv("DATABASE_URL"),
			Host:     getenvOrDefault("DB_HOST", "localhost"),
			Port:     getenvOrDefault("DB_PORT", "5432"),
			Name:     getenvOrDefault("DB_NAME", "neurosis"),
			User:     getenvOrDefault("DB_USER", "postgres"),
			Password: getenvOrDefault("DB_PASSWORD", "postgres"),
			SSLMode:  getenvOrDefault("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Host:            getenvOrDefault("SERVER_HOST", ""),
			Port:            getenvOrDefault("SERVER_PORT", "8080"),
			ReadTimeout:     durationOrDefault("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    durationOrDefault("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:     durationOrDefault("SERVER_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: durationOrDefault("SERVER_SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		Auth: AuthConfig{
			JWT: JWTConfig{
				Secret:   os.Getenv("JWT_SECRET"),
				Issuer:   getenvOrDefault("JWT_ISSUER", "neurosis-go-api"),
				Audience: getenvOrDefault("JWT_AUDIENCE", "neurosis-api"),
				TTL:      durationOrDefault("JWT_TTL", 15*time.Minute),
			},
			RefreshToken: RefreshTokenConfig{
				TTL: durationOrDefault("REFRESH_TOKEN_TTL", 30*24*time.Hour),
			},
			RefreshCookie: RefreshCookieConfig{
				Name:     getenvOrDefault("REFRESH_COOKIE_NAME", "refresh_token"),
				Path:     getenvOrDefault("REFRESH_COOKIE_PATH", "/auth"),
				Domain:   os.Getenv("REFRESH_COOKIE_DOMAIN"),
				Secure:   boolOrDefault("REFRESH_COOKIE_SECURE", false),
				HTTPOnly: boolOrDefault("REFRESH_COOKIE_HTTP_ONLY", true),
				SameSite: strings.ToLower(getenvOrDefault("REFRESH_COOKIE_SAME_SITE", "lax")),
			},
		},
	}
}

func (c DatabaseConfig) DSN() string {
	if c.URL != "" {
		return c.URL
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}

func (c ServerConfig) Address() string {
	if c.Host == "" {
		return ":" + c.Port
	}
	return c.Host + ":" + c.Port
}

func getenvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func durationOrDefault(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func boolOrDefault(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}
