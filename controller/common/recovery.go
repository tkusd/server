package common

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/tkusd/server/util"
)

const (
	errorFormat = "PANIC: %s\n%s"
	stackSize   = 1024 * 8
	stackAll    = false
)

func Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			HandleAPIError(c, &util.APIError{
				Code:    util.ServerError,
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})

			// Print error stack
			stack := make([]byte, stackSize)
			stack = stack[:runtime.Stack(stack, stackAll)]

			util.Log().Errorf(errorFormat, err, stack)
		}
	}()

	c.Next()
}
