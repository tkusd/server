package v1

import (
	"net/http/httptest"

	"testing"

	"log"
	"net/http"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/tommy351/app-studio-server/model"
	"github.com/tommy351/app-studio-server/util"
)

func createTestUser(data interface{}, body interface{}) *httptest.ResponseRecorder {
	r := request(&requestOptions{
		Method: "POST",
		URL:    "/users",
		Body:   body,
	})

	if err := parseJSON(r.Body, data); err != nil {
		log.Fatal(err)
	}

	return r
}

func TestUserCreate(t *testing.T) {
	Convey("Success", t, func() {
		user := new(model.User)
		now := time.Now().Truncate(time.Second)
		r := createTestUser(user, fixtureUsers[0])
		defer model.DeleteUser(user)

		So(r.Code, ShouldEqual, http.StatusCreated)
		So(user.Name, ShouldEqual, fixtureUsers[0].Name)
		So(user.Email, ShouldEqual, fixtureUsers[0].Email)
		So(user.IsActivated, ShouldBeFalse)
		So(user.Avatar, ShouldEqual, util.Gravatar(fixtureUsers[0].Email))
		So(user.Password, ShouldBeEmpty)
		So(user.CreatedAt, ShouldHappenOnOrAfter, now)
		So(user.UpdatedAt, ShouldHappenOnOrAfter, now)

		realUser, err := model.GetUser(user.ID)

		if err != nil {
			log.Fatal(err)
		}

		err = realUser.Authenticate(fixtureUsers[0].Password)

		So(err, ShouldBeNil)
	})

	Convey("Name is required", t, func() {
		err := new(util.APIError)
		r := createTestUser(err, map[string]interface{}{})

		So(r.Code, ShouldEqual, http.StatusBadRequest)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.RequiredError,
			Message: "Name is required.",
			Field:   "name",
		})
	})

	Convey("Email is required", t, func() {
		err := new(util.APIError)
		r := createTestUser(err, map[string]interface{}{
			"name": "abc",
		})

		So(r.Code, ShouldEqual, http.StatusBadRequest)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.RequiredError,
			Message: "Email is required.",
			Field:   "email",
		})
	})

	Convey("Password is required", t, func() {
		err := new(util.APIError)
		r := createTestUser(err, map[string]interface{}{
			"name":  "abc",
			"email": "abc@example.com",
		})

		So(r.Code, ShouldEqual, http.StatusBadRequest)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.RequiredError,
			Message: "Password is required.",
			Field:   "password",
		})
	})
}
