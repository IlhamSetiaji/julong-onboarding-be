package response

import (
	"time"

	"github.com/google/uuid"
)

type SkillResponse struct {
	ID            uuid.UUID `json:"id"`
	UserProfileID uuid.UUID `json:"user_profile_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Certificate   *string   `json:"certificate"`
	Level         *int      `json:"level"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
