package config

import (
	"fmt"
	"log"
	"os"

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
	// Si existe la variable completa
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	// Opción B: Construcción manual desde variables separadas
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Supabase Pooler
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

// Obtener la URL del proyecto
func GetSupabaseURL() string {
	url := os.Getenv("SUPABASE_URL")
	if url == "" {
		log.Fatal("❌ Error: SUPABASE_URL es requerida en el .env")
	}
	return url
}
