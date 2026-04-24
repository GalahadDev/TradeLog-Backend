package accounts

import (
	"log"
	"net/http"

	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"
	"samll-trading-back/api/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DeleteAccount(c *gin.Context) {
	accountID := c.GetString("accountID")
	db := database.GetDB()

	err := db.Transaction(func(tx *gorm.DB) error {
		// Soft delete en cascada de todos los trades de la cuenta
		if err := tx.Where("account_id = ?", accountID).
			Delete(&domains.Trade{}).Error; err != nil {
			return err
		}
		// Soft delete de la cuenta
		return tx.Delete(&domains.TradingAccount{}, "id = ?", accountID).Error
	})

	if err != nil {
		log.Printf("DeleteAccount tx error accountID=%s: %v", accountID, err)
		response.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cuenta eliminada correctamente"})
}
