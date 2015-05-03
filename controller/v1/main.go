package v1

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/tommy351/app-studio-server/controller/common"
)

// Router returns a http.Handler.
func Router() http.Handler {
	n := negroni.New()
	r := httprouter.New()

	r.POST("/users", common.WrapHandlerFunc(UserCreate))
	r.GET("/users/:user_id", common.WrapHandlerFunc(UserShow))
	r.PUT("/users/:user_id", common.WrapHandlerFunc(UserUpdate))
	r.DELETE("/users/:user_id", common.WrapHandlerFunc(UserDestroy))

	r.GET("/users/:user_id/projects", common.ChainHandler(CheckUserExist, ProjectList))
	r.POST("/users/:user_id/projects", common.ChainHandler(CheckUserExist, ProjectCreate))
	r.GET("/projects/:project_id", common.WrapHandlerFunc(ProjectShow))
	r.PUT("/projects/:project_id", common.WrapHandlerFunc(ProjectUpdate))
	r.DELETE("/projects/:project_id", common.WrapHandlerFunc(ProjectDestroy))

	r.POST("/tokens", common.WrapHandlerFunc(TokenCreate))
	r.DELETE("/tokens/:key", common.WrapHandlerFunc(TokenDestroy))

	r.NotFound = common.NotFound

	n.Use(common.NewRecovery())
	n.UseHandler(r)

	return n
}
