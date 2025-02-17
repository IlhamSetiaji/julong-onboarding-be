package response

import (
	"time"

	"github.com/google/uuid"
)

type WorkExperienceResponse struct {
	ID             uuid.UUID `json:"id"`
	UserProfileID  uuid.UUID `json:"user_profile_id"`
	Name           string    `json:"name"`
	CompanyName    string    `json:"company_name"`
	YearExperience int       `json:"year_experience"`
	JobDescription string    `json:"job_description"`
	Certificate    *string   `json:"certificate"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
