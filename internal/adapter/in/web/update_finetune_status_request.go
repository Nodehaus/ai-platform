package web

import "ai-platform/internal/application/domain/entities"

type UpdateFinetuneStatusRequest struct {
	Status entities.FinetuneStatus `json:"status" binding:"required"`
}