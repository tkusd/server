package common

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
)

type clearContext struct{}

func (m *clearContext) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	defer context.Clear(req)
	next(res, req)
}

// ClearContext creates a middleware which clears the context at the end of a request.
func ClearContext() negroni.Handler {
	return &clearContext{}
}
