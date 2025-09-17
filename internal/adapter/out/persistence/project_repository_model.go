package persistence

import (
	"ai-platform/internal/application/domain/entities"
	"time"

	"github.com/google/uuid"
)

type ProjectRepositoryModel struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	OwnerID   uuid.UUID `db:"owner_id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (m *ProjectRepositoryModel) ToEntity() *entities.Project {
	return &entities.Project{
		ID:               m.ID,
		Name:             m.Name,
		OwnerID:          m.OwnerID,
		TrainingDataset:  nil,
		Finetune:         nil,
		Status:           entities.ProjectStatus(m.Status),
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

func FromEntity(project *entities.Project) *ProjectRepositoryModel {
	return &ProjectRepositoryModel{
		ID:        project.ID,
		Name:      project.Name,
		OwnerID:   project.OwnerID,
		Status:    string(project.Status),
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
	}
}