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

func createTestToken(data interface{}, body interface{}) *httptest.ResponseRecorder {
	r := request(&requestOptions{
		Method: "POST",
		URL:    "/tokens",
		Body:   body,
	})

	if err := parseJSON(r.Body, data); err != nil {
		log.Fatal(err)
	}

	return r
}

func TestTokenCreate(t *testing.T) {
	user := new(model.User)
	createTestUser(user, fixtureUsers[0])
	defer user.Delete()

	Convey("Success", t, func() {
		now := time.Now().Truncate(time.Second)
		token := new(model.Token)
		r := createTestToken(token, fixtureUsers[0])
		defer token.Delete()

		So(r.Code, ShouldEqual, http.StatusCreated)
		So(r.Header().Get("Pragma"), ShouldEqual, "no-cache")
		So(r.Header().Get("Cache-Control"), ShouldEqual, "no-cache, no-store, must-revalidate")
		So(r.Header().Get("Expires"), ShouldEqual, "0")
		So(token.ID, ShouldNotBeEmpty)
		So(token.UserID, ShouldResemble, user.ID)
		So(token.CreatedAt.Time, ShouldHappenOnOrAfter, now)
	})

	Convey("Email is required", t, func() {
		err := new(util.APIError)
		r := createTestToken(err, map[string]interface{}{})

		So(r.Code, ShouldEqual, http.StatusBadRequest)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.RequiredError,
			Message: "Email is required.",
			Field:   "email",
		})
	})

	Convey("Email is invalid", t, func() {
		err := new(util.APIError)
		r := createTestToken(err, map[string]interface{}{
			"email": "abc@",
		})

		So(r.Code, ShouldEqual, http.StatusBadRequest)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.EmailError,
			Message: "Email is invalid.",
			Field:   "email",
		})
	})

	Convey("User does not exist", t, func() {
		err := new(util.APIError)
		r := createTestToken(err, map[string]interface{}{
			"email": "nothing@nothing.com",
		})

		So(r.Code, ShouldEqual, http.StatusBadRequest)
		So(err, ShouldResemble, &util.APIError{
			Field:   "email",
			Code:    util.UserNotFoundError,
			Message: "User does not exist.",
		})
	})

	Convey("Password is wrong", t, func() {
		err := new(util.APIError)
		r := createTestToken(err, map[string]interface{}{
			"email":    fixtureUsers[0].Email,
			"password": "werjwoerijwer",
		})

		So(r.Code, ShouldEqual, http.StatusBadRequest)
		So(err, ShouldResemble, &util.APIError{
			Field:   "password",
			Code:    util.WrongPasswordError,
			Message: "Password is wrong.",
		})
	})
}

func TestTokenDestroy(t *testing.T) {
	user := new(model.User)
	createTestUser(user, fixtureUsers[0])
	defer user.Delete()

	Convey("Success", t, func() {
		token := new(model.Token)
		createTestToken(token, fixtureUsers[0])
		r := request(&requestOptions{
			Method: "DELETE",
			URL:    "/tokens/" + token.ID.String(),
		})

		So(r.Code, ShouldEqual, http.StatusNoContent)

		_, err := model.GetToken(token.ID)
		So(err, ShouldNotBeNil)
	})

	Convey("Token not found", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "DELETE",
			URL:    "/tokens/nothing",
		})

		So(r.Code, ShouldEqual, http.StatusNotFound)

		if e := parseJSON(r.Body, err); e != nil {
			log.Fatal(e)
		}

		So(err, ShouldResemble, &util.APIError{
			Code:    util.TokenNotFoundError,
			Message: "Token does not exist.",
		})
	})
}
