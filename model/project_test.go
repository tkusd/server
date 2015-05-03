package model

import (
	"strings"
	"testing"

	"log"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/tommy351/app-studio-server/model/types"
	"github.com/tommy351/app-studio-server/util"
)

func createTestProject(user *User) (*Project, error) {
	project := &Project{
		Title:  "Test project",
		UserID: user.ID,
	}

	if err := project.Save(); err != nil {
		return nil, err
	}

	return project, nil
}

func TestProject(t *testing.T) {
	user, err := createTestUser(fixtureUsers[0])
	defer user.Delete()

	if err != nil {
		log.Fatal(err)
	}

	Convey("BeforeSave", t, func() {
		now := time.Now()
		project := &Project{}
		project.BeforeSave()
		So(project.UpdatedAt.Time, ShouldHappenOnOrAfter, now)
	})

	Convey("BeforeCreate", t, func() {
		now := time.Now()
		project := &Project{}
		project.BeforeCreate()
		So(project.CreatedAt.Time, ShouldHappenOnOrAfter, now)
	})

	Convey("Save", t, func() {
		Convey("Trim title", func() {
			project := &Project{Title: "   abc   "}
			project.Save()
			defer project.Delete()

			So(project.Title, ShouldEqual, "abc")
		})

		Convey("Title is required", func() {
			project := &Project{}
			err := project.Save()
			So(err, ShouldResemble, &util.APIError{
				Field:   "title",
				Code:    util.RequiredError,
				Message: "Title is required.",
			})
		})

		Convey("Maximum length of title is 255", func() {
			project := &Project{Title: strings.Repeat("a", 256)}
			err := project.Save()
			So(err, ShouldResemble, &util.APIError{
				Field:   "title",
				Code:    util.LengthError,
				Message: "Maximum length of title is 255.",
			})
		})

		Convey("Create", func() {
			now := time.Now().Truncate(time.Second)
			project := &Project{
				Title:  "test",
				UserID: user.ID,
			}
			err := project.Save()
			defer project.Delete()

			if err != nil {
				log.Fatal(err)
			}

			p, _ := GetProject(project.ID)
			So(p.ID, ShouldResemble, project.ID)
			So(p.Title, ShouldEqual, project.Title)
			So(p.CreatedAt.Time, ShouldHappenOnOrAfter, now)
			So(p.UpdatedAt.Time, ShouldHappenOnOrAfter, now)
			So(p.UserID, ShouldResemble, user.ID)
		})

		Convey("Update", func() {
			project, err := createTestProject(user)
			defer project.Delete()

			if err != nil {
				log.Fatal(err)
			}

			newTitle := "new title"
			now := time.Now().Truncate(time.Second)
			project.Title = newTitle
			err = project.Save()

			if err != nil {
				log.Fatal(err)
			}

			p, _ := GetProject(project.ID)
			So(p.Title, ShouldEqual, newTitle)
			So(p.UpdatedAt.Time, ShouldHappenOnOrAfter, now)
		})
	})

	Convey("Delete", t, func() {
		project, err := createTestProject(user)

		if err != nil {
			log.Fatal(err)
		}

		project.Delete()

		p, _ := GetProject(project.ID)
		So(p, ShouldBeNil)
	})

	Convey("Exists", t, func() {
		Convey("true", func() {
			project, err := createTestProject(user)
			defer project.Delete()

			if err != nil {
				log.Fatal(err)
			}

			So(project.Exists(), ShouldBeTrue)
		})

		Convey("false", func() {
			project := &Project{ID: types.NewRandomUUID()}
			So(project.Exists(), ShouldBeFalse)
		})
	})
}

func TestGetProjectList(t *testing.T) {
	// TODO: need tests
}

func TestGetProject(t *testing.T) {
	user, err := createTestUser(fixtureUsers[0])
	defer user.Delete()

	if err != nil {
		log.Fatal(err)
	}

	project, err := createTestProject(user)
	defer project.Delete()

	if err != nil {
		log.Fatal(err)
	}

	Convey("Success", t, func() {
		p, err := GetProject(project.ID)

		if err != nil {
			log.Fatal(err)
		}

		So(p.ID, ShouldResemble, project.ID)
	})

	Convey("Failed", t, func() {
		p, err := GetProject(types.NewRandomUUID())

		So(p, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})
}
