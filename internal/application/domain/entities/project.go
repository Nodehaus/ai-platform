package entities

import (
	"time"

	"github.com/google/uuid"
)

type ProjectStatus string

const (
	ProjectStatusActive   ProjectStatus = "ACTIVE"
	ProjectStatusArchived ProjectStatus = "ARCHIVED"
	ProjectStatusDeleted  ProjectStatus = "DELETED"
)

type Project struct {
	ID               uuid.UUID     `json:"id"`
	Name             string        `json:"name"`
	OwnerID          uuid.UUID     `json:"owner_id"`
	TrainingDataset  interface{}   `json:"training_dataset"`
	Finetune         interface{}   `json:"finetune"`
	Status           ProjectStatus `json:"status"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}