package config

import (
	"os"
)

type Config struct {
	AppEnv    string
	Port      string
	MongoURI  string
	MongoDB   string
	JWTSecret string
}

func LoadConfig() (*Config, error) {

	cfg := &Config{
		AppEnv:    getEnv("APP_ENV", "development"),
		Port:      getEnv("PORT", "8080"),
		MongoURI:  getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:   getEnv("MONGO_DB", "fmiis"),
		JWTSecret: getEnv("JWT_SECRET", "change-me"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
