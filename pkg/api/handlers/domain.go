package handlers

import (
	"DomainMan/pkg/api/handlers/methods"
	"DomainMan/pkg/api/status"
	"DomainMan/pkg/database"
	"DomainMan/pkg/errors"
	"DomainMan/pkg/models"
	"DomainMan/pkg/mq"
	"DomainMan/pkg/whois"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type DomainHandler struct {
}

func (h DomainHandler) Register(g *gin.RouterGroup) {
	g.GET("", h.List())
	g.POST("", h.Add())
	g.GET(":domainID", h.Preload(), h.Show())
	g.DELETE(":domainID", h.Preload(), h.Delete())
}

func (DomainHandler) Preload() gin.HandlerFunc {
	return func(c *gin.Context) {
		domainID := c.Param("domainID")
		db := database.DB
		var domain models.Domain
		if v := db.Where("id = ?", domainID).First(&domain); errors.Is(v.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status": status.NotFound,
			})
		} else if v.Error != nil {
			panic(v.Error)
		}
		c.Set("requestDomain", &domain)
		c.Next()
	}
}

func (DomainHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		offset, limit, err := methods.GetPageInfo(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": status.BadParameter,
			})
			return
		}
		db := database.DB
		var (
			domains []models.Domain
			count   int64
		)
		if v := db.Model(&models.Domain{}).Count(&count).Offset(int(offset)).Limit(int(limit)).Find(&domains); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  status.OK,
			"total":   count,
			"domains": domains,
		})
	}
}

func (DomainHandler) Add() gin.HandlerFunc {
	return func(c *gin.Context) {
		type req struct {
			Domains []string `json:"domains"`
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
		var addedDomains []models.Domain
		err = database.DB.Transaction(func(tx *gorm.DB) error {
			for _, domain := range in.Domains {
				domainToAdd := models.Domain{
					Name: domain,
				}
				if v := tx.Create(&domainToAdd); v.Error != nil {
					return v.Error
				}
				addedDomains = append(addedDomains, domainToAdd)
				_ = mq.MQ.Push(func() {
					err := database.DB.Transaction(func(tx *gorm.DB) error {
						domainToAdd.IsUpdatingWhois = true
						if v := tx.Save(&domainToAdd); v.Error != nil {
							return v.Error
						}
						w, err := whois.Lookup(domainToAdd.Name)
						if err != nil {
							log.Printf("Failed to fetch %v's whois information\n", domainToAdd.Name)
						} else {
							if v := tx.Create(w); v.Error != nil {
								return tx.Error
							}
							domainToAdd.NameServer = w.NameServer
						}
						domainToAdd.IsUpdatingWhois = false
						if v := tx.Save(&domainToAdd); v.Error != nil {
							return v.Error
						}
						return nil
					})
					if err != nil {
						panic(err)
					}
				})
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  status.OK,
			"domains": addedDomains,
		})
	}
}

func (DomainHandler) Show() gin.HandlerFunc {
	return func(c *gin.Context) {
		domain := c.MustGet("requestDomain").(*models.Domain)
		c.JSON(http.StatusOK, gin.H{
			"status": status.OK,
			"domain": domain,
		})
	}
}

func (DomainHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		domain := c.MustGet("requestDomain").(*models.Domain)
		db := database.DB
		if v := db.Delete(domain); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":          status.OK,
			"deletedDomainID": domain.ID,
		})
	}
}
