package v1

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"code.google.com/p/go-uuid/uuid"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/tommy351/app-studio-server/model"
	"github.com/tommy351/app-studio-server/model/types"
	"github.com/tommy351/app-studio-server/util"
)

func createTestElement(project *model.Project, token *model.Token, data interface{}, body interface{}) *httptest.ResponseRecorder {
	r := request(&requestOptions{
		Method: "POST",
		URL:    "/projects/" + project.ID.String() + "/elements",
		Body:   body,
		Headers: map[string]string{
			"Authorization": "Bearer " + token.ID.String(),
		},
	})

	if err := parseJSON(r.Body, data); err != nil {
		log.Fatal(err)
	}

	return r
}

func TestElementList(t *testing.T) {
	//
}

func TestElementCreate(t *testing.T) {
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

	Convey("Success", t, func() {
		element := new(model.Element)
		now := time.Now().Truncate(time.Second)
		r := createTestElement(p1, t1, element, fixtureElements[0])
		defer element.Delete()

		So(r.Code, ShouldEqual, http.StatusCreated)
		So(element.Name, ShouldEqual, fixtureElements[0].Name)
		So(element.Type, ShouldEqual, fixtureElements[0].Type)
		So(element.CreatedAt.Time, ShouldHappenOnOrAfter, now)
		So(element.UpdatedAt.Time, ShouldHappenOnOrAfter, now)
		So(element.ProjectID, ShouldResemble, p1.ID)
		So(element.Attributes, ShouldResemble, fixtureElements[0].Attributes)
	})

	Convey("Forbidden", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "POST",
			URL:    "/projects/" + p1.ID.String() + "/elements",
			Headers: map[string]string{
				"Authorization": "Bearer " + t2.ID.String(),
			},
			Body: map[string]interface{}{},
		})

		So(r.Code, ShouldEqual, http.StatusForbidden)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.UserForbiddenError,
			Message: "You are forbidden to access.",
		})
	})

	Convey("Project not found", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "POST",
			URL:    "/projects/" + uuid.New() + "/elements",
		})

		So(r.Code, ShouldEqual, http.StatusNotFound)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.ProjectNotFoundError,
			Message: "Project not found.",
		})
	})
}

func TestElementShow(t *testing.T) {
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

	e1 := new(model.Element)
	createTestElement(p1, t1, e1, fixtureElements[0])
	defer e1.Delete()

	e2 := new(model.Element)
	createTestElement(p2, t2, e2, fixtureElements[0])
	defer e2.Delete()

	Convey("Success", t, func() {
		element := new(model.Element)
		r := request(&requestOptions{
			Method: "GET",
			URL:    "/elements/" + e1.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + t1.ID.String(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusOK)
		parseJSON(r.Body, element)
		So(element.ID, ShouldResemble, e1.ID)
		So(element.Name, ShouldEqual, e1.Name)
		So(element.ProjectID, ShouldResemble, e1.ProjectID)
		So(element.Type, ShouldEqual, e1.Type)
		So(element.CreatedAt.Time, ShouldResemble, e1.CreatedAt.Truncate(time.Second))
		So(element.UpdatedAt.Time, ShouldResemble, e1.UpdatedAt.Truncate(time.Second))
		So(element.Attributes, ShouldResemble, e1.Attributes)
	})

	Convey("Public", t, func() {
		element := new(model.Element)
		r := request(&requestOptions{
			Method: "GET",
			URL:    "/elements/" + e1.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + t2.ID.String(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusOK)
		parseJSON(r.Body, element)
		So(element.ID, ShouldResemble, e1.ID)
	})

	Convey("Private", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "GET",
			URL:    "/elements/" + e2.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + t1.ID.String(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusForbidden)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.UserForbiddenError,
			Message: "You are forbidden to access.",
		})
	})

	Convey("Element not found", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "GET",
			URL:    "/elements/" + uuid.New(),
		})

		So(r.Code, ShouldEqual, http.StatusNotFound)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.ElementNotFoundError,
			Message: "Element not found.",
		})
	})
}

func TestElementUpdate(t *testing.T) {
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

	e1 := new(model.Element)
	createTestElement(p1, t1, e1, fixtureElements[0])
	defer e1.Delete()

	Convey("Success", t, func() {
		element := new(model.Element)
		now := time.Now().Truncate(time.Second)
		newName := "New name"
		newType := types.ElementTypeText
		newAttrs := types.JSONObject(map[string]interface{}{
			"bar": "123",
		})

		r := request(&requestOptions{
			Method: "PUT",
			URL:    "/elements/" + e1.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + t1.ID.String(),
			},
			Body: map[string]interface{}{
				"name":       newName,
				"type":       newType,
				"attributes": newAttrs,
			},
		})

		So(r.Code, ShouldEqual, http.StatusOK)
		parseJSON(r.Body, element)
		So(element.ID, ShouldResemble, e1.ID)
		So(element.Name, ShouldEqual, newName)
		So(element.Type, ShouldEqual, newType)
		So(element.Attributes, ShouldResemble, newAttrs)
		So(element.CreatedAt.Time, ShouldResemble, e1.CreatedAt.Truncate(time.Second))
		So(element.UpdatedAt.Time, ShouldHappenOnOrAfter, now)
	})

	Convey("Forbidden", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "PUT",
			URL:    "/elements/" + e1.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + t2.ID.String(),
			},
			Body: map[string]interface{}{},
		})

		So(r.Code, ShouldEqual, http.StatusForbidden)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.UserForbiddenError,
			Message: "You are forbidden to access.",
		})
	})

	Convey("Element not found", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "PUT",
			URL:    "/elements/" + uuid.New(),
			Headers: map[string]string{
				"Authorization": "Bearer " + t1.ID.String(),
			},
			Body: map[string]interface{}{},
		})

		So(r.Code, ShouldEqual, http.StatusNotFound)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.ElementNotFoundError,
			Message: "Element not found.",
		})
	})
}

func TestElementDestroy(t *testing.T) {
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

	e1 := new(model.Element)
	createTestElement(p1, t1, e1, fixtureElements[0])
	defer e1.Delete()

	Convey("Success", t, func() {
		r := request(&requestOptions{
			Method: "DELETE",
			URL:    "/elements/" + e1.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + t1.ID.String(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusNoContent)

		_, err := model.GetElement(e1.ID)
		So(err, ShouldNotBeNil)

		createTestElement(p1, t1, e1, fixtureElements[0])
	})

	Convey("Forbidden", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "DELETE",
			URL:    "/elements/" + e1.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + t2.ID.String(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusForbidden)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.UserForbiddenError,
			Message: "You are forbidden to access.",
		})
	})

	Convey("Element not found", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "DELETE",
			URL:    "/elements/" + uuid.New(),
			Headers: map[string]string{
				"Authorization": "Bearer " + t1.ID.String(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusNotFound)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.ElementNotFoundError,
			Message: "Element not found.",
		})
	})
}
