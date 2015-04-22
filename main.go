package main

import (
	"github.com/codegangsta/negroni"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/tommy351/app-studio-server/controller"
)

func main() {
	n := negroni.Classic()

	// Middlewares
	n.Use(gzip.Gzip(gzip.DefaultCompression))

	// Register routes
	n.UseHandler(controller.Router())

	// Start listening
	n.Run(":3000")
}
