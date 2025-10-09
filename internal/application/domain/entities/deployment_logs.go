package entities

import (
	"time"

	"github.com/google/uuid"
)

type DeploymentLogs struct {
	ID            uuid.UUID `json:"id"`
	DeploymentID  uuid.UUID `json:"deployment_id"`
	TokensIn      int       `json:"tokens_in"`
	TokensOut     int       `json:"tokens_out"`
	Input         string    `json:"input"`
	Output        string    `json:"output"`
	DelayTime     int       `json:"delay_time"`
	ExecutionTime int       `json:"execution_time"`
	Source        string    `json:"source"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
