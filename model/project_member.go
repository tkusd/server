package model

import "code.google.com/p/go-uuid/uuid"

type ProjectMember struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ProjectID uuid.UUID
	Role      int8
}
