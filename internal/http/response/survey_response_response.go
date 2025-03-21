package response

import (
	"time"

	"github.com/google/uuid"
)

type SurveyResponseResponse struct {
	ID               uuid.UUID `json:"id"`
	SurveyTemplateID uuid.UUID `json:"survey_template_id"`
	EmployeeTaskID   uuid.UUID `json:"employee_task_id"`
	QuestionID       uuid.UUID `json:"question_id"`
	Answer           string    `json:"answer"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
