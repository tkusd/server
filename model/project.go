package model

import (
	"time"

	"code.google.com/p/go-uuid/uuid"
)

type Project struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsPrivate   bool      `json:"is_private"`
}
