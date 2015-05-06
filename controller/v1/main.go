package v1

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/tkusd/server/controller/common"
)

// Params
const (
	userIDParam    = "user_id"
	projectIDParam = "project_id"
	elementIDParam = "element_id"
	keyParam       = "key"
)

// URL patterns
const (
	userCollectionURL         = "/users"
	userSingularURL           = "/users/:" + userIDParam
	projectCollectionURL      = userSingularURL + "/projects"
	projectSingularURL        = "/projects/:" + projectIDParam
	elementCollectionURL      = projectSingularURL + "/elements"
	elementSingularURL        = "/elements/:" + elementIDParam
	childElementCollectionURL = elementSingularURL + "/elements"
	tokenCollectionURL        = "/tokens"
	tokenSingularURL          = "/tokens/:" + keyParam
)

// Router returns a http.Handler.
func Router() http.Handler {
	n := negroni.New()
	r := httprouter.New()

	r.POST(userCollectionURL, common.WrapHandlerFunc(UserCreate))
	r.GET(userSingularURL, common.WrapHandlerFunc(UserShow))
	r.PUT(userSingularURL, common.WrapHandlerFunc(UserUpdate))
	r.DELETE(userSingularURL, common.WrapHandlerFunc(UserDestroy))

	r.GET(projectCollectionURL, common.ChainHandler(CheckUserExist, ProjectList))
	r.POST(projectCollectionURL, common.ChainHandler(CheckUserExist, ProjectCreate))
	r.GET(projectSingularURL, common.WrapHandlerFunc(ProjectShow))
	r.PUT(projectSingularURL, common.WrapHandlerFunc(ProjectUpdate))
	r.DELETE(projectSingularURL, common.WrapHandlerFunc(ProjectDestroy))

	r.GET(elementCollectionURL, common.ChainHandler(CheckProjectExist, ElementList))
	r.POST(elementCollectionURL, common.WrapHandlerFunc(ElementCreate))
	r.GET(elementSingularURL, common.WrapHandlerFunc(ElementShow))
	r.PUT(elementSingularURL, common.WrapHandlerFunc(ElementUpdate))
	r.DELETE(elementSingularURL, common.WrapHandlerFunc(ElementDestroy))

	r.GET(childElementCollectionURL, common.WrapHandlerFunc(ChildElementList))
	r.POST(childElementCollectionURL, common.WrapHandlerFunc(ChildElementCreate))

	r.POST(tokenCollectionURL, common.WrapHandlerFunc(TokenCreate))
	r.DELETE(tokenSingularURL, common.WrapHandlerFunc(TokenDestroy))

	r.NotFound = common.NotFound
	r.HandleMethodNotAllowed = false

	n.Use(common.NewRecovery())
	n.UseHandler(r)

	return n
}
