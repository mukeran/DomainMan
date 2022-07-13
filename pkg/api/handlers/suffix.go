package handlers

import "C"
import (
	"DomainMan/pkg/api/handlers/methods"
	"DomainMan/pkg/api/status"
	"DomainMan/pkg/database"
	"DomainMan/pkg/errors"
	"DomainMan/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type SuffixHandler struct {
}

func (h SuffixHandler) Register(g *gin.RouterGroup) {
	g.GET("", h.List())
	g.POST("", h.Add())
	g.GET(":suffixID", h.Preload(), h.Show())
	g.PATCH(":suffixID", h.Preload(), h.Modify())
	g.DELETE(":suffixID", h.Preload(), h.Delete())
}

func (SuffixHandler) Preload() gin.HandlerFunc {
	return func(c *gin.Context) {
		suffixID := c.Param("suffixID")
		db := database.DB
		var suffix models.Suffix
		if v := db.Where("id = ?", suffixID).First(&suffix); errors.Is(v.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status": status.NotFound,
			})
		} else if v.Error != nil {
			panic(v.Error)
		}
		c.Set("requestSuffix", &suffix)
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
		var (
			suffixes []models.Suffix
			count    int64
		)
		db := database.DB
		if v := db.Model(&models.Suffix{}).Where("name like ?", "%"+query+"%").Count(&count).
			Offset(int(offset)).Limit(int(limit)).Find(&suffixes); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   status.OK,
			"total":    count,
			"suffixes": suffixes,
		})
	}
}

func (SuffixHandler) Add() gin.HandlerFunc {
	return func(c *gin.Context) {
		type req struct {
			Name        string `binding:"alphanum|containsany=.-" json:"name"`
			Memo        string `json:"memo"`
			Description string `json:"description"`
			WhoisServer string `json:"whoisServer" binding:"uri"`
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
		suffix := c.MustGet("requestSuffix").(*models.Suffix)
		c.JSON(http.StatusOK, gin.H{
			"status": status.OK,
			"suffix": suffix,
		})
	}
}

func (SuffixHandler) Modify() gin.HandlerFunc {
	return func(c *gin.Context) {
		type req struct {
			Memo        *string `json:"memo"`
			Description *string `json:"description"`
			WhoisServer *string `json:"whoisServer" binding:"uri"`
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
		suffix := c.MustGet("requestSuffix").(*models.Suffix)
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
		suffix := c.MustGet("requestSuffix").(*models.Suffix)
		db := database.DB
		if v := db.Delete(suffix); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":          status.OK,
			"deletedSuffixID": suffix.ID,
		})
	}
}
