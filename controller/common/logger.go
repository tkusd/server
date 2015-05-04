package common

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/tkusd/server/util"
)

type logger struct{}

func (l *logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	defer func() {
		res := rw.(negroni.ResponseWriter)

		util.Log().WithFields(logrus.Fields{
			"start":    start,
			"duration": time.Since(start),
			"method":   r.Method,
			"code":     res.Status(),
		}).Info(r.URL)
	}()

	next(rw, r)
}

func NewLogger() negroni.Handler {
	return &logger{}
}
