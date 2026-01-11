package main

import (
	"log"
	"samll-trading-back/api/config"
	"samll-trading-back/api/database"
	"samll-trading-back/api/handlers/admin"
	"samll-trading-back/api/handlers/dashboard"
	"samll-trading-back/api/handlers/health"
	"samll-trading-back/api/handlers/trades"
	"samll-trading-back/api/handlers/users"
	"samll-trading-back/api/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// Cargar variables de entorno (.env)
	config.LoadConfig()

	// Conectar a PostgreSQL (GORM)
	database.ConnectDB()

	// Inicializar Router de Gin
	r := gin.Default()

	// Configuración de CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://tradelog-app.vercel.app","https://cron-job.org"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health Check: Para verificar que el servidor está vivo (Ping)
	r.GET("/health", health.Check)

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		//  MÓDULO DE USUARIOS (Self-Service)
		usersGroup := api.Group("/users")
		{
			// Obtener perfil del usuario logueado (Rol, Status, Bio)
			usersGroup.GET("/me", users.GetMe)
			// Editar perfil propio (Bio, Teléfono, Foto, etc.)
			usersGroup.PATCH("/me", users.UpdateMyProfile)
		}

		// DASHBOARD & ANALYTICS
		dashboardGroup := api.Group("/dashboard")
		{
			// Datos para el Heatmap/Calendario (PnL diario)
			dashboardGroup.GET("/calendar", dashboard.GetCalendarMetrics)
			// KPIs Avanzados (WinRate, Profit Factor, Sharpe Ratio, Expectancy)
			dashboardGroup.GET("/stats", dashboard.GetStats)
		}

		//  MÓDULO DE TRADES (Core del Trading Journal)
		tradesGroup := api.Group("/trades")
		{
			// Crear una nueva operación
			tradesGroup.POST("", trades.CreateTrade)
			// Listar operaciones (Soporta paginación ?page=1&limit=10)
			tradesGroup.GET("", trades.GetTrades)
			// Obtener detalle de una operación específica
			tradesGroup.GET("/:id", trades.GetTradeByID)
			// Editar operación (Notas, Precios, Status)
			tradesGroup.PATCH("/:id", trades.UpdateTrade)
			// Eliminar operación (Soft Delete o Hard Delete según implementación)
			tradesGroup.DELETE("/:id", trades.DeleteTrade)
		}

		// PANEL DE ADMINISTRACIÓN (Backoffice)
		// Requiere Rol 'admin' además de estar autenticado
		adminGroup := api.Group("/admin")
		adminGroup.Use(middleware.AdminOnly())
		{
			// Listar todos los usuarios registrados
			adminGroup.GET("/users", admin.GetAllUsers)
			// Ver detalle completo de un usuario
			adminGroup.GET("/users/:id", admin.GetUserByID)
			// Edición total (Aprobar/Banear, cambiar Roles, editar info)
			adminGroup.PUT("/users/:id", admin.UpdateUser)
			// Eliminar usuario de la base de datos
			adminGroup.DELETE("/users/:id", admin.DeleteUser)
		}
	}

	port := config.GetPort()

	log.Println("---------------------------------------------------------")
	log.Printf("🚀 Servidor de TradeLog corriendo en el puerto:%s", port)
	log.Println("---------------------------------------------------------")

	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Error al iniciar el servidor:", err)
	}
}


