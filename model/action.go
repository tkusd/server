package model

import (
	"github.com/asaskevich/govalidator"
	"github.com/tkusd/server/model/types"
	"github.com/tkusd/server/util"
)

type Action struct {
	ID        types.UUID       `json:"id"`
	Name      string           `json:"name"`
	ProjectID types.UUID       `json:"project_id"`
	Action    string           `json:"action"`
	Data      types.JSONObject `json:"data"`
	CreatedAt types.Time       `json:"created_at"`
	UpdatedAt types.Time       `json:"updated_at"`
}

func (action *Action) Save() error {
	action.Name = govalidator.Trim(action.Name, "")

	if len(action.Name) > 255 {
		return &util.APIError{
			Field:   "name",
			Code:    util.LengthError,
			Message: "Maximum length of name is 255.",
		}
	}

	if action.Action == "" {
		return &util.APIError{
			Field:   "action",
			Code:    util.RequiredError,
			Message: "Action is required.",
		}
	}

	return db.Save(action).Error
}

func (action *Action) Delete() error {
	return db.Delete(action).Error
}

func GetAction(id types.UUID) (*Action, error) {
	action := new(Action)

	if err := db.Where("id = ?", id.String()).First(action).Error; err != nil {
		return nil, err
	}

	return action, nil
}

func GetActionList(projectID types.UUID) ([]*Action, error) {
	var list []*Action

	if err := db.Where("project_id = ?", projectID.String()).Find(&list).Error; err != nil {
		return nil, err
	}

	if list == nil {
		list = make([]*Action, 0)
	}

	return list, nil
}
