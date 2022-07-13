package methods

import (
	"github.com/gin-gonic/gin"
)

func GetPageInfo(c *gin.Context) (offset int64, limit int64, err error) {
	offset, err = GetIntQuery(c, "offset", 0)
	if err != nil {
		return
	}
	limit, err = GetIntQuery(c, "limit", 10)
	if err != nil {
		return
	}
	return
}
