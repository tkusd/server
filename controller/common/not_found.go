package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tkusd/server/util"
)

func NotFound(c *gin.Context) {
	HandleAPIError(c, &util.APIError{
		Code:    util.NotFoundError,
		Status:  http.StatusNotFound,
		Message: "Not found",
	})
}
