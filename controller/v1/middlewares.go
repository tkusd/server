package v1

import (
	"net/http"

	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/tkusd/server/controller/common"
	"github.com/tkusd/server/model"
	"github.com/tkusd/server/model/types"
	"github.com/tkusd/server/util"
)

var rToken = regexp.MustCompile(`Bearer ([A-Za-z0-9\-\._~\+\/]+=*)`)

func GetIDParam(c *gin.Context, param string) (*types.UUID, error) {
	id := c.Param(param)

	if id == "" {
		return nil, &util.APIError{
			Code:    util.RequiredError,
			Message: "ID is required.",
		}
	}

	uid := types.ParseUUID(id)

	if !uid.Valid() {
		return nil, &util.APIError{
			Code:    util.UUIDError,
			Message: "UUID is invalid.",
		}
	}

	return &uid, nil
}

func GetToken(c *gin.Context) (*model.Token, error) {
	id, err := GetIDParam(c, tokenIDParam)

	if err != nil {
		return nil, err
	}

	if token, err := model.GetToken(*id); err == nil {
		return token, nil
	}

	return nil, &util.APIError{
		Code:    util.TokenNotFoundError,
		Message: "Token not found.",
		Status:  http.StatusNotFound,
	}
}

// CheckToken checks the Authorization header and gets the token from the database.
func CheckToken(c *gin.Context) (*model.Token, error) {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" || !rToken.MatchString(authHeader) {
		return nil, &util.APIError{
			Code:    util.TokenRequiredError,
			Message: "Token is required.",
			Status:  http.StatusUnauthorized,
		}
	}

	//key := rToken.FindStringSubmatch(authHeader)[1]
	id := rToken.FindStringSubmatch(authHeader)[1]
	token, err := model.GetToken(types.ParseUUID(id))

	if err != nil {
		return nil, &util.APIError{
			Code:    util.TokenInvalidError,
			Message: "Token is invalid.",
			Status:  http.StatusUnauthorized,
		}
	}

	return token, nil
}

// GetUser parses user_id in the URL and gets the user data from the database.
func GetUser(c *gin.Context) (*model.User, error) {
	id, err := GetIDParam(c, userIDParam)

	if err != nil {
		return nil, err
	}

	if user, err := model.GetUser(*id); err == nil {
		return user, nil
	}

	return nil, &util.APIError{
		Code:    util.UserNotFoundError,
		Message: "User not found.",
		Status:  http.StatusNotFound,
	}
}

// CheckUserPermission checks whether the current token matching user ID.
func CheckUserPermission(c *gin.Context, userID types.UUID) error {
	token, err := CheckToken(c)

	if err != nil {
		return err
	}

	if token.UserID.Equal(userID) {
		return nil
	}

	return &util.APIError{
		Code:    util.UserForbiddenError,
		Message: "You are forbidden to access.",
		Status:  http.StatusForbidden,
	}
}

// CheckUserExist checks whether the user exists or not.
func CheckUserExist(c *gin.Context) {
	id, err := GetIDParam(c, userIDParam)

	if err != nil {
		common.HandleAPIError(c, err)
		return
	}

	user := &model.User{ID: *id}

	if user.Exists() {
		c.Next()
		return
	}

	common.HandleAPIError(c, &util.APIError{
		Code:    util.UserNotFoundError,
		Message: "User not found.",
		Status:  http.StatusNotFound,
	})
}

// GetProject parses project_id in the URL and gets the project data from the database.
func GetProject(c *gin.Context) (*model.Project, error) {
	id, err := GetIDParam(c, projectIDParam)

	if err != nil {
		return nil, err
	}

	if project, err := model.GetProject(*id); err == nil {
		return project, nil
	}

	return nil, &util.APIError{
		Code:    util.ProjectNotFoundError,
		Message: "Project not found.",
		Status:  http.StatusNotFound,
	}
}

// CheckProjectExist checks whether the project exists or not.
func CheckProjectExist(c *gin.Context) {
	id, err := GetIDParam(c, projectIDParam)

	if err != nil {
		common.HandleAPIError(c, err)
		return
	}

	project := &model.Project{ID: *id}

	if project.Exists() {
		c.Next()
		return
	}

	common.HandleAPIError(c, &util.APIError{
		Code:    util.ProjectNotFoundError,
		Message: "Project not found.",
		Status:  http.StatusNotFound,
	})
}

// GetElement parses element_id in the URL and gets the element data from the database.
func GetElement(c *gin.Context) (*model.Element, error) {
	id, err := GetIDParam(c, elementIDParam)

	if err != nil {
		return nil, err
	}

	if element, err := model.GetElement(*id); err == nil {
		return element, nil
	}

	return nil, &util.APIError{
		Code:    util.ElementNotFoundError,
		Message: "Element not found.",
		Status:  http.StatusNotFound,
	}
}

// CheckProjectPermission checks whether the current user is able to edit the project.
func CheckProjectPermission(c *gin.Context, projectID types.UUID, strict bool) error {
	token, err := CheckToken(c)

	if strict && err != nil {
		return err
	}

	project, err := model.GetProject(projectID)

	if err != nil {
		return nil
	}

	if (token != nil && project.UserID.Equal(token.UserID)) || (!strict && !project.IsPrivate) {
		return nil
	}

	return &util.APIError{
		Code:    util.UserForbiddenError,
		Message: "You are forbidden to access.",
		Status:  http.StatusForbidden,
	}
}

func CheckElementExist(c *gin.Context) {
	id, err := GetIDParam(c, elementIDParam)

	if err != nil {
		common.HandleAPIError(c, err)
		return
	}

	element := &model.Element{ID: *id}

	if element.Exists() {
		c.Next()
		return
	}

	common.HandleAPIError(c, &util.APIError{
		Code:    util.ElementNotFoundError,
		Message: "Element not found.",
		Status:  http.StatusNotFound,
	})
}
