package handlers

import "C"
import (
	"DomainMan/models"
	"DomainMan/pkg/database"
	"DomainMan/pkg/server/handlers/methods"
	"DomainMan/pkg/server/status"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type SuffixHandler struct {
}

func (h SuffixHandler) Register(g *gin.RouterGroup) {
	g.GET("", h.List())
	g.POST("", h.Add())
	g.GET(":suffix_id", h.Preload(), h.Show())
	g.PATCH(":suffix_id", h.Preload(), h.Modify())
	g.DELETE(":suffix_id", h.Preload(), h.Delete())
}

func (SuffixHandler) Preload() gin.HandlerFunc {
	return func(c *gin.Context) {
		suffixID := c.Param("suffix_id")
		db := database.DB
		var suffix models.Suffix
		if v := db.Where("id = ?", suffixID).First(&suffix); gorm.IsRecordNotFoundError(v.Error) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status": status.NotFound,
			})
		} else if v.Error != nil {
			panic(v.Error)
		}
		c.Set("request_suffix", &suffix)
		c.Next()
	}
}

func (SuffixHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		offset, limit, err := methods.GetPageInfo(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": status.BadParameter,
			})
			return
		}
		query := c.DefaultQuery("query", "")
		var suffixes []models.Suffix
		db := database.DB
		if v := db.Where("name like ?", "%"+query+"%").Offset(offset).Limit(limit).Find(&suffixes); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   status.OK,
			"suffixes": suffixes,
		})
	}
}

func (SuffixHandler) Add() gin.HandlerFunc {
	return func(c *gin.Context) {
		type req struct {
			Name        string `binding:"alphanum|containsany=.-"`
			Memo        string
			Description string
			WhoisServer string `json:"whois_server" binding:"uri"`
		}
		var (
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
		suffix := models.Suffix{
			Name:        in.Name,
			Memo:        in.Memo,
			Description: in.Description,
			WhoisServer: in.WhoisServer,
		}
		db := database.DB
		if v := db.Create(&suffix); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status": status.OK,
			"suffix": suffix,
		})
	}
}

func (SuffixHandler) Show() gin.HandlerFunc {
	return func(c *gin.Context) {
		suffix := c.MustGet("request_suffix").(*models.Suffix)
		c.JSON(http.StatusOK, gin.H{
			"status": status.OK,
			"suffix": suffix,
		})
	}
}

func (SuffixHandler) Modify() gin.HandlerFunc {
	return func(c *gin.Context) {
		type req struct {
			Memo        *string
			Description *string
			WhoisServer *string `json:"whois_server" binding:"uri"`
		}
		var (
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
		suffix := c.MustGet("request_suffix").(*models.Suffix)
		if in.Memo != nil {
			suffix.Memo = *in.Memo
		}
		if in.Description != nil {
			suffix.Description = *in.Description
		}
		if in.WhoisServer != nil {
			suffix.WhoisServer = *in.WhoisServer
		}
		db := database.DB
		if v := db.Save(suffix); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status": status.OK,
			"suffix": suffix,
		})
	}
}

func (SuffixHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		suffix := c.MustGet("request_suffix").(*models.Suffix)
		db := database.DB
		if v := db.Delete(suffix); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":            status.OK,
			"deleted_suffix_id": suffix.ID,
		})
	}
}
