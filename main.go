package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"os"

	"samll-trading-back/api/config"
	"samll-trading-back/api/database"
	"samll-trading-back/api/handlers/accounts"
	"samll-trading-back/api/handlers/admin"
	"samll-trading-back/api/handlers/dashboard"
	"samll-trading-back/api/handlers/health"
	"samll-trading-back/api/handlers/trades"
	"samll-trading-back/api/handlers/users"
	"samll-trading-back/api/logger"
	"samll-trading-back/api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadConfig()

	if err := config.ValidateConfig(); err != nil {
		logger.L().Error("configuración inválida", slog.String("err", err.Error()))
		panic(err)
	}

	if err := middleware.InitJWKS(); err != nil {
		logger.L().Error("no se pudo inicializar JWKS", slog.String("err", err.Error()))
		panic(err)
	}

	database.ConnectDB()

	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	_ = r.SetTrustedProxies(nil)

	r.MaxMultipartMemory = 1 << 20

	r.Use(gin.Recovery())

	r.Use(middleware.RequestID())

	r.Use(middleware.RequestLogger())

	r.Use(middleware.SecurityHeaders())

	r.Use(middleware.RateLimit())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.GetAllowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", health.Check)

	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	{
		usersGroup := api.Group("/users")
		{
			usersGroup.GET("/me", users.GetMe)
			usersGroup.PATCH("/me", users.UpdateMyProfile)
		}

		accountsGroup := api.Group("/accounts")
		{
			accountsGroup.GET("", accounts.ListAccounts)
			accountsGroup.POST("", accounts.CreateAccount)

			accountItem := accountsGroup.Group("/:id")
			accountItem.Use(middleware.AccountOwnership())
			{
				accountItem.GET("", accounts.GetAccount)
				accountItem.PATCH("", accounts.UpdateAccount)
				accountItem.DELETE("", accounts.DeleteAccount)

				accountItem.POST("/trades", trades.CreateTrade)
				accountItem.GET("/trades", trades.GetTrades)
				accountItem.GET("/trades/:trade_id", trades.GetTradeByID)
				accountItem.PATCH("/trades/:trade_id", trades.UpdateTrade)
				accountItem.DELETE("/trades/:trade_id", trades.DeleteTrade)

				accountItem.GET("/dashboard/stats", dashboard.GetStats)
				accountItem.GET("/dashboard/calendar", dashboard.GetCalendarMetrics)
			}
		}

		adminGroup := api.Group("/admin")
		adminGroup.Use(middleware.AdminOnly())
		adminGroup.Use(middleware.AuditAdmin("user"))
		{
			adminGroup.GET("/users", admin.GetAllUsers)
			adminGroup.GET("/users/:id", admin.GetUserByID)
			adminGroup.PUT("/users/:id", admin.UpdateUser)
			adminGroup.DELETE("/users/:id", admin.DeleteUser)
		}
	}

	port := config.GetPort()

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Iniciar servidor en goroutine separada
	go func() {
		logger.L().Info("servidor iniciado", slog.String("port", port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.L().Error("error al iniciar el servidor", slog.String("err", err.Error()))
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	stop()

	logger.L().Info("apagando servidor gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.L().Error("error en el shutdown", slog.String("err", err.Error()))
	}

	logger.L().Info("servidor apagado correctamente")
}
