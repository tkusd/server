package common

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/tkusd/server/util"
)

func Logger(c *gin.Context) {
	start := time.Now()

	defer func() {
		util.Log().WithFields(logrus.Fields{
			"start":    start,
			"duration": time.Since(start),
			"method":   c.Request.Method,
			"code":     c.Writer.Status(),
			"ip":       c.ClientIP(),
		}).Info(c.Request.RequestURI)
	}()

	c.Next()
}
