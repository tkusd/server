package common

import (
	"net/http"

	"github.com/tommy351/app-studio-server/util"
)

func NotFound(res http.ResponseWriter, req *http.Request) {
	HandleAPIError(res, &util.APIError{
		Code:    util.NotFoundError,
		Status:  http.StatusNotFound,
		Message: "Not found",
	})
}
