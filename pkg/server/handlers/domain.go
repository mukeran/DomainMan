package handlers

import (
	"DomainMan/models"
	"DomainMan/pkg/database"
	"DomainMan/pkg/mq"
	"DomainMan/pkg/server/handlers/methods"
	"DomainMan/pkg/server/status"
	"DomainMan/pkg/whois"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

type DomainHandler struct {
}

func (h DomainHandler) Register(g *gin.RouterGroup) {
	g.GET("", h.List())
	g.POST("", h.Add())
	g.GET(":domain_id", h.Preload(), h.Show())
	g.DELETE(":domain_id", h.Preload(), h.Delete())
}

func (DomainHandler) Preload() gin.HandlerFunc {
	return func(c *gin.Context) {
		domainID := c.Param("domain_id")
		db := database.DB
		var domain models.Domain
		if v := db.Where("id = ?", domainID).First(&domain); gorm.IsRecordNotFoundError(v.Error) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status": status.NotFound,
			})
		} else if v.Error != nil {
			panic(v.Error)
		}
		c.Set("request_domain", &domain)
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
		var domains []models.Domain
		if v := db.Offset(offset).Limit(limit).Find(&domains); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  status.OK,
			"domains": domains,
		})
	}
}

func (DomainHandler) Add() gin.HandlerFunc {
	return func(c *gin.Context) {
		type req struct {
			Domains []string
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
		db := database.DB.Begin()
		defer db.RollbackUnlessCommitted()
		var addedDomains []models.Domain
		for _, domain := range in.Domains {
			domainToAdd := models.Domain{
				Name: domain,
			}
			if v := db.Create(&domainToAdd); v.Error != nil {
				panic(v.Error)
			}
			addedDomains = append(addedDomains, domainToAdd)
			_ = mq.MQ.Push(func() {
				db := database.DB.Begin()
				defer db.RollbackUnlessCommitted()
				domainToAdd.IsUpdatingWhois = true
				if v := db.Save(&domainToAdd); v.Error != nil {
					panic(v.Error)
				}
				w, err := whois.Lookup(domainToAdd.Name)
				if err != nil {
					log.Printf("Failed to fetch %v's whois information\n", domainToAdd.Name)
				} else {
					if v := db.Create(w); v.Error != nil {
						panic(v.Error)
					}
					domainToAdd.NameServer = w.NameServer
				}
				domainToAdd.IsUpdatingWhois = false
				if v := db.Save(&domainToAdd); v.Error != nil {
					panic(v.Error)
				}
				if v := db.Commit(); v.Error != nil {
					panic(v.Error)
				}
			})
		}
		if v := db.Commit(); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  status.OK,
			"domains": addedDomains,
		})
	}
}

func (DomainHandler) Show() gin.HandlerFunc {
	return func(c *gin.Context) {
		domain := c.MustGet("request_domain").(*models.Domain)
		c.JSON(http.StatusOK, gin.H{
			"status": status.OK,
			"domain": domain,
		})
	}
}

func (DomainHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		domain := c.MustGet("request_domain").(*models.Domain)
		db := database.DB
		if v := db.Delete(domain); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":            status.OK,
			"deleted_domain_id": domain.ID,
		})
	}
}
