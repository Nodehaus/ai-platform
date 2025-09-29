package in

import "github.com/google/uuid"

type DownloadModelCommand struct {
	ProjectID  uuid.UUID `json:"project_id"`
	FinetuneID uuid.UUID `json:"finetune_id"`
}