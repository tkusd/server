package common

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

const paramsKey = "params"

// WrapHandler wraps the http.Handler to fit httprouter.Handle.
func WrapHandler(handler http.Handler) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		setContextParams(req, params)
		handler.ServeHTTP(res, req)
	}
}

// WrapHandlerFunc wraps the http.HandlerFunc to fit httprouter.Handle.
func WrapHandlerFunc(handler http.HandlerFunc) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		setContextParams(req, params)
		handler(res, req)
	}
}

func setContextParams(req *http.Request, params httprouter.Params) {
	context.Set(req, paramsKey, params)
}

// GetParam returns the URL param by name.
func GetParam(req *http.Request, name string) string {
	params, ok := context.GetOk(req, paramsKey)

	if ok {
		if p, ok := params.(httprouter.Params); ok {
			return p.ByName(name)
		}
	}

	return ""
}
