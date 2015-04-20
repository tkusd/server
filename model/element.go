package model

import (
	"time"

	"github.com/jmoiron/sqlx/types"

	"code.google.com/p/go-uuid/uuid"
)

type Element struct {
	ID         uuid.UUID
	ProjectID  uuid.UUID
	ElementID  uuid.UUID
	Name       string
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Attributes types.JsonText
}
