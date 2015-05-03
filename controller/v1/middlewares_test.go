package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/base64"

	"code.google.com/p/go-uuid/uuid"
	"github.com/julienschmidt/httprouter"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tommy351/app-studio-server/controller/common"
	"github.com/tommy351/app-studio-server/model"
	"github.com/tommy351/app-studio-server/model/types"
	"github.com/tommy351/app-studio-server/util"
)

func TestGetToken(t *testing.T) {
	Convey("Success", t, func() {
		user := new(model.User)
		createTestUser(user, fixtureUsers[0])
		defer user.Delete()

		token := new(model.Token)
		createTestToken(token, fixtureUsers[0])
		defer token.Delete()

		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+token.ID.String())

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
	defer user.Delete()

	router := httprouter.New()

	router.GET("/users/:user_id", common.WrapHandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if user, err := GetUser(res, req); err != nil {
			common.HandleAPIError(res, err)
		} else {
			common.RenderJSON(res, http.StatusOK, user)
		}
	}))

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
	defer user.Delete()

	token := new(model.Token)
	createTestToken(token, fixtureUsers[0])
	defer token.Delete()

	Convey("Success", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+token.ID.String())

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
		req.Header.Set("Authorization", "Bearer "+token.ID.String())

		err := CheckUserPermission(res, req, types.NewRandomUUID())
		So(err, ShouldResemble, &util.APIError{
			Code:    util.UserForbiddenError,
			Message: "You are forbidden to access this user.",
			Status:  http.StatusForbidden,
		})
	})
}

func TestCheckUserExist(t *testing.T) {
	user := new(model.User)
	createTestUser(user, fixtureUsers[0])
	defer user.Delete()

	router := httprouter.New()
	router.GET("/users/:user_id", common.ChainHandler(CheckUserExist, func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNoContent)
	}))

	Convey("Success", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/"+user.ID.String(), nil)

		router.ServeHTTP(res, req)

		So(res.Code, ShouldEqual, http.StatusNoContent)
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

func TestGetProject(t *testing.T) {
	user := new(model.User)
	createTestUser(user, fixtureUsers[0])
	defer user.Delete()

	token := new(model.Token)
	createTestToken(token, fixtureUsers[0])
	defer token.Delete()

	project := new(model.Project)
	createTestProject(user, token, project, fixtureProjects[0])
	defer project.Delete()

	router := httprouter.New()
	router.GET("/projects/:project_id", common.WrapHandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if project, err := GetProject(res, req); err != nil {
			common.HandleAPIError(res, err)
		} else {
			common.RenderJSON(res, http.StatusOK, project)
		}
	}))

	Convey("Success", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/projects/"+project.ID.String(), nil)

		router.ServeHTTP(res, req)

		p := new(model.Project)
		So(res.Code, ShouldEqual, http.StatusOK)
		parseJSON(res.Body, p)
		So(p.ID, ShouldResemble, project.ID)
	})

	Convey("Failed", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/projects/"+uuid.New(), nil)

		router.ServeHTTP(res, req)

		err := new(util.APIError)
		parseJSON(res.Body, err)
		So(res.Code, ShouldEqual, http.StatusNotFound)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.ProjectNotFoundError,
			Message: "Project not found.",
		})
	})
}
