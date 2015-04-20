package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tommy351/app-studio-server/controller/v1"
	"github.com/tommy351/app-studio-server/util"
)

func RegisterRoute(r *mux.Router) {
	r.HandleFunc("/", home)
	v1.RegisterRoute(r.PathPrefix("/v1").Subrouter())
}

func home(res http.ResponseWriter, req *http.Request) {
	util.RenderJSON(res, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
