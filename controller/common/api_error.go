package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tkusd/server/util"
)

// HandleAPIError handles API errors.
func HandleAPIError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *util.APIError:
		if e.Status == 0 {
			e.Status = http.StatusBadRequest
		}

		if err := APIResponse(c, e.Status, e); err != nil {
			panic(err)
		}
		break

	default:
		panic(e)
	}
}
