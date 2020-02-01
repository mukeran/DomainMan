package middlewares

import (
	"DomainMan/models"
	"DomainMan/pkg/database"
	"DomainMan/pkg/server/status"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func AccessControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Access-Token")
		db := database.DB.Begin()
		defer db.RollbackUnlessCommitted()
		var accessToken models.AccessToken
		if v := db.Where("token = ?", token).First(&accessToken); gorm.IsRecordNotFoundError(v.Error) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status": status.AccessDenied,
			})
		} else if v.Error != nil {
			panic(v.Error)
		}
		if v := db.Commit(); v.Error != nil {
			panic(v.Error)
		}
		c.Set("access_token", &accessToken)
		c.Next()
	}
}
