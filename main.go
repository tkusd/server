package main

import (
	"net/http"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/tkusd/server/controller"
	"github.com/tkusd/server/util"
)

func main() {
	r := controller.Router()
	log := util.Log()
	addr := ":3000"

	log.Infof("Listening on %s", addr)
	gracehttp.Serve(&http.Server{
		Addr:    addr,
		Handler: r,
	})
}
