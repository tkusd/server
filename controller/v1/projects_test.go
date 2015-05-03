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
	"github.com/tommy351/app-studio-server/util"
)

func createTestProject(user *model.User, token *model.Token, data interface{}, body interface{}) *httptest.ResponseRecorder {
	r := request(&requestOptions{
		Method: "POST",
		URL:    "/users/" + user.ID.String() + "/projects",
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

func TestProjectList(t *testing.T) {
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

	Convey("Success", t, func() {
		var list []*model.Project
		r := request(&requestOptions{
			Method: "GET",
			URL:    "/users/" + u1.ID.String() + "/projects",
			Headers: map[string]string{
				"Authorization": "Bearer " + t1.ID.String(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusOK)
		parseJSON(r.Body, &list)
		So(len(list), ShouldEqual, 1)
		So(list[0].ID, ShouldResemble, p1.ID)
	})

	Convey("Show owner's all projects", t, func() {
		var list []*model.Project
		r := request(&requestOptions{
			Method: "GET",
			URL:    "/users/" + u2.ID.String() + "/projects",
			Headers: map[string]string{
				"Authorization": "Bearer " + t2.ID.String(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusOK)
		parseJSON(r.Body, &list)
		So(len(list), ShouldEqual, 1)
		So(list[0].ID, ShouldResemble, p2.ID)
	})

	Convey("Hide private projects to others", t, func() {
		var list []*model.Project
		r := request(&requestOptions{
			Method: "GET",
			URL:    "/users/" + u2.ID.String() + "/projects",
			Headers: map[string]string{
				"Authorization": "Bearer " + t1.ID.String(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusOK)
		parseJSON(r.Body, &list)
		So(len(list), ShouldEqual, 0)
	})

	Convey("User not found", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "GET",
			URL:    "/users/" + uuid.New() + "/projects",
		})

		So(r.Code, ShouldEqual, http.StatusNotFound)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.UserNotFoundError,
			Message: "User not found.",
		})
	})
}

func TestProjectCreate(t *testing.T) {
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

	Convey("Success", t, func() {
		project := new(model.Project)
		now := time.Now().Truncate(time.Second)
		r := createTestProject(u1, t1, project, fixtureProjects[0])
		defer project.Delete()

		So(r.Code, ShouldEqual, http.StatusCreated)
		So(project.Title, ShouldEqual, fixtureProjects[0].Title)
		So(project.Description, ShouldEqual, fixtureProjects[0].Description)
		So(project.IsPrivate, ShouldEqual, fixtureProjects[0].IsPrivate)
		So(project.CreatedAt.Time, ShouldHappenOnOrAfter, now)
		So(project.UpdatedAt.Time, ShouldHappenOnOrAfter, now)
	})

	Convey("Title is required", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "POST",
			URL:    "/users/" + u1.ID.String() + "/projects",
			Body:   map[string]interface{}{},
			Headers: map[string]string{
				"Authorization": "Bearer " + t1.ID.String(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusBadRequest)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Field:   "title",
			Code:    util.RequiredError,
			Message: "Title is required.",
		})
	})

	Convey("User not found", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "POST",
			URL:    "/users/" + uuid.New() + "/projects",
			Headers: map[string]string{
				"Authorization": "Bearer " + t1.ID.String(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusNotFound)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.UserNotFoundError,
			Message: "User not found.",
		})
	})

	Convey("Forbidden", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "POST",
			URL:    "/users/" + u1.ID.String() + "/projects",
			Headers: map[string]string{
				"Authorization": "Bearer " + t2.ID.String(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusForbidden)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.UserForbiddenError,
			Message: "You are forbidden to access this user.",
		})
	})
}

func TestProjectShow(t *testing.T) {
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

	Convey("Success", t, func() {
		project := new(model.Project)
		r := request(&requestOptions{
			Method: "GET",
			URL:    "/projects/" + p1.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + t1.ID.String(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusOK)
		parseJSON(r.Body, project)
		So(project.ID, ShouldResemble, p1.ID)
		So(project.Title, ShouldEqual, p1.Title)
		So(project.Description, ShouldEqual, p1.Description)
		So(project.UserID, ShouldResemble, p1.UserID)
		So(project.CreatedAt.Time, ShouldResemble, p1.CreatedAt.Truncate(time.Second))
		So(project.UpdatedAt.Time, ShouldResemble, p1.UpdatedAt.Truncate(time.Second))
		So(project.IsPrivate, ShouldEqual, p1.IsPrivate)
	})

	Convey("Public project", t, func() {
		project := new(model.Project)
		r := request(&requestOptions{
			Method: "GET",
			URL:    "/projects/" + p1.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + t2.ID.String(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusOK)
		parseJSON(r.Body, project)
		So(project.ID, ShouldResemble, p1.ID)
	})

	Convey("Forbidden", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "GET",
			URL:    "/projects/" + p2.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + t1.ID.String(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusForbidden)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.UserForbiddenError,
			Message: "You are forbidden to access this project.",
		})
	})

	Convey("Project not found", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "GET",
			URL:    "/projects/" + uuid.New(),
		})

		So(r.Code, ShouldEqual, http.StatusNotFound)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.ProjectNotFoundError,
			Message: "Project not found.",
		})
	})
}

