package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	JWT      JWTConfig      `json:"jwt"`
}

type ServerConfig struct {
	Port string `json:"port"`
	Mode string `json:"mode"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbName"`
	SSLMode  string `json:"sslMode"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type JWTConfig struct {
	Secret    string `json:"secret"`
	ExpiresIn int    `json:"expiresIn"`
}

// Load loads configuration from a JSON file, then applies environment variable overrides.
// The JSON config file path is resolved in the following order:
//  1. The filePath argument (if non-empty)
//  2. The CONFIG_FILE environment variable
//  3. Default: "config.json" in the current working directory
//
// If the config file does not exist, built-in defaults are used.
// Environment variables always take precedence over file values.
func Load(filePath ...string) *Config {
	cfg := defaultConfig()

	// Determine config file path
	path := ""
	if len(filePath) > 0 && filePath[0] != "" {
		path = filePath[0]
	} else if envPath := os.Getenv("CONFIG_FILE"); envPath != "" {
		path = envPath
	} else {
		path = "config.json"
	}

	// Load from JSON file (ignore errors — file is optional)
	if err := loadFromFile(path, cfg); err != nil {
		// Only log if the file was explicitly specified (not the default)
		if (len(filePath) > 0 && filePath[0] != "") || os.Getenv("CONFIG_FILE") != "" {
			fmt.Fprintf(os.Stderr, "Warning: failed to load config file %q: %v\n", path, err)
		}
	}

	// Apply environment variable overrides
	applyEnvOverrides(cfg)

	return cfg
}

// defaultConfig returns a Config populated with sensible defaults.
func defaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8080",
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "123456",
			DBName:   "shihai",
			SSLMode:  "disable",
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     "6379",
			Password: "",
			DB:       0,
		},
		JWT: JWTConfig{
			Secret:    "your-secret-key",
			ExpiresIn: 86400,
		},
	}
}

// loadFromFile reads the JSON config file at the given path and unmarshals it into cfg.
func loadFromFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("parse config file: %w", err)
	}

	return nil
}

// applyEnvOverrides overrides config values with environment variables when they are set.
func applyEnvOverrides(cfg *Config) {
	// Server
	overrideStr(&cfg.Server.Port, "SERVER_PORT")
	overrideStr(&cfg.Server.Mode, "GIN_MODE")

	// Database
	overrideStr(&cfg.Database.Host, "DB_HOST")
	overrideStr(&cfg.Database.Port, "DB_PORT")
	overrideStr(&cfg.Database.User, "DB_USER")
	overrideStr(&cfg.Database.Password, "DB_PASSWORD")
	overrideStr(&cfg.Database.DBName, "DB_NAME")
	overrideStr(&cfg.Database.SSLMode, "DB_SSLMODE")

	// Redis
	overrideStr(&cfg.Redis.Host, "REDIS_HOST")
	overrideStr(&cfg.Redis.Port, "REDIS_PORT")
	overrideStr(&cfg.Redis.Password, "REDIS_PASSWORD")
	overrideInt(&cfg.Redis.DB, "REDIS_DB")

	// JWT
	overrideStr(&cfg.JWT.Secret, "JWT_SECRET")
	overrideInt(&cfg.JWT.ExpiresIn, "JWT_EXPIRES_IN")
}

// overrideStr sets *target to the value of the named environment variable, if it is set and non-empty.
func overrideStr(target *string, envKey string) {
	if v := os.Getenv(envKey); v != "" {
		*target = v
	}
}

// overrideInt sets *target to the integer value of the named environment variable, if it is set and valid.
func overrideInt(target *int, envKey string) {
	if v := os.Getenv(envKey); v != "" {
		if intVal, err := strconv.Atoi(v); err == nil {
			*target = intVal
		}
	}
}
