package v1

import (
	"net/http/httptest"

	"testing"

	"log"
	"net/http"
	"time"

	"code.google.com/p/go-uuid/uuid"
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

func TestUserShow(t *testing.T) {
	user := new(model.User)
	createTestUser(user, fixtureUsers[0])
	defer model.DeleteUser(user)

	token := new(model.Token)
	createTestToken(token, fixtureUsers[0])
	defer model.DeleteToken(token)

	Convey("Success (private)", t, func() {
		u := new(model.User)
		r := request(&requestOptions{
			Method: "GET",
			URL:    "/users/" + user.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + token.GetBase64ID(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusOK)
		parseJSON(r.Body, u)
		So(u.ID, ShouldResemble, user.ID)
		So(u.Name, ShouldEqual, user.Name)
		So(u.Email, ShouldEqual, user.Email)
		So(u.Avatar, ShouldEqual, user.Avatar)
		So(u.CreatedAt, ShouldResemble, user.CreatedAt.Truncate(time.Second))
		So(u.UpdatedAt, ShouldResemble, user.UpdatedAt.Truncate(time.Second))
		So(u.IsActivated, ShouldEqual, user.IsActivated)
		So(u.Password, ShouldBeEmpty)
		So(u.ActivationToken, ShouldBeEmpty)
	})

	Convey("Success (public)", t, func() {
		u := new(model.User)
		r := request(&requestOptions{
			Method: "GET",
			URL:    "/users/" + user.ID.String(),
		})

		So(r.Code, ShouldEqual, http.StatusOK)
		parseJSON(r.Body, u)
		So(u.Email, ShouldBeEmpty)
	})

	Convey("User not found", t, func() {
		r := request(&requestOptions{
			Method: "GET",
			URL:    "/users/" + uuid.New(),
			Headers: map[string]string{
				"Authorization": "Bearer " + token.GetBase64ID(),
			},
		})

		err := new(util.APIError)
		So(r.Code, ShouldEqual, http.StatusNotFound)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.UserNotFoundError,
			Message: "User not found.",
		})
	})
}

func TestUserUpdate(t *testing.T) {
	user := new(model.User)
	createTestUser(user, fixtureUsers[0])
	defer model.DeleteUser(user)

	user2 := new(model.User)
	createTestUser(user2, fixtureUsers[1])
	defer model.DeleteUser(user)

	token := new(model.Token)
	createTestToken(token, fixtureUsers[0])
	defer model.DeleteToken(token)

	token2 := new(model.Token)
	createTestToken(token2, fixtureUsers[1])
	defer model.DeleteToken(token2)

	Convey("Success", t, func() {
		now := time.Now().Truncate(time.Second)
		u := new(model.User)
		r := request(&requestOptions{
			Method: "PUT",
			URL:    "/users/" + user.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + token.GetBase64ID(),
			},
			Body: map[string]interface{}{},
		})

		So(r.Code, ShouldEqual, http.StatusOK)
		parseJSON(r.Body, u)
		So(u.ID, ShouldResemble, user.ID)
		So(u.Name, ShouldEqual, user.Name)
		So(u.Email, ShouldEqual, user.Email)
		So(u.Avatar, ShouldEqual, user.Avatar)
		So(u.CreatedAt, ShouldResemble, user.CreatedAt.Truncate(time.Second))
		So(u.UpdatedAt, ShouldHappenOnOrAfter, now)
		So(u.IsActivated, ShouldEqual, user.IsActivated)
		So(u.Password, ShouldBeEmpty)
		So(u.ActivationToken, ShouldBeEmpty)

		user, _ = model.GetUser(user.ID)
	})

	Convey("Update name", t, func() {
		newName := "johnnnn"
		u := new(model.User)
		r := request(&requestOptions{
			Method: "PUT",
			URL:    "/users/" + user.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + token.GetBase64ID(),
			},
			Body: map[string]interface{}{
				"name": newName,
			},
		})

		So(r.Code, ShouldEqual, http.StatusOK)
		parseJSON(r.Body, u)
		So(u.Name, ShouldEqual, newName)

		user, _ = model.GetUser(user.ID)
		So(user.Name, ShouldEqual, newName)
	})

	Convey("Update email", t, func() {
		newEmail := "jgdfjgdfg@jgeorj.com"
		u := new(model.User)
		r := request(&requestOptions{
			Method: "PUT",
			URL:    "/users/" + user.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + token.GetBase64ID(),
			},
			Body: map[string]interface{}{
				"email": newEmail,
			},
		})

		So(r.Code, ShouldEqual, http.StatusOK)
		parseJSON(r.Body, u)
		So(u.Email, ShouldEqual, newEmail)

		user, _ = model.GetUser(user.ID)
		So(user.Email, ShouldEqual, newEmail)
	})

	Convey("Update password", t, func() {
		newPassword := "fejfosdijfsd"
		u := new(model.User)
		r := request(&requestOptions{
			Method: "PUT",
			URL:    "/users/" + user.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + token.GetBase64ID(),
			},
			Body: map[string]interface{}{
				"password":     newPassword,
				"old_password": fixtureUsers[0].Password,
			},
		})

		So(r.Code, ShouldEqual, http.StatusOK)
		parseJSON(r.Body, u)
		So(u.Password, ShouldBeEmpty)

		user, _ = model.GetUser(user.ID)
		err := user.Authenticate(newPassword)
		So(err, ShouldBeNil)
	})

	Convey("Current password is required", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "PUT",
			URL:    "/users/" + user.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + token.GetBase64ID(),
			},
			Body: map[string]interface{}{
				"password": fixtureUsers[0].Password,
			},
		})

		So(r.Code, ShouldEqual, http.StatusBadRequest)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Field:   "old_password",
			Code:    util.RequiredError,
			Message: "Current password is required.",
		})
	})

	Convey("Password is wrong", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "PUT",
			URL:    "/users/" + user.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + token.GetBase64ID(),
			},
			Body: map[string]interface{}{
				"password":     fixtureUsers[0].Password,
				"old_password": "weroijweorijweorjweorjwe",
			},
		})

		So(r.Code, ShouldEqual, http.StatusBadRequest)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Field:   "old_password",
			Code:    util.WrongPasswordError,
			Message: "Password is wrong.",
		})
	})

	Convey("Forbidden", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "PUT",
			URL:    "/users/" + user.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + token2.GetBase64ID(),
			},
			Body: map[string]interface{}{},
		})

		So(r.Code, ShouldEqual, http.StatusForbidden)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.UserForbiddenError,
			Message: "You are forbidden to access this user.",
		})
	})

	Convey("User not found", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "PUT",
			URL:    "/users/" + uuid.New(),
			Headers: map[string]string{
				"Authorization": "Bearer " + token.GetBase64ID(),
			},
			Body: map[string]interface{}{},
		})

		So(r.Code, ShouldEqual, http.StatusNotFound)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.UserNotFoundError,
			Message: "User not found.",
		})
	})
}

