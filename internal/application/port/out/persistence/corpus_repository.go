package persistence

import (
	"context"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

type CorpusRepository interface {
	Create(ctx context.Context, corpus *entities.Corpus) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Corpus, error)
	GetByName(ctx context.Context, name string) (*entities.Corpus, error)
	Update(ctx context.Context, corpus *entities.Corpus) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*entities.Corpus, error)
}