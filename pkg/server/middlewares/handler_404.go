package middlewares

import (
	"DomainMan/pkg/server/status"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Handler404() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": status.NotFound,
		})
		c.Next()
	}
}
