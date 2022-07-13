package handlers

import "C"
import (
	"DomainMan/pkg/api/handlers/methods"
	"DomainMan/pkg/api/status"
	"DomainMan/pkg/database"
	"DomainMan/pkg/errors"
	"DomainMan/pkg/models"
	"DomainMan/pkg/whois"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WhoisHandler struct {
}

func (h WhoisHandler) Register(g *gin.RouterGroup) {
	g.GET("", h.List())
	g.POST("", h.Query())
	g.GET(":whoisID", h.Preload(), h.Show())
	g.DELETE(":whoisID", h.Preload(), h.Delete())
}

func (WhoisHandler) Preload() gin.HandlerFunc {
	return func(c *gin.Context) {
		whoisID := c.Param("whoisID")
		db := database.DB
		var w models.Whois
		if v := db.Where("id = ?", whoisID).First(&w); v.Error != nil {
			panic(v.Error)
		}
		c.Set("requestWhois", &w)
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
		var (
			whoises []models.Whois
			count   int64
		)
		if v := db.Model(&models.Whois{}).
			Where("domain_name like ? or registrant like ? or registrant_email like ? or "+
				"registrar like ?", query, query, query, query).Count(&count).
			Offset(int(offset)).Limit(int(limit)).Find(&whoises); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  status.OK,
			"total":   count,
			"whoises": whoises,
		})
	}
}

func (WhoisHandler) Query() gin.HandlerFunc {
	return func(c *gin.Context) {
		type req struct {
			Domain      string `binding:"required" json:"domain"`
			ForceUpdate bool   `json:"forceUpdate,omitempty"`
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
			case errors.ErrBadWhoisFormatOrNotRegistered:
				c.JSON(http.StatusNotFound, gin.H{
					"status": status.BadWhoisFormatOrNotRegistered,
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
		w := c.MustGet("requestWhois").(*models.Whois)
		c.JSON(http.StatusOK, gin.H{
			"status": status.OK,
			"whois":  w,
		})
	}
}

func (WhoisHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.MustGet("requestWhois").(*models.Whois)
		db := database.DB
		if v := db.Delete(w); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":         status.OK,
			"deletedWhoisID": w.ID,
		})
	}
}
