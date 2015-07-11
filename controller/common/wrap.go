package common

import "github.com/gin-gonic/gin"

type HandlerFuncWithError func(c *gin.Context) error

func Wrap(handler HandlerFuncWithError) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := handler(c); err != nil {
			HandleAPIError(c, err)
		}
	}
}
