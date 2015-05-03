package model

import (
	"time"

	"github.com/tommy351/app-studio-server/model/types"
)

// Element represents the data structure of a element.
type Element struct {
	ID         types.UUID
	ProjectID  types.UUID
	ElementID  types.UUID
	Name       string
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Attributes types.JSONObject
}
