package in

import (
	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
)

type UpdateFinetuneStatusCommand struct {
	FinetuneID uuid.UUID              `json:"finetune_id"`
	Status     entities.FinetuneStatus `json:"status"`
}