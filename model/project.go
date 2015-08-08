package model

import (
	"database/sql"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
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
	MainScreen  types.UUID `json:"main_screen"`
	Theme       string     `json:"theme"`

	// Virtual attributes
	Owner struct {
		ID     types.UUID `json:"id"`
		Name   string     `json:"name"`
		Avatar string     `json:"avatar"`
	} `json:"owner,omitempty" sql:"-"`
}

type ProjectCollection struct {
	Data    []*Project `json:"data"`
	HasMore bool       `json:"has_more"`
	Count   int        `json:"count"`
	Limit   int        `json:"limit"`
	Offset  int        `json:"offset"`
}

// ProjectQueryOption is the query options for projects.
type ProjectQueryOption struct {
	QueryOption
	UserID    *types.UUID
	Private   bool
	WithOwner bool
}

// BeforeSave is called when the data is about to be saved.
func (p *Project) BeforeSave(tx *gorm.DB) error {
	if p.MainScreen.Valid() {
		var exist sql.NullBool
		tx.Raw("SELECT exists(SELECT project_id = ? FROM elements WHERE id = ?)", p.ID.String(), p.MainScreen.String()).Row().Scan(&exist)

		if !exist.Bool {
			return &util.APIError{
				Field:   "main_screen",
				Code:    util.ElementNotOwnedByProjectError,
				Message: "Element is not owned by the project.",
			}
		}
	}

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

	if err := db.Save(p).Error; err != nil {
		return err
	}

	var user User

	if err := db.Where("id = ?", p.UserID.String()).Select([]string{"id", "name", "avatar"}).First(&user).Error; err != nil {
		return err
	}

	p.Owner.ID = user.ID
	p.Owner.Name = user.Name
	p.Owner.Avatar = user.Avatar

	return nil
}

// Delete deletes data from the database.
func (p *Project) Delete() error {
	return db.Delete(p).Error
}

// Exists returns true if the record exists.
func (p *Project) Exists() bool {
	return exists("projects", p.ID.String())
}

func generateProjectWithOwnerQuery() *gorm.DB {
	return db.Table("projects").
		Joins("JOIN users ON users.id = projects.user_id").
		Select([]string{
		"projects.id",
		"projects.title",
		"projects.description",
		"projects.user_id",
		"projects.created_at",
		"projects.updated_at",
		"projects.is_private",
		"projects.main_screen",
		"projects.theme",
		"users.id",
		"users.name",
		"users.avatar",
	})
}

func scanProjectsWithOwner(rows *sql.Rows) ([]*Project, error) {
	var list []*Project

	defer rows.Close()

	for rows.Next() {
		project := new(Project)

		err := rows.Scan(
			&project.ID,
			&project.Title,
			&project.Description,
			&project.UserID,
			&project.CreatedAt,
			&project.UpdatedAt,
			&project.IsPrivate,
			&project.MainScreen,
			&project.Theme,
			&project.Owner.ID,
			&project.Owner.Name,
			&project.Owner.Avatar,
		)

		if err != nil {
			return nil, err
		}

		list = append(list, project)
	}

	return list, nil
}

// GetProjectList gets a list of projects.
func GetProjectList(option *ProjectQueryOption) (*ProjectCollection, error) {
	var count int
	query := map[string]interface{}{}

	if !option.Private {
		query["is_private"] = false
	}

	if option.UserID != nil {
		query["user_id"] = option.UserID.String()
	}

	if option.Limit == 0 || option.Limit > maxLimit {
		option.Limit = defaultLimit
	}

	if option.Order == "" {
		option.Order = "-created_at"
	}

	order := option.ParseOrder()

	// Get count
	if err := db.Table("projects").Where(query).Count(&count).Error; err != nil {
		return nil, err
	}

	rows, err := generateProjectWithOwnerQuery().
		Where(query).
		Order(order).
		Offset(option.Offset).
		Limit(option.Limit).
		Rows()

	if err != nil {
		return nil, err
	}

	projects, err := scanProjectsWithOwner(rows)

	if err != nil {
		return nil, err
	}

	if projects == nil {
		projects = make([]*Project, 0)
	}

	return &ProjectCollection{
		Data:    projects,
		Limit:   option.Limit,
		Offset:  option.Offset,
		Count:   count,
		HasMore: count > option.Offset+option.Limit,
	}, nil
}

// GetProject gets the project data.
func GetProject(id types.UUID) (*Project, error) {
	project := new(Project)

	if err := db.Where("id = ?", id.String()).First(project).Error; err != nil {
		return nil, err
	}

	return project, nil
}

// GetProjectWithOwner gets the project with owner data.
func GetProjectWithOwner(id types.UUID) (*Project, error) {
	rows, err := generateProjectWithOwnerQuery().
		Where("projects.id = ?", id.String()).
		Limit(1).
		Rows()

	if err != nil {
		return nil, err
	}

	projects, err := scanProjectsWithOwner(rows)

	if err != nil {
		return nil, err
	}

	if len(projects) == 0 {
		return nil, nil
	}

	return projects[0], nil
}
