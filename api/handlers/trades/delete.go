package trades

import (
	"net/http"
	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"

	"github.com/gin-gonic/gin"
)

func DeleteTrade(c *gin.Context) {
	userID := c.GetString("userID")
	tradeID := c.Param("id")

	db := database.GetDB()

	// Tenant Isolation en el DELETE:
	result := db.Delete(&domains.Trade{}, "id = ? AND user_id = ?", tradeID, userID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trade no encontrado o no te pertenece"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trade eliminado correctamente"})
}
