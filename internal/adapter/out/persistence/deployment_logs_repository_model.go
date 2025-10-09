package persistence

import (
	"time"

	"ai-platform/internal/application/domain/entities"
	"github.com/google/uuid"
)

type DeploymentLogsRepositoryModel struct {
	ID            uuid.UUID `db:"id"`
	DeploymentID  uuid.UUID `db:"deployment_id"`
	TokensIn      int       `db:"tokens_in"`
	TokensOut     int       `db:"tokens_out"`
	Input         string    `db:"input"`
	Output        string    `db:"output"`
	DelayTime     int       `db:"delay_time"`
	ExecutionTime int       `db:"execution_time"`
	Source        string    `db:"source"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (m *DeploymentLogsRepositoryModel) ToEntity() *entities.DeploymentLogs {
	return &entities.DeploymentLogs{
		ID:            m.ID,
		DeploymentID:  m.DeploymentID,
		TokensIn:      m.TokensIn,
		TokensOut:     m.TokensOut,
		Input:         m.Input,
		Output:        m.Output,
		DelayTime:     m.DelayTime,
		ExecutionTime: m.ExecutionTime,
		Source:        m.Source,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
