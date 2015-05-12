package common

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

const paramsKey = "params"

type Handle func(res http.ResponseWriter, req *http.Request) error

type handleForNegroni struct {
	handle Handle
}

func (h *handleForNegroni) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if err := h.handle(res, req); err == nil {
		next(res, req)
	} else {
		HandleAPIError(res, req, err)
	}
}

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

func WrapCommonHandle(handler Handle) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		setContextParams(req, params)

		if err := handler(res, req); err != nil {
			HandleAPIError(res, req, err)
		}
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
