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

type jsonpData struct {
	Status int
	Data   interface{}
}

func (j jsonpData) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"meta": map[string]interface{}{
			"status": j.Status,
		},
		"data": j.Data,
	})
}

// RenderJSONP renders JSON-P.
func RenderJSONP(res http.ResponseWriter, status int, callback string, value interface{}) {
	res.Header().Set("Content-Type", "application/javascript")

	data := &jsonpData{
		Status: status,
		Data:   value,
	}

	if result, err := json.Marshal(data); err == nil {
		result = append([]byte(callback+"("), result...)
		result = append(result, ')')

		res.WriteHeader(http.StatusOK)
		res.Header().Set("Content-Length", strconv.Itoa(len(result)))
		res.Write(result)
	} else {
		panic(err)
	}
}

// APIResponse handles API responses.
func APIResponse(res http.ResponseWriter, req *http.Request, status int, value interface{}) {
	callback := req.URL.Query().Get("callback")

	if callback != "" {
		RenderJSONP(res, status, callback, value)
	} else {
		RenderJSON(res, status, value)
	}
}
