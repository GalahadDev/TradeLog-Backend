package database

import (
	"log"
	"os"
	"time"

	"samll-trading-back/api/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dsn := config.GetDBURL()

	// Nivel de logging
	logLevel := logger.Silent
	if os.Getenv("APP_ENV") != "production" {
		logLevel = logger.Info
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		log.Fatal("❌ Error conectando a la base de datos:", err)
	}

	// CONFIGURACIÓN DEL POOL DE CONEXIONES
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("❌ Error obteniendo instancia SQL:", err)
	}

	// Ajustes de Supabase
	sqlDB.SetMaxIdleConns(5)                  // Mantener 5 conexiones libres listas
	sqlDB.SetMaxOpenConns(20)                 // Máximo 20 conexiones simultáneas
	sqlDB.SetConnMaxLifetime(time.Hour)       // Renovar conexiones cada hora
	sqlDB.SetConnMaxIdleTime(5 * time.Minute) // Cerrar conexiones idle >5 min para no agotar el pool

	log.Println("✅ Conexión a Base de Datos exitosa (Pool configurado)")
}

func GetDB() *gorm.DB {
	return DB
}
