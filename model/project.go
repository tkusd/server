package model

import (
	"github.com/asaskevich/govalidator"
	"github.com/tkusd/server/model/types"
	"github.com/tkusd/server/util"
)

// Project represents the data structure of a project.
type Project struct {
	ID          types.UUID `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	UserID      types.UUID `json:"user_id"`
	CreatedAt   types.Time `json:"created_at"`
	UpdatedAt   types.Time `json:"updated_at"`
	IsPrivate   bool       `json:"is_private"`
}

// ProjectQueryOption is the query options for projects.
type ProjectQueryOption struct {
	QueryOption
	UserID  *types.UUID
	Private bool
}

// BeforeSave is called when the data is about to be saved.
func (p *Project) BeforeSave() error {
	p.UpdatedAt = types.Now()
	return nil
}

// BeforeCreate is called when the data is about to be created.
func (p *Project) BeforeCreate() error {
	p.CreatedAt = types.Now()
	return nil
}

// Save creates or updates data in the database.
func (p *Project) Save() error {
	p.Title = govalidator.Trim(p.Title, "")

	if p.Title == "" {
		return &util.APIError{
			Field:   "title",
			Code:    util.RequiredError,
			Message: "Title is required.",
		}
	}

	if len(p.Title) > 255 {
		return &util.APIError{
			Field:   "title",
			Code:    util.LengthError,
			Message: "Maximum length of title is 255.",
		}
	}

	return db.Save(p).Error
}

// Delete deletes data from the database.
func (p *Project) Delete() error {
	return db.Delete(p).Error
}

// Exists returns true if the record exists.
func (p *Project) Exists() bool {
	return exists("projects", p.ID.String())
}

/*
func (p *Project) Elements() ([]*Element, error) {

	var list []*Element
	var result []*Element

	if err := db.Where("project_id = ?", p.ID.String()).Find(&list).Error; err != nil {
		return nil, err
	}

	for i, item := range list {
		//
	}

	return result, nil

	db.Raw(``, p.ID.String())
}*/

// GetProjectList gets a list of projects.
func GetProjectList(option *ProjectQueryOption) ([]*Project, error) {
	var list []*Project

	query := map[string]interface{}{}

	if !option.Private {
		query["is_private"] = false
	}

	if option.UserID != nil {
		query["user_id"] = option.UserID.String()
	}

	if option.Limit == 0 {
		option.Limit = defaultLimit
	}

	if option.Order == "" {
		option.Order = "created_at desc"
	}

	if err := db.Where(query).Order(option.Order).Offset(option.Offset).Limit(option.Limit).Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

// GetProject gets the project data.
func GetProject(id types.UUID) (*Project, error) {
	project := new(Project)

	if err := db.Where("id = ?", id.String()).First(project).Error; err != nil {
		return nil, err
	}

	return project, nil
}
