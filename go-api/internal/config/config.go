package config

import "os"

// Config holds the application configuration loaded from environment variables.
type Config struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	ServerPort string
}

// LoadConfig reads configuration from environment variables with sensible defaults.
func LoadConfig() Config {
	return Config{
		DBHost:     getEnv("DB_SERVER", "mysql"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "task-list"),
		DBUser:     getEnv("DB_USER", "secretuser"),
		DBPassword: getEnv("DB_PASSWORD", "thisisasupersecretpassworddontyouthink"),
		ServerPort: getEnv("APP_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
