package middlewares

import (
	"DomainMan/pkg/api/status"
	"DomainMan/pkg/database"
	"DomainMan/pkg/errors"
	"DomainMan/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func AccessControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Access-Token")
		db := database.DB
		var accessToken models.AccessToken
		if v := db.Where("token = ?", token).First(&accessToken); errors.Is(v.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status": status.AccessDenied,
			})
		} else if v.Error != nil {
			panic(v.Error)
		}
		c.Set("accessToken", &accessToken)
		c.Next()
	}
}
