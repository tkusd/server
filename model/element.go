package model

import (
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

	// Virtual attributes
	Elements []*Element `json:"elements,omitempty" sql:"-"`
}

// ElementQueryOption is the query option for elements.
type ElementQueryOption struct {
	ProjectID *types.UUID
	ElementID *types.UUID
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
	query := tx.Table("elements").Select("order_id").Order("order_id desc").Limit(1)

	if !e.ElementID.IsEmpty() {
		query = query.Where("element_id = ?", e.ElementID.String())
	} else if !e.ProjectID.IsEmpty() {
		query = query.Where("project_id = ?", e.ProjectID.String()).Where("element_id IS NULL")
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

func (e *Element) UpdateOrder(arr []types.UUID) error {
	return nil
}

// GetElement gets the element data.
func GetElement(id types.UUID) (*Element, error) {
	e := new(Element)

	if err := db.Where("id = ?", id.String()).First(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}

// GetElementList gets a list of elements.
func GetElementList(option *ElementQueryOption) ([]*Element, error) {
	var list []*Element
	var id string
	var elementID types.UUID

	raw := `WITH RECURSIVE tree AS (
SELECT elements.*, 1 AS depth FROM elements `

	if option.ElementID != nil {
		elementID = *option.ElementID
		id = elementID.String()
		raw += "WHERE element_id = ?"
	} else if option.ProjectID != nil {
		id = option.ProjectID.String()
		raw += "WHERE project_id = ? AND element_id IS NULL"
	}

	raw += ` UNION ALL
SELECT elements.*, tree.depth + 1 FROM elements, tree
WHERE elements.element_id = tree.id
)
SELECT * FROM tree ORDER BY depth, order_id;`

	if err := db.Raw(raw, id).Find(&list).Error; err != nil {
		return nil, err
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

	return result
}
