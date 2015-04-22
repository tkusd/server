package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/base64"

	"code.google.com/p/go-uuid/uuid"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tommy351/app-studio-server/model"
	"github.com/tommy351/app-studio-server/util"
)

func TestGetToken(t *testing.T) {
	Convey("Success", t, func() {
		user := new(model.User)
		createTestUser(user, fixtureUsers[0])
		defer model.DeleteUser(user)

		token := new(model.Token)
		createTestToken(token, fixtureUsers[0])
		defer model.DeleteToken(token)

		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+token.GetBase64ID())

		t, _ := GetToken(res, req)
		So(t.ID, ShouldResemble, token.ID)
		So(t.UserID, ShouldResemble, token.UserID)
	})

	Convey("Token not found", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		fakeKey := base64.StdEncoding.EncodeToString(uuid.NewRandom())
		req.Header.Set("Authorization", "Bearer "+fakeKey)

		_, err := GetToken(res, req)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.TokenInvalidError,
			Message: "Token is invalid.",
			Status:  http.StatusUnauthorized,
		})
	})

	Convey("Wrong header format", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", uuid.New())

		_, err := GetToken(res, req)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.TokenRequiredError,
			Message: "Token is required.",
			Status:  http.StatusUnauthorized,
		})
	})

	Convey("Authorization header does not exist", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)

		_, err := GetToken(res, req)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.TokenRequiredError,
			Message: "Token is required.",
			Status:  http.StatusUnauthorized,
		})
	})
}

func TestGetUser(t *testing.T) {
	user := new(model.User)
	createTestUser(user, fixtureUsers[0])
	defer model.DeleteUser(user)

	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", func(res http.ResponseWriter, req *http.Request) {
		if user, err := GetUser(res, req); err != nil {
			util.HandleAPIError(res, err)
		} else {
			util.RenderJSON(res, http.StatusOK, user)
		}
	})

	Convey("Success", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/"+user.ID.String(), nil)

		router.ServeHTTP(res, req)

		u := new(model.User)
		So(res.Code, ShouldEqual, http.StatusOK)
		parseJSON(res.Body, u)
		So(u.ID, ShouldResemble, user.ID)
	})

	Convey("Failed", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/"+uuid.New(), nil)

		router.ServeHTTP(res, req)

		err := new(util.APIError)
		parseJSON(res.Body, err)
		So(res.Code, ShouldEqual, http.StatusNotFound)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.UserNotFoundError,
			Message: "User not found.",
		})
	})
}

func TestCheckUserPermission(t *testing.T) {
	user := new(model.User)
	createTestUser(user, fixtureUsers[0])
	defer model.DeleteUser(user)

	token := new(model.Token)
	createTestToken(token, fixtureUsers[0])
	defer model.DeleteToken(token)

	Convey("Success", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+token.GetBase64ID())

		err := CheckUserPermission(res, req, user.ID)
		So(err, ShouldBeNil)
	})

	Convey("Token not found", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)

		err := CheckUserPermission(res, req, user.ID)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.TokenRequiredError,
			Message: "Token is required.",
			Status:  http.StatusUnauthorized,
		})
	})

	Convey("User ID does not match", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+token.GetBase64ID())

		err := CheckUserPermission(res, req, util.NewRandomUUID())
		So(err, ShouldResemble, &util.APIError{
			Code:    util.UserForbiddenError,
			Message: "You are forbidden to access this user.",
			Status:  http.StatusForbidden,
		})
	})
}
