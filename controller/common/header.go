package common

import "github.com/gin-gonic/gin"

// NoCacheHeader adds no-cache data to the response header.
func NoCacheHeader(c *gin.Context) {
	c.Header("Pragma", "no-cache")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")
}
