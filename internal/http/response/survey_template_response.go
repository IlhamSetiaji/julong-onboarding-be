package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
)

type SurveyTemplateResponse struct {
	ID           uuid.UUID                       `json:"id"`
	SurveyNumber string                          `json:"survey_number"`
	Title        string                          `json:"title"`
	Status       entity.SurveyTemplateStatusEnum `json:"status"`
	CreatedAt    time.Time                       `json:"created_at"`
	UpdatedAt    time.Time                       `json:"updated_at"`

	Questions []QuestionResponse `json:"questions"`
}
