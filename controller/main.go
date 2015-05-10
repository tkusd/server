package controller

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/tkusd/server/controller/common"
	"github.com/tkusd/server/controller/v1"
)

// Router returns a http.Handler.
func Router() http.Handler {
	n := negroni.New()
	r := httprouter.New()

	r.GET("/", home)
	r.NotFound = common.NotFound
	r.HandleMethodNotAllowed = false

	n.Use(common.ClearContext())
	n.Use(common.NewLogger())
	n.Use(common.NewRecovery())
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.Use(common.CSR())
	n.Use(common.NewSubRoute("/v1", v1.Router()))
	n.UseHandler(r)

	return n
}

func home(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	common.APIResponse(res, req, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
