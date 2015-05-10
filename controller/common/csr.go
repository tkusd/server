package common

import (
	"net/http"

	"github.com/codegangsta/negroni"
)

type csrMiddleware struct{}

func (m *csrMiddleware) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	res.Header().Set("Content-Security-Policy", "default-src 'none'")
	next(res, req)
}

// CSR adds Content-Security-Policy header to the response.
func CSR() negroni.Handler {
	return &csrMiddleware{}
}
