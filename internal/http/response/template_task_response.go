package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
)

type TemplateTaskResponse struct {
	ID               uuid.UUID                       `json:"id"`
	CoverPath        *string                         `json:"cover_path"`
	CoverPathOrigin  *string                         `json:"cover_path_origin"`
	Name             string                          `json:"name"`
	SurveyTemplateID *uuid.UUID                      `json:"survey_template_id"`
	Priority         entity.TemplateTaskPriorityEnum `json:"priority"`
	DueDuration      *int                            `json:"due_duration"`
	Status           entity.TemplateTaskStatusEnum   `json:"status"`
	Description      string                          `json:"description"`
	Source           string                          `json:"source"`
	OrganizationType string                          `json:"organization_type"`
	CreatedAt        time.Time                       `json:"created_at"`
	UpdatedAt        time.Time                       `json:"updated_at"`

	TemplateTaskAttachments []TemplateTaskAttachmentResponse `json:"template_task_attachments"`
	TemplateTaskChecklists  []TemplateTaskChecklistResponse  `json:"template_task_checklists"`
	SurveyTemplate          *SurveyTemplateResponse          `json:"survey_template"`
}
