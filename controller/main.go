package controller

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/tommy351/app-studio-server/controller/v1"
	"github.com/tommy351/app-studio-server/util"
)

func Router() http.Handler {
	n := negroni.New()
	r := mux.NewRouter()

	r.HandleFunc("/", home)
	r.PathPrefix("/v1").Handler(v1.Router())
	n.UseHandler(r)

	return n
}

func home(res http.ResponseWriter, req *http.Request) {
	util.RenderJSON(res, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
