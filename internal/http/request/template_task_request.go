package request

type CreateTemplateTaskRequest struct {
	// CoverFile               *multipart.FileHeader           `form:"cover_file" validate:"required"`
	CoverPath               string                          `form:"cover_path" validate:"required"`
	Name                    string                          `form:"name" validate:"required"`
	Priority                string                          `form:"priority" validate:"required,template_task_priority_validation"`
	DueDuration             *int                            `form:"due_duration" validate:"omitempty,numeric"`
	Status                  string                          `form:"status" validate:"required,template_task_status_validation"`
	Description             string                          `form:"description" validate:"omitempty"`
	TemplateTaskAttachments []TemplateTaskAttachmentRequest `form:"template_task_attachments" validate:"omitempty,dive"`
	TemplateTaskChecklists  []TemplateTaskChecklistRequest  `form:"template_task_checklists" validate:"omitempty,dive"`
}

type UpdateTemplateTaskRequest struct {
	ID string `form:"id" validate:"required"`
	// CoverFile               *multipart.FileHeader           `form:"cover_file" validate:"required"`
	CoverPath               string                          `form:"cover_path" validate:"required"`
	Name                    string                          `form:"name" validate:"required"`
	Priority                string                          `form:"priority" validate:"required,template_task_priority_validation"`
	DueDuration             *int                            `form:"due_duration" validate:"omitempty,numeric"`
	Status                  string                          `form:"status" validate:"required,template_task_status_validation"`
	Description             string                          `form:"description" validate:"omitempty"`
	TemplateTaskAttachments []TemplateTaskAttachmentRequest `form:"template_task_attachments" validate:"omitempty,dive"`
	TemplateTaskChecklists  []TemplateTaskChecklistRequest  `form:"template_task_checklists" validate:"omitempty,dive"`
}
