package main

import (
	"net/http"

	"github.com/tommy351/app-studio-server/controller"
	"github.com/tommy351/app-studio-server/util"
)

func main() {
	r := controller.Router()
	log := util.Log()
	addr := ":3000"

	log.Infof("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
