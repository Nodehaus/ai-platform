package web

import "ai-platform/internal/application/domain/entities"

type UpdateTrainingDatasetStatusRequest struct {
	Status entities.TrainingDatasetStatus `json:"status" binding:"required"`
}