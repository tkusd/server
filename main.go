package main

import (
	"time"

	"github.com/tkusd/server/controller"
	"github.com/tkusd/server/util"
	"gopkg.in/tylerb/graceful.v1"
)

func main() {
	r := controller.Router()
	log := util.Log()
	addr := ":3000"

	log.Infof("Listening on %s", addr)
	graceful.Run(addr, 10*time.Second, r)
}
