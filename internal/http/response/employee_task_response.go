package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
)

type EmployeeTaskResponse struct {
	ID             uuid.UUID                       `json:"id"`
	CoverPath      *string                         `json:"cover_path"`
	EmployeeID     *uuid.UUID                      `json:"employee_id"`
	TemplateTaskID *uuid.UUID                      `json:"template_task_id"`
	VerifiedBy     *uuid.UUID                      `json:"verified_by"`
	Name           string                          `json:"name"`
	Priority       entity.EmployeeTaskPriorityEnum `json:"priority"`
	Description    string                          `json:"description"`
	StartDate      time.Time                       `json:"start_date"`
	EndDate        time.Time                       `json:"end_date"`
	IsDone         string                          `json:"is_done"`
	Proof          *string                         `json:"proof"`
	Status         entity.EmployeeTaskStatusEnum   `json:"status"`
	Kanban         entity.EmployeeTaskKanbanEnum   `json:"kanban"`
	Notes          string                          `json:"notes"`
	Source         string                          `json:"source"`
	CreatedAt      time.Time                       `json:"created_at"`
	UpdatedAt      time.Time                       `json:"updated_at"`

	VerifiedByName string `json:"verified_by_name"`
	EmployeeName   string `json:"employee_name"`

	TemplateTask            *TemplateTaskResponse            `json:"template_task"`
	EmployeeTaskAttachments []EmployeeTaskAttachmentResponse `json:"employee_task_attachments"`
	EmployeeTaskChecklists  []EmployeeTaskChecklistResponse  `json:"employee_task_checklists"`
}

type EmployeeTaskKanbanResponse struct {
	ToDo       []EmployeeTaskResponse `json:"to_do"`
	InProgress []EmployeeTaskResponse `json:"in_progress"`
	NeedReview []EmployeeTaskResponse `json:"need_review"`
	Completed  []EmployeeTaskResponse `json:"completed"`
}
