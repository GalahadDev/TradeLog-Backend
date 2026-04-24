package health

import (
	"net/http"

	"samll-trading-back/api/database"

	"github.com/gin-gonic/gin"
)

// Check verifica el estado del servidor y la conectividad con la base de datos.
// Devuelve 200 si todo está bien, 503 si la DB no responde.
func Check(c *gin.Context) {
	dbStatus := "ok"
	httpStatus := http.StatusOK
	overallStatus := "ok"

	sqlDB, err := database.GetDB().DB()
	if err != nil || sqlDB.PingContext(c.Request.Context()) != nil {
		dbStatus = "error"
		httpStatus = http.StatusServiceUnavailable
		overallStatus = "degraded"
	}

	c.JSON(httpStatus, gin.H{
		"status": overallStatus,
		"checks": gin.H{
			"db": dbStatus,
		},
	})
}
