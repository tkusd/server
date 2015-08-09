package model

import "github.com/tkusd/server/model/types"

type PermissionFlag byte

const (
	Readable PermissionFlag = 1 << iota
	Writable
	Admin
)

type Collaborator struct {
	ID         types.UUID
	UserID     types.UUID
	ProjectID  types.UUID
	Permission PermissionFlag
	CreatedAt  types.Time
	UpdatedAt  types.Time
}
