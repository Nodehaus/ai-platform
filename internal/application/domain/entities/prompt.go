package entities

import (
	"time"

	"github.com/google/uuid"
)

type Prompt struct {
	ID        uuid.UUID `json:"id"`
	Version   int       `json:"version"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}