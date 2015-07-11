package model

import "github.com/tkusd/server/model/types"

type Asset struct {
	ID          types.UUID `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ProjectID   types.UUID `json:"project_id"`
	CreatedAt   types.Time `json:"created_at"`
	UpdatedAt   types.Time `json:"updated_at"`
}
