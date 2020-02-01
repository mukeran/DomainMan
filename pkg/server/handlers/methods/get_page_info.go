package methods

import (
	"github.com/gin-gonic/gin"
)

func GetPageInfo(c *gin.Context) (offset uint64, limit uint64, err error) {
	offset, err = GetUintQuery(c, "offset", 0)
	if err != nil {
		return
	}
	limit, err = GetUintQuery(c, "limit", 10)
	if err != nil {
		return
	}
	return
}
