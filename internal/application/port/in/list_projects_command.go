package in

import "github.com/google/uuid"

type ListProjectsCommand struct {
	OwnerID uuid.UUID
}