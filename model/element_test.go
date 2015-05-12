package model

import (
	"fmt"
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

func createTestChildElement(element *Element) (*Element, error) {
	e := &Element{
		Name:      "Test child element",
		ProjectID: element.ProjectID,
		ElementID: element.ID,
		Type:      types.ElementTypeText,
	}

	if err := e.Save(); err != nil {
		return nil, err
	}

	return e, nil
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
			So(e.OrderID, ShouldEqual, 1)
		})

		Convey("Order should be increased", func() {
			e1 := &Element{
				Name:      "element 1",
				ProjectID: project.ID,
				Type:      types.ElementTypeScreen,
			}
			defer e1.Delete()

			if err := e1.Save(); err != nil {
				log.Fatal(err)
			}

			e2 := &Element{
				Name:      "element 2",
				ProjectID: project.ID,
				Type:      types.ElementTypeScreen,
			}
			defer e2.Delete()

			if err := e2.Save(); err != nil {
				log.Fatal(err)
			}

			So(e1.OrderID, ShouldEqual, 1)
			So(e2.OrderID, ShouldEqual, 2)
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
	})

	Convey("Failed", t, func() {
		e, err := GetElement(types.NewRandomUUID())

		So(e, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})
}

func TestGetElementList(t *testing.T) {
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

	/**
	project
	|- e1
		|- e2
		|- e3
			|- e4
	*/

	e1, err := createTestElement(project)
	defer e1.Delete()

	if err != nil {
		log.Fatal(err)
	}

	e2, err := createTestChildElement(e1)
	defer e2.Delete()

	if err != nil {
		log.Fatal(err)
	}

	e3, err := createTestChildElement(e1)
	defer e3.Delete()

	if err != nil {
		log.Fatal(err)
	}

	e4, err := createTestChildElement(e3)
	defer e4.Delete()

	if err != nil {
		log.Fatal(err)
	}

	Convey("Project", t, func() {
		list, err := GetElementList(&ElementQueryOption{
			ProjectID: &project.ID,
		})

		if err != nil {
			log.Fatal(err)
		}

		So(list[0].ID, ShouldResemble, e1.ID)
		So(list[0].Elements[0].ID, ShouldResemble, e2.ID)
		So(list[0].Elements[1].ID, ShouldResemble, e3.ID)
		So(list[0].Elements[1].Elements[0].ID, ShouldResemble, e4.ID)
	})

	Convey("Element", t, func() {
		list, err := GetElementList(&ElementQueryOption{
			ElementID: &e1.ID,
		})

		if err != nil {
			log.Fatal(err)
		}

		So(list[0].ID, ShouldResemble, e2.ID)
		So(list[1].ID, ShouldResemble, e3.ID)
		So(list[1].Elements[0].ID, ShouldResemble, e4.ID)
	})

	Convey("Flat", t, func() {
		list, err := GetElementList(&ElementQueryOption{
			ProjectID: &project.ID,
			Flat:      true,
		})

		if err != nil {
			log.Fatal(err)
		}

		So(list[0].ID, ShouldResemble, e1.ID)
		So(list[1].ID, ShouldResemble, e2.ID)
		So(list[2].ID, ShouldResemble, e3.ID)
		So(list[3].ID, ShouldResemble, e4.ID)
	})

	Convey("Depth", t, func() {
		list, err := GetElementList(&ElementQueryOption{
			ProjectID: &project.ID,
			Depth:     2,
		})

		if err != nil {
			log.Fatal(err)
		}

		So(len(list[0].Elements[1].Elements), ShouldEqual, 0)
	})
}

func TestUpdateElementOrder(t *testing.T) {
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

	/**
	project
	|- e1
		|- e2
		|- e3
			|- e4
	*/

	e1, err := createTestElement(project)
	defer e1.Delete()

	if err != nil {
		log.Fatal(err)
	}

	e2, err := createTestChildElement(e1)
	defer e2.Delete()

	if err != nil {
		log.Fatal(err)
	}

	e3, err := createTestChildElement(e1)
	defer e3.Delete()

	if err != nil {
		log.Fatal(err)
	}

	e4, err := createTestChildElement(e3)
	defer e4.Delete()

	if err != nil {
		log.Fatal(err)
	}

	Convey("Not one of the children", t, func() {
		option := &ElementQueryOption{
			ElementID: &e3.ID,
		}
		tree := []ElementTreeItem{
			ElementTreeItem{ID: e1.ID},
		}
		err := UpdateElementOrder(option, tree)

		So(err, ShouldResemble, &util.APIError{
			Code:    util.ElementNotInTreeError,
			Message: fmt.Sprintf("Element %s is not a child of the specified element.", e1.ID.String()),
		})
	})

	Convey("Not complete list of children", t, func() {
		option := &ElementQueryOption{
			ElementID: &e1.ID,
		}
		tree := []ElementTreeItem{
			ElementTreeItem{ID: e2.ID},
			ElementTreeItem{ID: e4.ID},
		}
		err := UpdateElementOrder(option, tree)

		So(err, ShouldResemble, &util.APIError{
			Code:    util.ElementTreeNotCompletedError,
			Message: "You didn't provide the full list of children.",
		})
	})

	Convey("Success", t, func() {
		option := &ElementQueryOption{
			ElementID: &e1.ID,
		}
		tree := []ElementTreeItem{
			ElementTreeItem{ID: e3.ID},
			ElementTreeItem{ID: e4.ID},
			ElementTreeItem{ID: e2.ID},
		}
		err := UpdateElementOrder(option, tree)

		if err != nil {
			log.Fatal(err)
		}

		list, _ := GetElementList(option)

		So(len(list), ShouldEqual, 3)
		So(list[0].ID, ShouldResemble, e3.ID)
		So(list[1].ID, ShouldResemble, e4.ID)
		So(list[1].ElementID, ShouldResemble, e1.ID)
		So(list[2].ID, ShouldResemble, e2.ID)
	})
}
