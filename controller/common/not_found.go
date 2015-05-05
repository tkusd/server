package common

import (
	"net/http"

	"github.com/tkusd/server/util"
)

// NotFound returns not found error to users.
func NotFound(res http.ResponseWriter, req *http.Request) {
	HandleAPIError(res, req, &util.APIError{
		Code:    util.NotFoundError,
		Status:  http.StatusNotFound,
		Message: "Not found",
	})
}
