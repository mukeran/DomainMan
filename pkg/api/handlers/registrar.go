package handlers

import (
	"DomainMan/pkg/api/handlers/methods"
	"DomainMan/pkg/api/status"
	"DomainMan/pkg/database"
	"DomainMan/pkg/errors"
	"DomainMan/pkg/models"
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		registrarID := c.Param("registrarID")
		db := database.DB
		var registrar models.Registrar
		if v := db.Where("id = ?", registrarID).First(&registrar); errors.Is(v.Error, gorm.ErrRecordNotFound) {
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
		if v := db.Offset(int(offset)).Limit(int(limit)).Find(&registrars); v.Error != nil {
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
		tx := database.DB.Begin()
		for {
			line, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				tx.Rollback()
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
			if v := tx.Create(&registrar); tx.Error != nil {
				tx.Rollback()
				panic(v.Error)
			}
		}
		if v := tx.Commit(); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status": status.OK,
		})
	}
}
