package model

import (
	"database/sql"
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
	ID         types.UUID       `json:"id"`
	ProjectID  types.UUID       `json:"project_id"`
	ElementID  types.UUID       `json:"element_id"`
	OrderID    int              `json:"order_id"`
	Name       string           `json:"name"`
	Type       string           `json:"type"`
	CreatedAt  types.Time       `json:"created_at"`
	UpdatedAt  types.Time       `json:"updated_at"`
	Attributes types.JSONObject `json:"attributes"`
	Styles     types.JSONObject `json:"styles"`
	Events     types.JSONArray  `json:"events"`
	IsVisible  bool             `json:"is_visible"`

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

func (e *Element) BeforeDelete(tx *gorm.DB) error {
	if e.Type != "screen" {
		return nil
	}

	var project Project

	if err := tx.Where("id = ?", e.ProjectID.String()).Select([]string{"id", "main_screen"}).First(&project).Error; err != nil {
		return err
	}

	if project.MainScreen.Equal(e.ID) {
		if err := tx.Exec(`UPDATE projects
SET main_screen = (SELECT id FROM elements WHERE project_id = ? AND element_id IS NULL AND id <> ? ORDER BY order_id LIMIT 1)
WHERE id = ?`, project.ID.String(), e.ID.String(), project.ID.String()).Error; err != nil {
			return err
		}
	}

	return nil
}

func (e *Element) AfterCreate(tx *gorm.DB) error {
	if e.Type != "screen" {
		return nil
	}

	var project Project

	if err := tx.Where("id = ?", e.ProjectID.String()).Select([]string{"id", "main_screen"}).First(&project).Error; err != nil {
		return err
	}

	if !project.MainScreen.Valid() {
		project.MainScreen = e.ID
		tx.Model(&project).UpdateColumns(project)
	}

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

	return db.Save(e).Error
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

	if list == nil {
		list = make([]*Element, 0)
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
	// Check whether all elements are owned by the same project
	projectID := *option.ProjectID
	tx := db.Begin()

	for _, elementID := range elements {
		var exist sql.NullBool
		tx.Raw("SELECT exists(SELECT 1 FROM elements WHERE id = ? AND project_id = ?)", elementID.String(), projectID.String()).Row().Scan(&exist)

		if !exist.Bool {
			tx.Rollback()
			return &util.APIError{
				Code:    util.ElementNotFoundError,
				Message: fmt.Sprintf("Element %s does not exist or is not the children of project %s", elementID.String(), projectID.String()),
			}
		}
	}

	// Start updating the order of elements
	var parentID types.UUID

	if option.ElementID != nil {
		parentID = *option.ElementID
	}

	for i, elementID := range elements {
		data := map[string]interface{}{
			"element_id": parentID,
			"order_id":   i + 1,
		}

		if err := tx.Table("elements").Where("id = ?", elementID.String()).UpdateColumns(data).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	tx.Commit()

	return nil
}
