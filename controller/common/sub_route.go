package common

import (
	"net/http"
	"strings"

	"github.com/gorilla/context"
)

const originalPathKey = "originalPath"

type SubRoute struct {
	Path    string
	Handler http.Handler
}

func (route *SubRoute) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if path := strings.TrimPrefix(req.URL.Path, route.Path); len(path) < len(req.URL.Path) {
		defer func() {
			req.URL.Path = GetOriginalPath(req)
			context.Delete(req, originalPathKey)
		}()

		context.Set(req, originalPathKey, req.URL.Path)
		req.URL.Path = path
		route.Handler.ServeHTTP(res, req)
	} else {
		next(res, req)
	}
}

func NewSubRoute(path string, handler http.Handler) *SubRoute {
	return &SubRoute{
		Path:    path,
		Handler: handler,
	}
}

// GetOriginalPath returns the original path in a sub-route.
func GetOriginalPath(req *http.Request) string {
	path, ok := context.GetOk(req, originalPathKey)

	if ok {
		if str, ok := path.(string); ok {
			return str
		}
	}

	return req.URL.Path
}
