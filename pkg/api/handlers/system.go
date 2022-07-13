package handlers

import (
	"DomainMan/pkg/api/status"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SystemHandler struct {
}

func (h SystemHandler) Register(g *gin.RouterGroup) {
	g.GET("ping", h.Ping())
}

func (SystemHandler) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  status.OK,
			"message": "pong",
		})
	}
}
