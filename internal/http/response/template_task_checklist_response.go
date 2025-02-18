package response

import (
	"time"

	"github.com/google/uuid"
)

type TemplateTaskChecklistResponse struct {
	ID             uuid.UUID `json:"id"`
	TemplateTaskID uuid.UUID `json:"template_task_id"`
	Name           string    `json:"name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
