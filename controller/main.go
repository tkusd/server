package controller

import (
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/tkusd/server/controller/common"
	"github.com/tkusd/server/controller/v1"
	"github.com/tkusd/server/util"
	"github.com/tommy351/gin-cors"
)

func Router() *gin.Engine {
	g := gin.New()

	g.Use(common.Recovery)
	g.Use(common.Logger)
	g.Use(gzip.Gzip(gzip.DefaultCompression))
	g.Use(cors.Middleware(cors.Options{}))
	g.GET("/", Home)
	v1.Router(g.Group("/v1"))
	g.Static("/uploads/assets", util.GetUploadDir())
	g.NoRoute(common.NotFound)

	return g
}
