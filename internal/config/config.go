package config

import (
	"log"
	"os"

	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv             string
	Port               string
	MongoURI           string
	MongoDB            string
	JWTSecret          string
	BrevoAPIKey        string
	EmailFrom          string
	EmailFromName      string
	SupabaseURL        string
	SupabaseServiceKey string
	SupabaseBucket     string
}

func LoadConfig() (*Config, error) {
	// This line is SAFE for both development and production
	err := godotenv.Load()
	if err != nil {
		// In production, this is normal - no .env file
		log.Printf("Note: .env file not found (this is normal in production)")
	} else {
		// In development, this runs
		log.Printf("✅ Loaded .env file for local development")
	}
	cfg := &Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		Port:               getEnv("PORT", "8080"),
		MongoURI:           getEnv("MONGO_URI", "mongodb://localhost:27018"),
		MongoDB:            getEnv("MONGO_DB", "fmiis"),
		JWTSecret:          getEnv("JWT_SECRET", "change-me"),
		BrevoAPIKey:        getEnv("BREVO_API_KEY", ""),
		EmailFrom:          getEnv("EMAIL_FROM", "fmiis.app@gmail.com"),
		EmailFromName:      getEnv("EMAIL_FROM_NAME", "FMIIS"),
		SupabaseURL:        getEnv("SUPABASE_URL", ""),
		SupabaseServiceKey: getEnv("SUPABASE_SERVICE_KEY", ""),
		SupabaseBucket:     getEnv("SUPABASE_BUCKET", ""),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return strings.TrimSpace(v)
	}
	return fallback
}
