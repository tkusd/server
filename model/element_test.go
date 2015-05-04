package model

import (
	"log"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/tkusd/server/model/types"
	"github.com/tkusd/server/util"
)

func createTestElement(project *Project) (*Element, error) {
	element := &Element{
		Name:      "Test element",
		ProjectID: project.ID,
		Type:      types.ElementTypeScreen,
	}

	if err := element.Save(); err != nil {
		return nil, err
	}

	return element, nil
}

func TestElement(t *testing.T) {
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

	Convey("BeforeSave", t, func() {
		now := time.Now()
		element := &Element{}
		element.BeforeSave()
		So(element.UpdatedAt.Time, ShouldHappenOnOrAfter, now)
	})

	Convey("BeforeCreate", t, func() {
		now := time.Now()
		element := &Element{}
		element.BeforeCreate()
		So(element.CreatedAt.Time, ShouldHappenOnOrAfter, now)
	})

	Convey("Save", t, func() {
		Convey("Trim name", func() {
			element := &Element{Name: "   abc   "}
			element.Save()
			So(element.Name, ShouldEqual, "abc")
		})

		Convey("Maximum length of name is 255", func() {
			element := &Element{Name: strings.Repeat("a", 256)}
			err := element.Save()
			So(err, ShouldResemble, &util.APIError{
				Field:   "name",
				Code:    util.LengthError,
				Message: "Maximum length of name is 255.",
			})
		})

		Convey("Type is required", func() {
			element := &Element{}
			err := element.Save()
			So(err, ShouldResemble, &util.APIError{
				Field:   "type",
				Code:    util.RequiredError,
				Message: "Element type is required.",
			})
		})

		Convey("Unsupported element type", func() {
			element := &Element{Type: -1}
			err := element.Save()
			So(err, ShouldResemble, &util.APIError{
				Field:   "type",
				Code:    util.UnsupportedElementTypeError,
				Message: "Unsupported element type.",
			})
		})

		Convey("Create", func() {
			now := time.Now().Truncate(time.Second)
			element := &Element{
				Name:      "test",
				ProjectID: project.ID,
				Type:      types.ElementTypeScreen,
				Attributes: map[string]interface{}{
					"foo": "bar",
				},
			}
			err := element.Save()
			defer element.Delete()

			if err != nil {
				log.Fatal(err)
			}

			e, _ := GetElement(element.ID)
			So(e.ID, ShouldResemble, element.ID)
			So(e.Name, ShouldEqual, element.Name)
			So(e.CreatedAt.Time, ShouldHappenOnOrAfter, now)
			So(e.UpdatedAt.Time, ShouldHappenOnOrAfter, now)
			So(e.ProjectID, ShouldResemble, project.ID)
			So(e.Type, ShouldEqual, element.Type)
			So(e.Attributes, ShouldResemble, element.Attributes)
		})

		Convey("Update", func() {
			element, err := createTestElement(project)
			defer element.Delete()

			if err != nil {
				log.Fatal(err)
			}

			newName := "New name"
			now := time.Now().Truncate(time.Second)
			element.Name = newName

			if err := element.Save(); err != nil {
				log.Fatal(err)
			}

			e, _ := GetElement(element.ID)
			So(e.Name, ShouldEqual, newName)
			So(e.UpdatedAt.Time, ShouldHappenOnOrAfter, now)
		})
	})

	Convey("Delete", t, func() {
		element, err := createTestElement(project)

		if err != nil {
			log.Fatal(err)
		}

		element.Delete()

		e, _ := GetElement(element.ID)
		So(e, ShouldBeNil)
	})
}

func TestGetElement(t *testing.T) {
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

	element, err := createTestElement(project)
	defer element.Delete()

	if err != nil {
		log.Fatal(err)
	}

	Convey("Success", t, func() {
		e, err := GetElement(element.ID)

		if err != nil {
			log.Fatal(err)
		}

		So(e.ID, ShouldResemble, element.ID)
		//So(e.UserID, ShouldResemble, project.UserID)
	})

	Convey("Failed", t, func() {
		e, err := GetElement(types.NewRandomUUID())

		So(e, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})
}

func TestGetElementList(t *testing.T) {
	// TODO: need tests
}
