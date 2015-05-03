package controller

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/tommy351/app-studio-server/controller/common"
	"github.com/tommy351/app-studio-server/controller/v1"
)

// Router returns a http.Handler.
func Router() http.Handler {
	n := negroni.New()
	r := httprouter.New()

	r.GET("/", home)
	r.NotFound = common.NotFound

	n.Use(common.ClearContext())
	n.Use(common.NewLogger())
	n.Use(common.NewRecovery())
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.Use(common.NewSubRoute("/v1", v1.Router()))
	n.UseHandler(r)

	return n
}

func home(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	common.RenderJSON(res, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
