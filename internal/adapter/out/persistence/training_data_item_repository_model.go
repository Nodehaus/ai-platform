package persistence

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

type TrainingDataItemRepositoryModel struct {
	ID                    uuid.UUID  `db:"id"`
	TrainingDatasetID     uuid.UUID  `db:"training_dataset_id"`
	ValuesJSON            string     `db:"values_json"`
	CorrectsID            *uuid.UUID `db:"corrects_id"`
	SourceDocument        *string    `db:"source_document"`
	SourceDocumentStart   *string    `db:"source_document_start"`
	SourceDocumentEnd     *string    `db:"source_document_end"`
	GenerationTimeSeconds float64    `db:"generation_time_seconds"`
	Deleted               bool       `db:"deleted"`
	CreatedAt             time.Time  `db:"created_at"`
	UpdatedAt             time.Time  `db:"updated_at"`
}

func (m *TrainingDataItemRepositoryModel) ToEntity() (*entities.TrainingDataItem, error) {
	var values []string
	if err := json.Unmarshal([]byte(m.ValuesJSON), &values); err != nil {
		return nil, err
	}

	return &entities.TrainingDataItem{
		ID:                    m.ID,
		Values:                values,
		CorrectsID:            m.CorrectsID,
		SourceDocument:        m.SourceDocument,
		SourceDocumentStart:   m.SourceDocumentStart,
		SourceDocumentEnd:     m.SourceDocumentEnd,
		GenerationTimeSeconds: m.GenerationTimeSeconds,
		Deleted:               m.Deleted,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
	}, nil
}

func FromTrainingDataItemEntity(tdi *entities.TrainingDataItem, trainingDatasetID uuid.UUID) (*TrainingDataItemRepositoryModel, error) {
	valuesJSON, err := json.Marshal(tdi.Values)
	if err != nil {
		return nil, err
	}

	return &TrainingDataItemRepositoryModel{
		ID:                    tdi.ID,
		TrainingDatasetID:     trainingDatasetID,
		ValuesJSON:            string(valuesJSON),
		CorrectsID:            tdi.CorrectsID,
		SourceDocument:        tdi.SourceDocument,
		SourceDocumentStart:   tdi.SourceDocumentStart,
		SourceDocumentEnd:     tdi.SourceDocumentEnd,
		GenerationTimeSeconds: tdi.GenerationTimeSeconds,
		Deleted:               tdi.Deleted,
		CreatedAt:             tdi.CreatedAt,
		UpdatedAt:             tdi.UpdatedAt,
	}, nil
}