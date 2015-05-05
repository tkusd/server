package common

import (
	"net/http"

	"github.com/tkusd/server/util"
)

// HandleAPIError handles API errors.
func HandleAPIError(res http.ResponseWriter, req *http.Request, err error) {
	switch e := err.(type) {
	case *util.APIError:
		if e.Status == 0 {
			e.Status = http.StatusBadRequest
		}

		APIResponse(res, req, e.Status, e)
		break

	default:
		panic(e)
	}
}
