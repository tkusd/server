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
	ElementID  *types.UUID       `json:"element_id"`
	Next       *types.UUID       `json:"-"`
	Name       string            `json:"name"`
	Type       types.ElementType `json:"type"`
	CreatedAt  types.Time        `json:"created_at"`
	UpdatedAt  types.Time        `json:"updated_at"`
	Attributes types.JSONObject  `json:"attributes"`

	// From project
	//UserID types.UUID `json:"-" sql:"-"`

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
func (e *Element) BeforeCreate() error {
	e.CreatedAt = types.Now()
	return nil
}

func (e *Element) AfterCreate(tx *gorm.DB) error {
	return tx.Where(&Element{
		ProjectID: e.ProjectID,
		ElementID: e.ElementID,
	}).Where("next is null").Update("next", e.ID.String()).Error
}

func (e *Element) BeforeDelete(tx *gorm.DB) error {
	return tx.Where(&Element{
		Next: &e.ID,
	}).Update("next", e.Next).Error
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
	/*
	   	err := db.Raw(`SELECT e.id, e.project_id, e.element_id, e.next, e.name, e.type, e.created_at, e.updated_at, e.attributes, p.user_id
	   FROM elements AS e
	   JOIN projects AS p ON e.project_id = p.id
	   WHERE e.id = ?
	   LIMIT 1`, id.String()).
	   		Row().
	   		Scan(&e.ID, &e.ProjectID, &e.ElementID, &e.Next, &e.Name, &e.Type, &e.CreatedAt, &e.UpdatedAt, &e.Attributes, &e.UserID)

	   	if err != nil {
	   		return nil, err
	   	}*/

	if err := db.Where("id = ?", id.String()).First(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}

// GetElementList gets a list of elements.
func GetElementList(option *ElementQueryOption) ([]*Element, error) {
	var list []*Element

	query := map[string]interface{}{}

	if option.ProjectID != nil {
		query["project_id"] = option.ProjectID.String()
	}

	if option.ElementID != nil {
		query["element_id"] = option.ElementID.String()
	}

	if err := db.Where(query).Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}
