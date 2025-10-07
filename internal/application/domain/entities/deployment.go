package entities

import (
	"time"

	"github.com/google/uuid"
)

type Deployment struct {
	ID         uuid.UUID  `json:"id"`
	ModelName  string     `json:"model_name"`
	APIKey     string     `json:"api_key"`
	ProjectID  uuid.UUID  `json:"project_id"`
	FinetuneID *uuid.UUID `json:"finetune_id,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
