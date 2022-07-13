package handlers

import (
	"DomainMan/pkg/api/status"
	"DomainMan/pkg/database"
	"DomainMan/pkg/errors"
	"DomainMan/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type ConfigHandler struct {
}

func (h ConfigHandler) Register(g *gin.RouterGroup) {
	g.GET("", h.Get())
	g.POST("", h.Set())
}

func (ConfigHandler) Set() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			in  map[string]string
			err error
		)
		err = c.ShouldBindJSON(&in)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": status.BadParameter,
			})
			return
		}
		var configs []models.Config
		err = database.DB.Transaction(func(tx *gorm.DB) error {
			for k, v := range in {
				var config models.Config
				if v := tx.Where(models.Config{Key: k}).FirstOrCreate(&config); v.Error != nil {
					return v.Error
				}
				config.Value = v
				if v := tx.Save(&config); v.Error != nil {
					return v.Error
				}
				configs = append(configs, config)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  status.OK,
			"configs": configs,
		})
	}
}

func (ConfigHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		keys, isSelected := c.GetQueryArray("key")
		db := database.DB
		var configs []models.Config
		if !isSelected {
			if v := db.Find(&configs); v.Error != nil {
				panic(v.Error)
			}
		} else {
			for _, key := range keys {
				var config models.Config
				if v := db.Where("key = ?", key).First(&config); errors.Is(v.Error, gorm.ErrRecordNotFound) {
					continue
				} else if v.Error != nil {
					panic(v.Error)
				}
				configs = append(configs, config)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  status.OK,
			"configs": configs,
		})
	}
}
