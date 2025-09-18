package persistence

import (
	"time"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

type PromptRepositoryModel struct {
	ID        uuid.UUID `db:"id"`
	Version   int       `db:"version"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (m *PromptRepositoryModel) ToEntity() *entities.Prompt {
	return &entities.Prompt{
		ID:        m.ID,
		Version:   m.Version,
		Text:      m.Text,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func FromPromptEntity(prompt *entities.Prompt) *PromptRepositoryModel {
	return &PromptRepositoryModel{
		ID:        prompt.ID,
		Version:   prompt.Version,
		Text:      prompt.Text,
		CreatedAt: prompt.CreatedAt,
		UpdatedAt: prompt.UpdatedAt,
	}
}