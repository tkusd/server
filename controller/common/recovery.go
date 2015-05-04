package common

import (
	"net/http"
	"runtime"

	"github.com/codegangsta/negroni"
	"github.com/tkusd/server/util"
)

const (
	errorFormat = "PANIC: %s\n%s"
	stackSize   = 1024 * 8
	stackAll    = false
)

type recovery struct{}

func (rec *recovery) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			HandleAPIError(res, &util.APIError{
				Code:    util.ServerError,
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})

			// Print error stack
			stack := make([]byte, stackSize)
			stack = stack[:runtime.Stack(stack, stackAll)]

			util.Log().Errorf(errorFormat, err, stack)
		}
	}()

	next(res, req)
}

func NewRecovery() negroni.Handler {
	return &recovery{}
}
