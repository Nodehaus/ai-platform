package in

import "github.com/google/uuid"

type CreateProjectCommand struct {
	Name    string
	OwnerID uuid.UUID
}