package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tkusd/server/controller/common"
)

func Home(c *gin.Context) {
	common.APIResponse(c, http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
