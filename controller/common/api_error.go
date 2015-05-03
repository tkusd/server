package common

import (
	"net/http"

	"github.com/tommy351/app-studio-server/util"
)

func HandleAPIError(res http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *util.APIError:
		if e.Status == 0 {
			e.Status = http.StatusBadRequest
		}

		RenderJSON(res, e.Status, e)
		break

	default:
		panic(e)
	}
}
