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
	tokenIDParam   = "token_id"
)

// URL patterns
const (
	userCollectionURL         = "/users"
	userSingularURL           = "/users/:" + userIDParam
	projectCollectionURL      = userSingularURL + "/projects"
	projectSingularURL        = "/projects/:" + projectIDParam
	projectFullURL            = projectSingularURL + "/full"
	elementCollectionURL      = projectSingularURL + "/elements"
	elementSingularURL        = "/elements/:" + elementIDParam
	childElementCollectionURL = elementSingularURL + "/elements"
	tokenCollectionURL        = "/tokens"
	tokenSingularURL          = "/tokens/:" + tokenIDParam
)

// Router returns a http.Handler.
func Router() http.Handler {
	n := negroni.New()
	r := httprouter.New()

	r.POST(userCollectionURL, common.WrapCommonHandle(UserCreate))
	r.GET(userSingularURL, common.WrapCommonHandle(UserShow))
	r.PUT(userSingularURL, common.WrapCommonHandle(UserUpdate))
	r.DELETE(userSingularURL, common.WrapCommonHandle(UserDestroy))

	r.GET(projectCollectionURL, common.ChainHandler(CheckUserExist, ProjectList))
	r.POST(projectCollectionURL, common.ChainHandler(CheckUserExist, ProjectCreate))
	r.GET(projectSingularURL, common.ChainHandler(CheckProjectExist, ProjectShow))
	r.PUT(projectSingularURL, common.WrapCommonHandle(ProjectUpdate))
	r.DELETE(projectSingularURL, common.WrapCommonHandle(ProjectDestroy))
	r.GET(projectFullURL, common.ChainHandler(CheckProjectExist, ProjectFull))

	r.GET(elementCollectionURL, common.ChainHandler(CheckProjectExist, ElementList))
	r.POST(elementCollectionURL, common.WrapCommonHandle(ElementCreate))
	r.GET(elementSingularURL, common.WrapCommonHandle(ElementShow))
	r.PUT(elementSingularURL, common.WrapCommonHandle(ElementUpdate))
	r.DELETE(elementSingularURL, common.WrapCommonHandle(ElementDestroy))

	r.GET(childElementCollectionURL, common.WrapCommonHandle(ChildElementList))
	r.POST(childElementCollectionURL, common.WrapCommonHandle(ChildElementCreate))

	r.POST(tokenCollectionURL, common.WrapCommonHandle(TokenCreate))
	r.PUT(tokenSingularURL, common.WrapCommonHandle(TokenUpdate))
	r.DELETE(tokenSingularURL, common.WrapCommonHandle(TokenDestroy))

	r.NotFound = common.NotFound
	r.HandleMethodNotAllowed = false

	n.Use(common.NewRecovery())
	n.UseHandler(r)

	return n
}
