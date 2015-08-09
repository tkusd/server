package common

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RenderJSON renders JSON and add corresponding header to the response.
func RenderJSON(c *gin.Context, status int, value interface{}) error {
	if result, err := json.Marshal(value); err == nil {
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Writer.WriteHeader(status)
		c.Writer.Write(result)
	} else {
		return err
	}

	return nil
}

type jsonpData struct {
	Meta *jsonpMeta   `json:"meta"`
	Data *interface{} `json:"data"`
}

type jsonpMeta struct {
	Status int `json:"status"`
}

// RenderJSONP renders JSON-P.
func RenderJSONP(c *gin.Context, status int, callback string, value interface{}) error {
	data := &jsonpData{
		Meta: &jsonpMeta{
			Status: status,
		},
		Data: &value,
	}

	if result, err := json.Marshal(data); err == nil {
		buf := bytes.Buffer{}

		buf.WriteString(callback)
		buf.WriteRune('(')
		buf.Write(result)
		buf.WriteRune(')')

		c.Header("Content-Type", "application/javascript; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		buf.WriteTo(c.Writer)
	} else {
		return err
	}

	return nil
}

// APIResponse handles API responses.
func APIResponse(c *gin.Context, status int, value interface{}) error {
	callback := c.Query("callback")

	if callback != "" {
		return RenderJSONP(c, status, callback, value)
	} else {
		return RenderJSON(c, status, value)
	}
}
