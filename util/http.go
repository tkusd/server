package util

import "net/http"

func NoCacheHeader(res http.ResponseWriter) {
	res.Header().Set("Pragma", "no-cache")
	res.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	res.Header().Set("Expires", "0")
}
