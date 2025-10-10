package in

import "github.com/google/uuid"

type DownloadDeploymentLogsCommand struct {
	ProjectID    uuid.UUID
	DeploymentID uuid.UUID
	OwnerID      uuid.UUID
}
