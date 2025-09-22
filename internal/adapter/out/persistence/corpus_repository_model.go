package persistence

import (
	"time"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

type CorpusRepositoryModel struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	S3Path      string    `db:"s3_path"`
	FilesSubset *[]string `db:"files_subset"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (m *CorpusRepositoryModel) ToEntity() *entities.Corpus {
	return &entities.Corpus{
		ID:          m.ID,
		Name:        m.Name,
		S3Path:      m.S3Path,
		FilesSubset: m.FilesSubset,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func FromCorpusEntity(corpus *entities.Corpus) *CorpusRepositoryModel {
	return &CorpusRepositoryModel{
		ID:          corpus.ID,
		Name:        corpus.Name,
		S3Path:      corpus.S3Path,
		FilesSubset: corpus.FilesSubset,
		CreatedAt:   corpus.CreatedAt,
		UpdatedAt:   corpus.UpdatedAt,
	}
}