package middlewares

import (
	"DomainMan/pkg/server/status"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Internal server error")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"status": status.ServerError,
				})
			}
		}()
		c.Next()
	}
}
