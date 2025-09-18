package in

import "github.com/google/uuid"

type GetProjectCommand struct {
	ProjectID uuid.UUID
	OwnerID   uuid.UUID
}