func TestUserDestroy(t *testing.T) {
	user := new(model.User)
	createTestUser(user, fixtureUsers[0])
	defer model.DeleteUser(user)

	user2 := new(model.User)
	createTestUser(user2, fixtureUsers[1])
	defer model.DeleteUser(user2)

	token := new(model.Token)
	createTestToken(token, fixtureUsers[0])
	defer model.DeleteToken(token)

	token2 := new(model.Token)
	createTestToken(token2, fixtureUsers[1])
	defer model.DeleteToken(token2)

	Convey("Success", t, func() {
		r := request(&requestOptions{
			Method: "DELETE",
			URL:    "/users/" + user.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + token.GetBase64ID(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusNoContent)

		_, err := model.GetUser(user.ID)
		So(err, ShouldNotBeNil)

		createTestUser(user, fixtureUsers[0])
		createTestToken(token, fixtureUsers[0])
	})

	Convey("Forbidden", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "DELETE",
			URL:    "/users/" + user.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + token2.GetBase64ID(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusForbidden)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.UserForbiddenError,
			Message: "You are forbidden to access this user.",
		})
	})

	Convey("User not found", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "DELETE",
			URL:    "/users/" + uuid.New(),
			Headers: map[string]string{
				"Authorization": "Bearer " + token.GetBase64ID(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusNotFound)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.UserNotFoundError,
			Message: "User not found.",
		})
	})
}
