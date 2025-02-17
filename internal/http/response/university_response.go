package response

import (
	"time"

	"github.com/google/uuid"
)

type UniversityResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Country      string    `json:"country"`
	AlphaTwoCode string    `json:"alpha_two_code"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
