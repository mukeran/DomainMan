package handlers

import "C"
import (
	"DomainMan/models"
	"DomainMan/pkg/database"
	"DomainMan/pkg/errors"
	"DomainMan/pkg/server/handlers/methods"
	"DomainMan/pkg/server/status"
	"DomainMan/pkg/whois"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WhoisHandler struct {
}

func (h WhoisHandler) Register(g *gin.RouterGroup) {
	g.GET("", h.List())
	g.POST("", h.Query())
	g.GET(":whois_id", h.Preload(), h.Show())
	g.DELETE(":whois_id", h.Preload(), h.Delete())
}

func (WhoisHandler) Preload() gin.HandlerFunc {
	return func(c *gin.Context) {
		whoisID := c.Param("whois_id")
		db := database.DB
		var w models.Whois
		if v := db.Where("id = ?", whoisID).First(&w); v.Error != nil {
			panic(v.Error)
		}
		c.Set("request_whois", &w)
		c.Next()
	}
}

func (WhoisHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		offset, limit, err := methods.GetPageInfo(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": status.BadParameter,
			})
			return
		}
		query := c.DefaultQuery("query", "")
		query = "%" + query + "%"
		db := database.DB
		var whoises []models.Whois
		if v := db.Where("domain_name like ? or registrant like ? or registrant_email like ? or "+
			"registrar like ?", query, query, query, query).
			Offset(offset).Limit(limit).Find(&whoises); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  status.OK,
			"whoises": whoises,
		})
	}
}

func (WhoisHandler) Query() gin.HandlerFunc {
	return func(c *gin.Context) {
		type req struct {
			Domain      string `binding:"required"`
			ForceUpdate bool   `json:"force_update,omitempty"`
		}
		var (
			w   *models.Whois
			in  req
			err error
		)
		err = c.ShouldBindJSON(&in)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": status.BadParameter,
			})
			return
		}
		if !in.ForceUpdate {
			w, err = whois.LookupWithCache(in.Domain)
		} else {
			w, err = whois.Lookup(in.Domain)
		}
		if err != nil {
			switch err {
			case errors.ErrMessageQueueFull:
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status": status.MessageQueueFull,
				})
			case errors.ErrUnsupportedSuffix:
				c.JSON(http.StatusBadRequest, gin.H{
					"status": status.UnsupportedSuffix,
				})
			default:
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status": status.ServerError,
				})
			}
			return
		}
		if w.ID == 0 {
			db := database.DB
			if v := db.Create(w); v.Error != nil {
				panic(v.Error)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"status": status.OK,
			"whois":  w,
		})
	}
}

func (WhoisHandler) Show() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.MustGet("request_whois").(*models.Whois)
		c.JSON(http.StatusOK, gin.H{
			"status": status.OK,
			"whois":  w,
		})
	}
}

func (WhoisHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.MustGet("request_whois").(*models.Whois)
		db := database.DB
		if v := db.Delete(w); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":           status.OK,
			"deleted_whois_id": w.ID,
		})
	}
}
