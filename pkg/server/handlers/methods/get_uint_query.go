package methods

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetUintQuery(c *gin.Context, key string, defaultValue uint64) (uint64, error) {
	return strconv.ParseUint(c.DefaultQuery(key, strconv.FormatUint(defaultValue, 10)), 10, 64)
}
