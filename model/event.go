package model

import (
	"github.com/lib/pq"
	"github.com/tkusd/server/model/types"
	"github.com/tkusd/server/util"
)

type Event struct {
	ID        types.UUID `json:"id"`
	ElementID types.UUID `json:"element_id"`
	ActionID  types.UUID `json:"action_id"`
	Event     string     `json:"event"`
	CreatedAt types.Time `json:"created_at"`
	UpdatedAt types.Time `json:"updated_at"`
}

func (event *Event) Save() error {
	if event.Event == "" {
		return &util.APIError{
			Field:   "event",
			Code:    util.RequiredError,
			Message: "Event is required.",
		}
	}

	if !event.ActionID.Valid() {
		return &util.APIError{
			Field:   "action_id",
			Code:    util.RequiredError,
			Message: "Action ID is required.",
		}
	}

	err := db.Save(event).Error

	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *pq.Error:
		switch e.Code.Name() {
		case ForeignKeyViolation:
			return &util.APIError{
				Code:    util.ActionNotOwnedByProjectError,
				Message: "Action is not owned by the project.",
				Field:   "action_id",
			}
		}
	}

	return err
}

func (event *Event) Delete() error {
	return db.Delete(event).Error
}

func GetEvent(id types.UUID) (*Event, error) {
	event := new(Event)

	if err := db.Where("id = ?", id.String()).First(event).Error; err != nil {
		return nil, err
	}

	return event, nil
}

func GetEventList(elementID types.UUID) ([]*Event, error) {
	var list []*Event

	if err := db.Where("element_id = ?", elementID.String()).Find(&list).Error; err != nil {
		return nil, err
	}

	if list == nil {
		list = make([]*Event, 0)
	}

	return list, nil
}
