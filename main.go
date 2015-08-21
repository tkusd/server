package main

import (
	"net/http"
	"strconv"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/tkusd/server/config"
	"github.com/tkusd/server/controller"
)

func main() {
	r := controller.Router()
	addr := config.Config.Server.Host + ":" + strconv.Itoa(config.Config.Server.Port)

	gracehttp.Serve(&http.Server{
		Addr:    addr,
		Handler: r,
	})
}
