package admin

import (
	"net/http"
	"strconv"

	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"
	"samll-trading-back/api/response"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	db := database.GetDB()
	var users []domains.User
	var total int64

	// Filtro opcional: /api/admin/users?pending=true
	query := db.Model(&domains.User{})
	if c.Query("pending") == "true" {
		query = query.Where("is_active = ?", false)
	}

	query.Count(&total)

	if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		response.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}
