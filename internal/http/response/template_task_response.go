package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
)

type TemplateTaskResponse struct {
	ID          uuid.UUID                       `json:"id"`
	CoverPath   *string                         `json:"cover_path`
	Name        string                          `json:"name"`
	Priority    entity.TemplateTaskPriorityEnum `json:"priority"`
	DueDuration *int                            `json:"due_duration"`
	Status      entity.TemplateTaskStatusEnum   `json:"status"`
	Description string                          `json:"description"`
	Source      string                          `json:"source"`
	CreatedAt   time.Time                       `json:"created_at"`
	UpdatedAt   time.Time                       `json:"updated_at"`

	TemplateTaskAttachments []TemplateTaskAttachmentResponse `json:"template_task_attachments"`
	TemplateTaskChecklists  []TemplateTaskChecklistResponse  `json:"template_task_checklists"`
}
