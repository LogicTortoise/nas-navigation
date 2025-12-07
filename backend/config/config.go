package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port       string
	DBPath     string
	ScriptDir  string
	StaticDir  string
	Debug      bool
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		DBPath:    getEnv("DB_PATH", "../data/navigation.db"),
		ScriptDir: getEnv("SCRIPT_DIR", "../scripts"),
		StaticDir: getEnv("STATIC_DIR", "../frontend/dist"),
		Debug:     getEnvBool("DEBUG", true),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return b
	}
	return defaultValue
}
