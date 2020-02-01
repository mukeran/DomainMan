package handlers

import (
	"DomainMan/models"
	"DomainMan/pkg/database"
	"DomainMan/pkg/server/status"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
		db := database.DB.Begin()
		defer db.RollbackUnlessCommitted()
		var configs []models.Config
		for k, v := range in {
			var config models.Config
			if v := db.Where(models.Config{Key: k}).FirstOrCreate(&config); v.Error != nil {
				panic(v.Error)
			}
			config.Value = v
			if v := db.Save(&config); v.Error != nil {
				panic(v.Error)
			}
			configs = append(configs, config)
		}
		if v := db.Commit(); v.Error != nil {
			panic(v.Error)
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
				if v := db.Where("key = ?", key).First(&config); gorm.IsRecordNotFoundError(v.Error) {
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
