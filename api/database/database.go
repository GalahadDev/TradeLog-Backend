package database

import (
	"log"
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

	// Configuración avanzada de GORM
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("❌ Error conectando a la base de datos:", err)
	}

	// CONFIGURACIÓN DEL POOL DE CONEXIONES
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("❌ Error obteniendo instancia SQL:", err)
	}

	// Ajustes para plan gratuito/hobby
	sqlDB.SetMaxIdleConns(5)            // Mantener 5 conexiones libres listas
	sqlDB.SetMaxOpenConns(20)           // Máximo 20 conexiones simultáneas
	sqlDB.SetConnMaxLifetime(time.Hour) // Renovar conexiones cada hora

	log.Println("✅ Conexión a Base de Datos exitosa (Pool configurado)")
}

func GetDB() *gorm.DB {
	return DB
}
