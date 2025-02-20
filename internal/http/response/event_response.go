package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
)

type EventResponse struct {
	ID             uuid.UUID              `json:"id"`
	TemplateTaskID uuid.UUID              `json:"template_task_id"`
	Name           string                 `json:"name"`
	StartDate      time.Time              `json:"start_date"`
	EndDate        time.Time              `json:"end_date"`
	Description    string                 `json:"description"`
	Status         entity.EventStatusEnum `json:"status"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`

	TemplateTask   *TemplateTaskResponse   `json:"template_task"`
	EventEmployees []EventEmployeeResponse `json:"event_employees"`
}
