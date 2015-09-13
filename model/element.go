package model

import (
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/tkusd/server/model/types"
	"github.com/tkusd/server/util"
)

// Element represents the data structure of a element.
type Element struct {
	ID         types.UUID       `json:"id"`
	ProjectID  types.UUID       `json:"project_id"`
	ElementID  types.UUID       `json:"element_id"`
	Index      int              `json:"index"`
	Name       string           `json:"name"`
	Type       string           `json:"type"`
	CreatedAt  types.Time       `json:"created_at"`
	UpdatedAt  types.Time       `json:"updated_at"`
	Attributes types.JSONObject `json:"attributes"`
	Styles     types.JSONObject `json:"styles"`
	IsVisible  bool             `json:"is_visible"`

	// Virtual attributes
	Elements []*Element `json:"elements,omitempty" sql:"-"`
	Events   []*Event   `json:"events,omitempty" sql:"-"`
}

// ElementQueryOption is the query option for elements.
type ElementQueryOption struct {
	ProjectID  *types.UUID
	ElementID  *types.UUID
	Flat       bool
	Depth      uint
	Select     []string
	WithEvents bool
}

func (e *Element) AfterCreate(tx *gorm.DB) error {
	tx.Table("elements").
		Select("index").
		Where("id = ?", e.ID.String()).
		Row().
		Scan(&e.Index)

	return nil
}

// Save creates or updates data in the database.
func (e *Element) Save() error {
	e.Name = govalidator.Trim(e.Name, "")

	if len(e.Name) > 255 {
		return &util.APIError{
			Field:   "name",
			Code:    util.LengthError,
			Message: "Maximum length of name is 255.",
		}
	}

	if e.Type == "" {
		return &util.APIError{
			Field:   "type",
			Code:    util.RequiredError,
			Message: "Element type is required.",
		}
	}

	if e.Attributes == nil {
		e.Attributes = map[string]interface{}{}
	}

	if e.Styles == nil {
		e.Styles = map[string]interface{}{}
	}

	err := db.Save(e).Error

	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *pq.Error:
		switch e.Code.Name() {
		case ForeignKeyViolation:
			return &util.APIError{
				Code:    util.ElementNotOwnedByProjectError,
				Message: "The parent element is not owned by the project.",
				Field:   "element_id",
			}
		}
	}

	return err
}

// Delete deletes data from the database.
func (e *Element) Delete() error {
	return db.Delete(e).Error
}

func (e *Element) Exists() bool {
	return exists("elements", e.ID.String())
}

// GetElement gets the element data.
func GetElement(id types.UUID) (*Element, error) {
	e := new(Element)

	if err := db.Where("id = ?", id.String()).First(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func appendIfMissing(arr []string, key ...string) []string {
	return arr
}

// GetElementList gets a list of elements.
func GetElementList(option *ElementQueryOption) ([]*Element, error) {
	var list []*Element
	var id string
	var elementID types.UUID
	var columns []string

	if len(option.Select) == 0 {
		option.Select = []string{"*"}
	} else {
		// TODO: Add missing fields
		//option.Select = appendIfMissing(option.Select, "id", "project_id", "element_id", "index")
	}

	for _, col := range option.Select {
		columns = append(columns, "elements."+col)
	}

	selectColumns := strings.Join(columns, ",")

	raw := `WITH RECURSIVE tree AS (
SELECT ` + selectColumns + `, 1 AS depth FROM elements WHERE `

	if option.ElementID != nil {
		elementID = *option.ElementID
		id = elementID.String()
		raw += "element_id = ?"
	} else if option.ProjectID != nil {
		id = option.ProjectID.String()
		raw += "project_id = ? AND element_id IS NULL"
	}

	raw += ` UNION ALL
SELECT ` + selectColumns + `, tree.depth + 1 FROM elements, tree
WHERE elements.element_id = tree.id`

	if option.Depth > 0 {
		raw += " AND tree.depth < " + strconv.Itoa(int(option.Depth))
	}

	raw += `)
SELECT * FROM tree ORDER BY depth, index;`

	if err := db.Raw(raw, id).Find(&list).Error; err != nil {
		return nil, err
	}

	if list == nil {
		list = make([]*Element, 0)
	}

	if option.WithEvents {
		var events []*Event

		err := db.Select([]string{
			"events.id",
			"events.element_id",
			"events.workspace",
			"events.event",
			"events.created_at",
			"events.updated_at",
		}).
			Joins("JOIN elements ON events.element_id = elements.id").
			Order("created_at").
			Where("project_id = ?", option.ProjectID.String()).
			Find(&events).
			Error

		if err != nil {
			return nil, err
		}

		for _, event := range events {
			for _, element := range list {
				if event.ElementID.Equal(element.ID) {
					element.Events = append(element.Events, event)
				}
			}
		}
	}

	if option.Flat {
		return list, nil
	}

	return buildElementTree(list, elementID), nil
}

func buildElementTree(list []*Element, parentID types.UUID) []*Element {
	var result []*Element

	for i, item := range list {
		if parentID.Equal(item.ElementID) {
			result = append(result, item)
			item.Elements = buildElementTree(list[i:], item.ID)
		}
	}

	if result == nil {
		result = make([]*Element, 0)
	}

	return result
}

func UpdateElementOrder(option *ElementQueryOption, elements []types.UUID) error {
	var parentID types.UUID
	tx := db.Begin()

	if option.ElementID != nil {
		parentID = *option.ElementID
	}

	for i, elementID := range elements {
		data := map[string]interface{}{
			"element_id": parentID,
			"index":      i + 1,
		}

		if err := tx.Table("elements").Where("id = ?", elementID.String()).UpdateColumns(data).Error; err != nil {
			tx.Rollback()

			switch e := err.(type) {
			case *pq.Error:
				switch e.Code.Name() {
				case ForeignKeyViolation:
					return &util.APIError{
						Code:    util.ElementNotOwnedByProjectError,
						Message: "Elements is not owned by the project.",
					}
				}
			}

			return err
		}
	}

	// Commit the transaction
	tx.Commit()

	return nil
}

func GetProjectIDForElement(elementID types.UUID) types.UUID {
	var projectID types.UUID
	db.Raw("SELECT project_id FROM elements WHERE id = ?", elementID.String()).Row().Scan(&projectID)
	return projectID
}
