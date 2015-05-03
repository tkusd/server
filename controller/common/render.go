package common

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// RenderJSON renders JSON and add corresponding header to the response.
func RenderJSON(res http.ResponseWriter, status int, value interface{}) {
	res.Header().Set("Content-Type", "application/json")

	if result, err := json.Marshal(value); err == nil {
		res.WriteHeader(status)
		res.Header().Set("Content-Length", strconv.Itoa(len(result)))
		res.Write(result)
	} else {
		panic(err)
	}
}
