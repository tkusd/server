package main

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/tommy351/app-studio-server/controller"
)

func main() {
	n := negroni.Classic()
	r := mux.NewRouter()

	// Middlewares
	n.Use(gzip.Gzip(gzip.DefaultCompression))

	// Register routes
	controller.RegisterRoute(r)
	n.UseHandler(r)

	// Start listening
	n.Run(":3000")
}
