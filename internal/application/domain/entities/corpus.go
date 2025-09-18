package entities

import (
	"time"

	"github.com/google/uuid"
)

type Corpus struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	S3Path    string    `json:"s3_path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}