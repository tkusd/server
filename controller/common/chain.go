package common

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
)

// ChainHandler chains middlewares and returns a handler.
func ChainHandler(handlers ...interface{}) httprouter.Handle {
	n := negroni.New()

	for _, handler := range handlers {
		switch h := handler.(type) {
		case http.HandlerFunc:
		case func(http.ResponseWriter, *http.Request):
			n.UseHandlerFunc(h)
			break

		case http.Handler:
			n.UseHandler(h)
			break

		case negroni.HandlerFunc:
		case func(http.ResponseWriter, *http.Request, http.HandlerFunc):
			n.UseFunc(h)
			break

		case negroni.Handler:
			n.Use(h)
			break
		}
	}

	return WrapHandler(n)
}
