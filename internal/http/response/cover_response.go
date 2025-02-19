package response

import (
	"time"

	"github.com/google/uuid"
)

type CoverResponse struct {
	ID         uuid.UUID `json:"id"`
	Path       string    `json:"path"`
	PathOrigin string    `json:"path_origin"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
