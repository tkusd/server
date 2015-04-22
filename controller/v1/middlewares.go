package v1

import (
	"net/http"

	"regexp"

	"github.com/gorilla/mux"
	"github.com/tommy351/app-studio-server/model"
	"github.com/tommy351/app-studio-server/util"
)

var rToken = regexp.MustCompile(`Bearer ([A-Za-z0-9\-\._~\+\/]+=*)`)

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

	if token, err := model.GetTokenBase64(key); err != nil {
		return nil, &util.APIError{
			Code:    util.TokenInvalidError,
			Message: "Token is invalid.",
			Status:  http.StatusUnauthorized,
		}
	} else {
		return token, nil
	}
}

func GetUser(res http.ResponseWriter, req *http.Request) (*model.User, error) {
	vars := mux.Vars(req)

	if id, ok := vars["id"]; ok {
		if user, err := model.GetUser(util.ParseUUID(id)); err == nil {
			return user, nil
		}
	}

	return nil, &util.APIError{
		Code:    util.UserNotFoundError,
		Message: "User not found.",
		Status:  http.StatusNotFound,
	}
}

func CheckUserPermission(res http.ResponseWriter, req *http.Request, userID util.UUID) error {
	if token, err := GetToken(res, req); err != nil {
		return err
	} else {
		if token.UserID.Equal(userID) {
			return nil
		} else {
			return &util.APIError{
				Code:    util.UserForbiddenError,
				Message: "You are forbidden to access this user.",
				Status:  http.StatusForbidden,
			}
		}
	}
}
