package in

import "github.com/google/uuid"

type GetFinetuneCommand struct {
	ProjectID  uuid.UUID `json:"project_id"`
	FinetuneID uuid.UUID `json:"finetune_id"`
	OwnerID    uuid.UUID `json:"owner_id"`
}