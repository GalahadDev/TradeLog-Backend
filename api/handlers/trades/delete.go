package trades

import (
	"net/http"

	"samll-trading-back/api/cache"
	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"

	"github.com/gin-gonic/gin"
)

func DeleteTrade(c *gin.Context) {
	accountID := c.GetString("accountID")
	tradeID := c.Param("trade_id")

	db := database.GetDB()

	result := db.Delete(&domains.Trade{}, "id = ? AND account_id = ?", tradeID, accountID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trade no encontrado o no te pertenece"})
		return
	}

	cache.Stats.Invalidate(accountID)
	c.JSON(http.StatusOK, gin.H{"message": "Trade eliminado correctamente"})
}