func TestProjectUpdate(t *testing.T) {
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
		project := new(model.Project)
		now := time.Now().Truncate(time.Second)
		newTitle := "New title"
		newDesc := "New description"
		r := request(&requestOptions{
			Method: "PUT",
			URL:    "/projects/" + p1.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + t1.ID.String(),
			},
			Body: map[string]interface{}{
				"title":       newTitle,
				"description": newDesc,
				"is_private":  true,
			},
		})

		So(r.Code, ShouldEqual, http.StatusOK)
		parseJSON(r.Body, project)
		So(project.ID, ShouldResemble, p1.ID)
		So(project.Title, ShouldEqual, newTitle)
		So(project.Description, ShouldEqual, newDesc)
		So(project.IsPrivate, ShouldBeTrue)
		So(project.CreatedAt.Time, ShouldResemble, p1.CreatedAt.Truncate(time.Second))
		So(project.UpdatedAt.Time, ShouldHappenOnOrAfter, now)
	})

	Convey("Forbidden", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "PUT",
			URL:    "/projects/" + p1.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + t2.ID.String(),
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

	Convey("Project not found", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "PUT",
			URL:    "/projects/" + uuid.New(),
			Headers: map[string]string{
				"Authorization": "Bearer " + t1.ID.String(),
			},
			Body: map[string]interface{}{},
		})

		So(r.Code, ShouldEqual, http.StatusNotFound)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.ProjectNotFoundError,
			Message: "Project not found.",
		})
	})
}

func TestProjectDestroy(t *testing.T) {
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
		r := request(&requestOptions{
			Method: "DELETE",
			URL:    "/projects/" + p1.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + t1.ID.String(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusNoContent)

		_, err := model.GetProject(p1.ID)
		So(err, ShouldNotBeNil)

		createTestProject(u1, t1, p1, fixtureProjects[0])
	})

	Convey("Forbidden", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "DELETE",
			URL:    "/projects/" + p1.ID.String(),
			Headers: map[string]string{
				"Authorization": "Bearer " + t2.ID.String(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusForbidden)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.UserForbiddenError,
			Message: "You are forbidden to access this user.",
		})
	})

	Convey("Project not found", t, func() {
		err := new(util.APIError)
		r := request(&requestOptions{
			Method: "DELETE",
			URL:    "/projects/" + uuid.New(),
			Headers: map[string]string{
				"Authorization": "Bearer " + t1.ID.String(),
			},
		})

		So(r.Code, ShouldEqual, http.StatusNotFound)
		parseJSON(r.Body, err)
		So(err, ShouldResemble, &util.APIError{
			Code:    util.ProjectNotFoundError,
			Message: "Project not found.",
		})
	})
}
