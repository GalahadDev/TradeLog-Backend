package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Nota: No se encontró archivo .env, asumiendo variables de entorno del sistema.")
	}
}

// GetDBURL construye el Data Source Name (DSN) para GORM
func GetDBURL() string {

	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		host, user, password, dbname, port)

	return dsn
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8080"
	}
	return port
}

// GetAllowedOrigins retorna los orígenes CORS permitidos desde ALLOWED_ORIGINS (separados por coma).
// Si no está definida, usa valores seguros por defecto según APP_ENV.
func GetAllowedOrigins() []string {
	if raw := os.Getenv("ALLOWED_ORIGINS"); raw != "" {
		origins := strings.Split(raw, ",")
		for i, o := range origins {
			origins[i] = strings.TrimSpace(o)
		}
		return origins
	}

	defaults := []string{"https://tradelog-app.vercel.app", "https://cron-job.org"}
	if os.Getenv("APP_ENV") != "production" {
		defaults = append(defaults, "http://localhost:8081")
	}
	return defaults
}

func ValidateConfig() error {
	required := []string{
		"SUPABASE_URL",
		"JWT_SECRET",
	}

	if os.Getenv("DATABASE_URL") == "" {
		required = append(required, "DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_PORT")
	}
	for _, key := range required {
		if os.Getenv(key) == "" {
			return fmt.Errorf("variable de entorno requerida: %s", key)
		}
	}
	return nil
}

func GetSupabaseURL() string {
	return os.Getenv("SUPABASE_URL")
}
