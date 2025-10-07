package in

import "github.com/google/uuid"

type CreateDeploymentCommand struct {
	ModelName  string
	ProjectID  uuid.UUID
	FinetuneID *uuid.UUID
	OwnerID    uuid.UUID
}
