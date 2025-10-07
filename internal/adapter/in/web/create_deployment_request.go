package web

import "github.com/google/uuid"

type CreateDeploymentRequest struct {
	ModelName  string     `json:"model_name" binding:"required"`
	FinetuneID *uuid.UUID `json:"finetune_id"`
}
