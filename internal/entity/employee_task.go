package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmployeeTaskPriorityEnum string

const (
	EMPLOYEE_TASK_PRIORITY_ENUM_LOW    EmployeeTaskPriorityEnum = "LOW"
	EMPLOYEE_TASK_PRIORITY_ENUM_MEDIUM EmployeeTaskPriorityEnum = "MEDIUM"
	EMPLOYEE_TASK_PRIORITY_ENUM_HIGH   EmployeeTaskPriorityEnum = "HIGH"
)

type EmployeeTaskStatusEnum string

const (
	EMPLOYEE_TASK_STATUS_ENUM_ACTIVE   EmployeeTaskStatusEnum = "ACTIVE"
	EMPLOYEE_TASK_STATUS_ENUM_INACTIVE EmployeeTaskStatusEnum = "INACTIVE"
)

type EmployeeTaskKanbanEnum string

const (
	EMPLOYEE_TASK_KANBAN_ENUM_TODO        EmployeeTaskKanbanEnum = "TO_DO"
	EPMLOYEE_TASK_KANBAN_ENUM_IN_PROGRESS EmployeeTaskKanbanEnum = "IN_PROGRESS"
	EMPLOYEE_TASK_KANBAN_ENUM_NEED_REVIEW EmployeeTaskKanbanEnum = "NEED_REVIEW"
	EMPLOYEE_TASK_KANBAN_ENUM_COMPLETED   EmployeeTaskKanbanEnum = "COMPLETED"
)

type EmployeeTask struct {
	gorm.Model       `json:"-"`
	ID               uuid.UUID                `json:"id" gorm:"type:char(36);primaryKey;"`
	CoverPath        *string                  `json:"cover_path" gorm:"type:varchar(255);default:null"`
	EmployeeID       *uuid.UUID               `json:"employee_id" gorm:"type:char(36);not null"`
	TemplateTaskID   *uuid.UUID               `json:"template_task_id" gorm:"type:char(36);default:null"`
	SurveyTemplateID *uuid.UUID               `json:"survey_template_id" gorm:"type:char(36);default:null"`
	VerifiedBy       *uuid.UUID               `json:"verified_by" gorm:"type:char(36);default:null"`
	Name             string                   `json:"name" gorm:"type:varchar(255);not null"`
	Priority         EmployeeTaskPriorityEnum `json:"priority" gorm:"type:varchar(255);not null"`
	Description      string                   `json:"description" gorm:"type:text;default:null"`
	StartDate        time.Time                `json:"start_date" gorm:"type:date;not null"`
	EndDate          time.Time                `json:"end_date" gorm:"type:date;not null"`
	IsDone           string                   `json:"is_done" gorm:"type:varchar(255);not null;default:'NO'"`
	Proof            *string                  `json:"proof" gorm:"type:varchar(255);default:null"`
	Status           EmployeeTaskStatusEnum   `json:"status" gorm:"type:varchar(255);not null;default:'ACTIVE'"`
	Kanban           EmployeeTaskKanbanEnum   `json:"kanban" gorm:"type:varchar(255);not null;default:'TO_DO'"`
	Notes            string                   `json:"notes" gorm:"type:text;default:null"`
	Source           string                   `json:"source" gorm:"type:varchar(255);default:null"`
	MidsuitID        *string                  `json:"midsuit_id" gorm:"type:varchar(255);default:null"`

	TemplateTask            *TemplateTask            `json:"template_task" gorm:"foreignKey:TemplateTaskID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	EmployeeTaskChecklists  []EmployeeTaskChecklist  `json:"employee_task_checklists" gorm:"foreignKey:EmployeeTaskID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	EmployeeTaskAttachments []EmployeeTaskAttachment `json:"employee_task_attachments" gorm:"foreignKey:EmployeeTaskID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	EmployeeTaskFiles       []EmployeeTaskFiles      `json:"employee_task_files" gorm:"foreignKey:EmployeeTaskID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SurveyResponses         []SurveyResponse         `json:"survey_responses" gorm:"foreignKey:EmployeeTaskID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SurveyTemplate          *SurveyTemplate          `json:"survey_template" gorm:"foreignKey:SurveyTemplateID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (e *EmployeeTask) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	e.CreatedAt = time.Now().In(loc)
	e.UpdatedAt = time.Now().In(loc)
	return nil
}

func (e *EmployeeTask) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	e.UpdatedAt = time.Now().In(loc)
	return nil
}

func (EmployeeTask) TableName() string {
	return "employee_tasks"
}
