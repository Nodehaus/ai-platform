package in

import "github.com/google/uuid"

type GetDeploymentCommand struct {
	DeploymentID uuid.UUID
	ProjectID    uuid.UUID
	OwnerID      uuid.UUID
}
