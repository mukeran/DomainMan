package handlers

import (
	"DomainMan/models"
	"DomainMan/pkg/database"
	"DomainMan/pkg/server/handlers/methods"
	"DomainMan/pkg/server/status"
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	IANARegistrarIDCSV = "https://www.iana.org/assignments/registrar-ids/registrar-ids-1.csv"
)

type RegistrarHandler struct {
}

func (h RegistrarHandler) Register(g *gin.RouterGroup) {
	g.GET("", h.List())
	g.PATCH("", h.Update())
}

func (RegistrarHandler) Preload() gin.HandlerFunc {
	return func(c *gin.Context) {
		registrarID := c.Param("registrar_id")
		db := database.DB
		var registrar models.Registrar
		if v := db.Where("id = ?", registrarID).First(&registrar); gorm.IsRecordNotFoundError(v.Error) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status": status.NotFound,
			})
		} else if v.Error != nil {
			panic(v.Error)
		}
		c.Next()
	}
}

func (RegistrarHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		offset, limit, err := methods.GetPageInfo(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": status.BadParameter,
			})
			return
		}
		db := database.DB
		var registrars []models.Registrar
		if v := db.Offset(offset).Limit(limit).Find(&registrars); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":     status.OK,
			"registrars": registrars,
		})
	}
}

func (RegistrarHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := http.Get(IANARegistrarIDCSV)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": status.ConnectionError,
			})
			return
		}
		r := csv.NewReader(resp.Body)
		_, err = r.Read()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": status.FormatError,
			})
			return
		}
		fetchedAt := time.Now()
		db := database.DB.Begin()
		defer db.RollbackUnlessCommitted()
		for {
			line, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status": status.FormatError,
				})
				return
			}
			ianaID, _ := strconv.ParseUint(line[0], 10, 64)
			registrar := models.Registrar{
				Name:        line[1],
				IANAID:      uint(ianaID),
				Status:      models.RegistrarStatusFromString[strings.ToLower(line[2])],
				RDAPBaseURL: line[3],
				IsFromIANA:  true,
				FetchedAt:   fetchedAt,
			}
			if v := db.Create(&registrar); v.Error != nil {
				panic(v.Error)
			}
		}
		if v := db.Commit(); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status": status.OK,
		})
	}
}
