package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
)

type EducationResponse struct {
	ID             uuid.UUID                 `json:"id"`
	UserProfileID  uuid.UUID                 `json:"user_profile_id"`
	EducationLevel entity.EducationLevelEnum `json:"education_level"`
	Major          string                    `json:"major"`
	SchoolName     string                    `json:"school_name"`
	GraduateYear   int                       `json:"graduate_year"`
	EndDate        time.Time                 `json:"end_date"`
	Certificate    *string                   `json:"certificate"`
	Gpa            *float64                  `json:"gpa"`
	CreatedAt      time.Time                 `json:"created_at"`
	UpdatedAt      time.Time                 `json:"updated_at"`
}
