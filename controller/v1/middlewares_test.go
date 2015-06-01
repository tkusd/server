package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/base64"

	"code.google.com/p/go-uuid/uuid"
	"github.com/julienschmidt/httprouter"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tkusd/server/controller/common"
	"github.com/tkusd/server/model"
	"github.com/tkusd/server/model/types"
	"github.com/tkusd/server/util"
)

func TestCheckToken(t *testing.T) {
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

		t, _ := CheckToken(res, req)
		So(t.ID, ShouldResemble, token.ID)
		So(t.UserID, ShouldResemble, token.UserID)
	})

	Convey("Token not found", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		fakeKey := base64.StdEncoding.EncodeToString(uuid.NewRandom())
		req.Header.Set("Authorization", "Bearer "+fakeKey)

		_, err := CheckToken(res, req)
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

		_, err := CheckToken(res, req)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.TokenRequiredError,
			Message: "Token is required.",
			Status:  http.StatusUnauthorized,
		})
	})

	Convey("Authorization header does not exist", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)

		_, err := CheckToken(res, req)
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

	router.GET(userSingularURL, common.WrapHandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if user, err := GetUser(res, req); err != nil {
			common.HandleAPIError(res, req, err)
		} else {
			common.APIResponse(res, req, http.StatusOK, user)
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
			Message: "You are forbidden to access.",
			Status:  http.StatusForbidden,
		})
	})
}

func TestCheckUserExist(t *testing.T) {
	user := new(model.User)
	createTestUser(user, fixtureUsers[0])
	defer user.Delete()

	router := httprouter.New()
	router.GET(userSingularURL, common.ChainHandler(CheckUserExist, func(res http.ResponseWriter, req *http.Request) {
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
	router.GET(projectSingularURL, common.WrapHandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if project, err := GetProject(res, req); err != nil {
			common.HandleAPIError(res, req, err)
		} else {
			common.APIResponse(res, req, http.StatusOK, project)
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

func TestCheckProjectExist(t *testing.T) {
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
	router.GET(projectSingularURL, common.ChainHandler(CheckProjectExist, func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNoContent)
	}))

	Convey("Success", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/projects/"+project.ID.String(), nil)

		router.ServeHTTP(res, req)

		So(res.Code, ShouldEqual, http.StatusNoContent)
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

func TestGetElement(t *testing.T) {
	user := new(model.User)
	createTestUser(user, fixtureUsers[0])
	defer user.Delete()

	token := new(model.Token)
	createTestToken(token, fixtureUsers[0])
	defer token.Delete()

	project := new(model.Project)
	createTestProject(user, token, project, fixtureProjects[0])
	defer project.Delete()

	element := new(model.Element)
	createTestElement(project, token, element, fixtureElements[0])
	defer element.Delete()

	router := httprouter.New()
	router.GET(elementSingularURL, common.WrapHandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if element, err := GetElement(res, req); err != nil {
			common.HandleAPIError(res, req, err)
		} else {
			common.APIResponse(res, req, http.StatusOK, element)
		}
	}))

	Convey("Success", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/elements/"+element.ID.String(), nil)

		router.ServeHTTP(res, req)

		e := new(model.Element)
		So(res.Code, ShouldEqual, http.StatusOK)
		parseJSON(res.Body, e)
		So(e.ID, ShouldResemble, element.ID)
	})

	Convey("Failed", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/elements/"+uuid.New(), nil)

		router.ServeHTTP(res, req)

		err := new(util.APIError)
		So(res.Code, ShouldEqual, http.StatusNotFound)
		parseJSON(res.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.ElementNotFoundError,
			Message: "Element not found.",
		})
	})
}

func TestCheckProjectPermission(t *testing.T) {
	u1 := new(model.User)
	createTestUser(u1, fixtureUsers[0])
	defer u1.Delete()

	u2 := new(model.User)
	createTestUser(u2, fixtureUsers[1])
	defer u2.Delete()

	t1 := new(model.Token)
	createTestToken(t1, fixtureUsers[0])
	defer t1.Delete()

	t2 := new(model.Token)
	createTestToken(t2, fixtureUsers[1])
	defer t2.Delete()

	p1 := new(model.Project)
	createTestProject(u1, t1, p1, fixtureProjects[0])
	defer p1.Delete()

	p2 := new(model.Project)
	createTestProject(u2, t2, p2, fixtureProjects[1])
	defer p2.Delete()

	Convey("Owner", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+t1.ID.String())

		err := CheckProjectPermission(res, req, p1.ID, false)
		So(err, ShouldBeNil)
	})

	Convey("Others + Public", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+t2.ID.String())

		err := CheckProjectPermission(res, req, p1.ID, false)
		So(err, ShouldBeNil)
	})

	Convey("Others + Public + Strict", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+t2.ID.String())

		err := CheckProjectPermission(res, req, p1.ID, true)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.UserForbiddenError,
			Message: "You are forbidden to access.",
			Status:  http.StatusForbidden,
		})
	})

	Convey("Others + Private", t, func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+t1.ID.String())

		err := CheckProjectPermission(res, req, p2.ID, false)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.UserForbiddenError,
			Message: "You are forbidden to access.",
			Status:  http.StatusForbidden,
		})
	})
}
