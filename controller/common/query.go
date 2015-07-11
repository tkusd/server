package common

import "github.com/gin-gonic/gin"

func QueryExist(c *gin.Context, param string) bool {
	if _, ok := c.Request.URL.Query()[param]; ok {
		return true
	}

	return false
}
