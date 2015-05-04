package v1

import (
	"net/http"

	"regexp"

	"github.com/tkusd/server/controller/common"
	"github.com/tkusd/server/model"
	"github.com/tkusd/server/model/types"
	"github.com/tkusd/server/util"
)

var rToken = regexp.MustCompile(`Bearer ([A-Za-z0-9\-\._~\+\/]+=*)`)

// GetToken checks the Authorization header and gets the token from the database.
func GetToken(res http.ResponseWriter, req *http.Request) (*model.Token, error) {
	authHeader := req.Header.Get("Authorization")

	if authHeader == "" || !rToken.MatchString(authHeader) {
		return nil, &util.APIError{
			Code:    util.TokenRequiredError,
			Message: "Token is required.",
			Status:  http.StatusUnauthorized,
		}
	}

	key := rToken.FindStringSubmatch(authHeader)[1]
	token, err := model.GetTokenBase64(key)

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
func GetUser(res http.ResponseWriter, req *http.Request) (*model.User, error) {
	if id := common.GetParam(req, userIDParam); id != "" {
		if user, err := model.GetUser(types.ParseUUID(id)); err == nil {
			return user, nil
		}
	}

	return nil, &util.APIError{
		Code:    util.UserNotFoundError,
		Message: "User not found.",
		Status:  http.StatusNotFound,
	}
}

// CheckUserPermission checks whether the current token matching user ID.
func CheckUserPermission(res http.ResponseWriter, req *http.Request, userID types.UUID) error {
	token, err := GetToken(res, req)

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
func CheckUserExist(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if id := common.GetParam(req, userIDParam); id != "" {
		user := &model.User{ID: types.ParseUUID(id)}

		if user.Exists() {
			next(res, req)
			return
		}
	}

	common.HandleAPIError(res, &util.APIError{
		Code:    util.UserNotFoundError,
		Message: "User not found.",
		Status:  http.StatusNotFound,
	})
}

// GetProject parses project_id in the URL and gets the project data from the database.
func GetProject(res http.ResponseWriter, req *http.Request) (*model.Project, error) {
	if id := common.GetParam(req, projectIDParam); id != "" {
		if project, err := model.GetProject(types.ParseUUID(id)); err == nil {
			return project, nil
		}
	}

	return nil, &util.APIError{
		Code:    util.ProjectNotFoundError,
		Message: "Project not found.",
		Status:  http.StatusNotFound,
	}
}

// CheckProjectExist checks whether the project exists or not.
func CheckProjectExist(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if id := common.GetParam(req, projectIDParam); id != "" {
		project := &model.Project{ID: types.ParseUUID(id)}

		if project.Exists() {
			next(res, req)
			return
		}
	}

	common.HandleAPIError(res, &util.APIError{
		Code:    util.ProjectNotFoundError,
		Message: "Project not found.",
		Status:  http.StatusNotFound,
	})
}

// GetElement parses element_id in the URL and gets the element data from the database.
func GetElement(res http.ResponseWriter, req *http.Request) (*model.Element, error) {
	if id := common.GetParam(req, elementIDParam); id != "" {
		if element, err := model.GetElement(types.ParseUUID(id)); err == nil {
			return element, nil
		}
	}

	return nil, &util.APIError{
		Code:    util.ElementNotFoundError,
		Message: "Element not found.",
		Status:  http.StatusNotFound,
	}
}

// CheckProjectPermission checks whether the current user is able to edit the project.
func CheckProjectPermission(res http.ResponseWriter, req *http.Request, projectID types.UUID, strict bool) error {
	token, err := GetToken(res, req)

	if err != nil {
		return err
	}

	project, err := model.GetProject(projectID)

	if err == nil {
		if project.UserID.Equal(token.UserID) || (!strict && !project.IsPrivate) {
			return nil
		}
	}

	return &util.APIError{
		Code:    util.UserForbiddenError,
		Message: "You are forbidden to access.",
		Status:  http.StatusForbidden,
	}
}
