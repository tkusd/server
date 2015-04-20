package util

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func RenderJSON(res http.ResponseWriter, status int, value interface{}) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	if result, err := json.Marshal(value); err == nil {
		res.WriteHeader(status)
		res.Header().Set("Content-Length", strconv.Itoa(len(result)))
		res.Write(result)
	} else {
		HandleAPIError(res, &APIError{
			Code:    ServerError,
			Message: "Server error",
			Status:  http.StatusInternalServerError,
		})
	}
}
