package request

import "mime/multipart"

type CreateEmployeeTaskRequest struct {
	CoverPath               *string                         `form:"cover_path" validate:"required"`
	EmployeeID              *string                         `form:"employee_id" validate:"required,uuid"`
	TemplateTaskID          *string                         `form:"template_task_id" validate:"omitempty,uuid"`
	Name                    string                          `form:"name" validate:"required"`
	Priority                string                          `form:"priority" validate:"required,employee_task_priority_validation"`
	Description             string                          `form:"description" validate:"omitempty"`
	StartDate               string                          `form:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate                 string                          `form:"end_date" validate:"required,datetime=2006-01-02"`
	EmployeeTaskAttachments []EmployeeTaskAttachmentRequest `form:"employee_task_attachments" validate:"omitempty,dive"`
	EmployeeTaskChecklists  []EmployeeTaskChecklistRequest  `form:"employee_task_checklists" validate:"omitempty,dive"`
}

type UpdateEmployeeTaskRequest struct {
	ID                      *string                         `form:"id" validate:"required"`
	CoverPath               *string                         `form:"cover_path" validate:"required"`
	EmployeeID              *string                         `form:"employee_id" validate:"required,uuid"`
	TemplateTaskID          *string                         `form:"template_task_id" validate:"omitempty,uuid"`
	Name                    string                          `form:"name" validate:"required"`
	Priority                string                          `form:"priority" validate:"required,employee_task_priority_validation"`
	Description             string                          `form:"description" validate:"omitempty"`
	StartDate               string                          `form:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate                 string                          `form:"end_date" validate:"required,datetime=2006-01-02"`
	VerifiedBy              *string                         `form:"verified_by" validate:"omitempty,uuid"`
	IsDone                  string                          `form:"is_done" validate:"omitempty,oneof=YES NO"`
	Proof                   *multipart.FileHeader           `form:"proof" validate:"omitempty"`
	ProofPath               *string                         `form:"proof_path" validate:"omitempty"`
	Status                  string                          `form:"status" validate:"omitempty,employee_task_status_validation"`
	Kanban                  string                          `form:"kanban" validate:"omitempty,employee_task_kanban_validation"`
	Notes                   string                          `form:"notes" validate:"omitempty"`
	EmployeeTaskAttachments []EmployeeTaskAttachmentRequest `form:"employee_task_attachments" validate:"omitempty,dive"`
	EmployeeTaskChecklists  []EmployeeTaskChecklistRequest  `form:"employee_task_checklists" validate:"omitempty,dive"`
}

type UpdateEmployeeTaskOnlyRequest struct {
	ID        *string `form:"id" validate:"required"`
	StartDate string  `form:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string  `form:"end_date" validate:"required,datetime=2006-01-02"`
}
