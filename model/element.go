package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/tkusd/server/model/types"
	"github.com/tkusd/server/util"
)

// Element represents the data structure of a element.
type Element struct {
	ID         types.UUID        `json:"id"`
	ProjectID  types.UUID        `json:"project_id"`
	ElementID  types.UUID        `json:"element_id"`
	OrderID    int               `json:"-"`
	Name       string            `json:"name"`
	Type       types.ElementType `json:"type"`
	CreatedAt  types.Time        `json:"created_at"`
	UpdatedAt  types.Time        `json:"updated_at"`
	Attributes types.JSONObject  `json:"attributes"`
	Styles     types.JSONObject  `json:"styles"`

	// Virtual attributes
	Elements []*Element `json:"elements,omitempty" sql:"-"`
}

// ElementQueryOption is the query option for elements.
type ElementQueryOption struct {
	ProjectID *types.UUID
	ElementID *types.UUID
	Flat      bool
	Depth     uint
	Select    []string
}

// BeforeSave is called when the data is about to be saved.
func (e *Element) BeforeSave() error {
	e.UpdatedAt = types.Now()
	return nil
}

// BeforeCreate is called when the data is about to be created.
func (e *Element) BeforeCreate(tx *gorm.DB) error {
	e.CreatedAt = types.Now()
	lastOrder := 0
	query := tx.Table("elements").
		Select("order_id").
		Order("order_id desc").
		Limit(1)

	if e.ElementID.Valid() {
		query = query.Where("element_id = ?", e.ElementID.String())
	} else if e.ProjectID.Valid() {
		query = query.Where("project_id = ?", e.ProjectID.String()).
			Where("element_id is null")
	} else {
		return errors.New("UUID is invalid")
	}

	query.Row().Scan(&lastOrder)
	e.OrderID = lastOrder + 1

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

	if e.Type == 0 {
		return &util.APIError{
			Field:   "type",
			Code:    util.RequiredError,
			Message: "Element type is required.",
		}
	}

	if e.Type.String() == "" {
		return &util.APIError{
			Field:   "type",
			Code:    util.UnsupportedElementTypeError,
			Message: "Unsupported element type.",
		}
	}

	return db.Save(e).Error
}

// Delete deletes data from the database.
func (e *Element) Delete() error {
	return db.Delete(e).Error
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
		//option.Select = appendIfMissing(option.Select, "id", "project_id", "element_id", "order_id")
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
SELECT * FROM tree ORDER BY depth, order_id;`

	if err := db.Raw(raw, id).Find(&list).Error; err != nil {
		return nil, err
	}

	if option.Flat {
		return list, nil
	}

	if list == nil {
		list = make([]*Element, 0)
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

func UpdateElementOrder(option *ElementQueryOption, tree []ElementTreeItem) error {
	if err := checkElementTree(option, tree); err != nil {
		return err
	}

	var parentID types.UUID
	tx := db.Begin()

	if option.ElementID != nil {
		parentID = *option.ElementID
	}

	if err := updateElementOrderInTx(tx, parentID, tree); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func checkElementTree(option *ElementQueryOption, tree []ElementTreeItem) error {
	children, err := GetElementList(&ElementQueryOption{
		ProjectID: option.ProjectID,
		ElementID: option.ElementID,
		Flat:      true,
		Select:    []string{"id", "order_id"}, // order_id is needed for sorting
	})

	if err != nil {
		return err
	}

	ids := pickElementIDFromTree(tree)

	if len(ids) != len(children) {
		return &util.APIError{
			Code:    util.ElementTreeNotCompletedError,
			Message: "You didn't provide the full list of children.",
		}
	}

	for _, id := range ids {
		found := false

		for _, elem := range children {
			if id.Equal(elem.ID) {
				found = true
				break
			}
		}

		if !found {
			return &util.APIError{
				Code:    util.ElementNotInTreeError,
				Message: fmt.Sprintf("Element %s is not a child of the specified element.", id.String()),
			}
		}
	}

	return nil
}

func pickElementIDFromTree(tree []ElementTreeItem) []types.UUID {
	var list []types.UUID

	for _, item := range tree {
		list = append(list, item.ID)
		list = append(list, pickElementIDFromTree(item.Elements)...)
	}

	return list
}

func updateElementOrderInTx(tx *gorm.DB, parentID types.UUID, tree []ElementTreeItem) error {
	for i, item := range tree {
		data := map[string]interface{}{
			"element_id": parentID,
			"order_id":   i + 1,
		}

		if err := tx.Table("elements").Where("id = ?", item.ID.String()).UpdateColumns(data).Error; err != nil {
			return err
		}

		if err := updateElementOrderInTx(tx, item.ID, item.Elements); err != nil {
			return err
		}
	}

	return nil
}
