package v1

import (
	"github.com/gin-gonic/gin"
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
	elementFullURL            = elementSingularURL + "/full"
	childElementCollectionURL = elementSingularURL + "/elements"
	tokenCollectionURL        = "/tokens"
	tokenSingularURL          = "/tokens/:" + tokenIDParam
)

// Router returns a http.Handler.
func Router(r *gin.RouterGroup) {
	r.POST(userCollectionURL, common.Wrap(UserCreate))
	r.GET(userSingularURL, common.Wrap(UserShow))
	r.PUT(userSingularURL, common.Wrap(UserUpdate))
	r.DELETE(userSingularURL, common.Wrap(UserDestroy))

	r.GET(projectCollectionURL, CheckUserExist, common.Wrap(ProjectList))
	r.POST(projectCollectionURL, CheckUserExist, common.Wrap(ProjectCreate))
	r.GET(projectSingularURL, CheckProjectExist, common.Wrap(ProjectShow))
	r.PUT(projectSingularURL, common.Wrap(ProjectUpdate))
	r.DELETE(projectSingularURL, common.Wrap(ProjectDestroy))
	r.GET(projectFullURL, common.Wrap(ProjectFull))

	r.GET(elementCollectionURL, CheckProjectExist, common.Wrap(ElementList))
	r.POST(elementCollectionURL, common.Wrap(ElementCreate))
	r.GET(elementSingularURL, common.Wrap(ElementShow))
	r.PUT(elementSingularURL, common.Wrap(ElementUpdate))
	r.DELETE(elementSingularURL, common.Wrap(ElementDestroy))
	r.GET(elementFullURL, common.Wrap(ElementFull))

	r.GET(childElementCollectionURL, common.Wrap(ChildElementList))
	r.POST(childElementCollectionURL, common.Wrap(ChildElementCreate))

	r.POST(tokenCollectionURL, common.Wrap(TokenCreate))
	r.PUT(tokenSingularURL, common.Wrap(TokenUpdate))
	r.DELETE(tokenSingularURL, common.Wrap(TokenDestroy))
}
