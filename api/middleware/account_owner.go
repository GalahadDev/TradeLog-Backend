package middleware

import (
	"errors"
	"log"
	"net/http"

	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"
	"samll-trading-back/api/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AccountOwnership valida que la cuenta en :id pertenezca al usuario autenticado.
func AccountOwnership() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		accountID := c.Param("id")

		var account domains.TradingAccount
		db := database.GetDB()

		err := db.Select("id").
			Where("id = ? AND user_id = ? AND deleted_at IS NULL", accountID, userID).
			First(&account).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Cuenta no encontrada"})
			} else {
				log.Printf("AccountOwnership db error accountID=%s userID=%s: %v", accountID, userID, err)
				response.AbortInternalError(c)
			}
			return
		}

		c.Set("accountID", account.ID)
		c.Next()
	}
}
