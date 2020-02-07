package handlers

import (
	"DomainMan/models"
	"DomainMan/pkg/database"
	"DomainMan/pkg/random"
	"DomainMan/pkg/server/handlers/methods"
	"DomainMan/pkg/server/status"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

const (
	AccessTokenLength = 32
)

type AccessTokenHandler struct {
}

func (h AccessTokenHandler) Register(g *gin.RouterGroup) {
	g.GET("", h.List())
	g.POST("", h.New())
	g.PATCH(":accessTokenID", h.Preload(), h.Modify())
	g.DELETE(":accessTokenID", h.Preload(), h.Delete())
}

func (AccessTokenHandler) Preload() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessTokenID := c.Param("accessTokenID")
		db := database.DB
		var accessToken models.AccessToken
		if v := db.Where("id = ?", accessTokenID).First(&accessToken); gorm.IsRecordNotFoundError(v.Error) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status": status.NotFound,
			})
		} else if v.Error != nil {
			panic(v.Error)
		}
		c.Set("requestAccessToken", &accessToken)
		c.Next()
	}
}

func (AccessTokenHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.MustGet("accessToken").(*models.AccessToken)
		offset, limit, err1 := methods.GetPageInfo(c)
		issuer, err2 := methods.GetUintQuery(c, "issuer", uint64(accessToken.ID))
		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": status.BadParameter,
			})
			return
		}
		if !accessToken.IsMaster && issuer != uint64(accessToken.ID) {
			c.JSON(http.StatusForbidden, gin.H{
				"status": status.AccessDenied,
			})
			return
		}
		db := database.DB
		var accessTokens []models.AccessToken
		query := db.Select([]string{"id", "created_at", "updated_at", "name", "is_master", "can_issue", "issuer_id"}).Offset(offset).Limit(limit)
		if issuer != 0 {
			query = query.Where("issuer_id = ?", issuer)
		}
		if v := query.Find(&accessTokens); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":        status.OK,
			"access_tokens": accessTokens,
		})
	}
}

func (AccessTokenHandler) New() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.MustGet("accessToken").(*models.AccessToken)
		if !accessToken.CanIssue {
			c.JSON(http.StatusForbidden, gin.H{
				"status": status.AccessDenied,
			})
			return
		}
		type req struct {
			Name     string `json:"name"`
			CanIssue bool   `json:"canIssue"`
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
		db := database.DB
		newAccessToken := models.AccessToken{
			Name:     in.Name,
			Token:    random.String(AccessTokenLength, random.DictAlphaNumber),
			CanIssue: in.CanIssue,
			IssuerID: accessToken.ID,
		}
		if v := db.Create(&newAccessToken); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":       status.OK,
			"access_token": newAccessToken,
		})
	}
}

func (AccessTokenHandler) Modify() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.MustGet("accessToken").(*models.AccessToken)
		requestAccessToken := c.MustGet("requestAccessToken").(*models.AccessToken)
		if !(accessToken.IsMaster || (accessToken.CanIssue && accessToken.ID == requestAccessToken.IssuerID)) {
			c.JSON(http.StatusForbidden, gin.H{
				"status": status.AccessDenied,
			})
			return
		}
		type req struct {
			Name     *string `json:"name,omitempty"`
			IsMaster *bool   `json:"isMaster,omitempty"`
			CanIssue *bool   `json:"canIssue,omitempty"`
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
		db := database.DB
		if in.Name != nil {
			requestAccessToken.Name = *in.Name
		}
		if accessToken.IsMaster && in.IsMaster != nil {
			requestAccessToken.IsMaster = *in.IsMaster
		}
		if in.CanIssue != nil {
			requestAccessToken.CanIssue = *in.CanIssue
		}
		if v := db.Save(requestAccessToken); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status": status.OK,
			"access_token": func() (r models.AccessToken) {
				r = *requestAccessToken
				r.Token = ""
				return
			}(),
		})
	}
}

func (AccessTokenHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.MustGet("accessToken").(*models.AccessToken)
		requestAccessToken := c.MustGet("requestAccessToken").(*models.AccessToken)
		if !(accessToken.IsMaster || (accessToken.CanIssue && accessToken.ID == requestAccessToken.IssuerID)) {
			c.JSON(http.StatusForbidden, gin.H{
				"status": status.AccessDenied,
			})
			return
		}
		db := database.DB
		if v := db.Delete(requestAccessToken); v.Error != nil {
			panic(v.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":               status.OK,
			"deletedAccessTokenID": requestAccessToken.ID,
		})
	}
}
