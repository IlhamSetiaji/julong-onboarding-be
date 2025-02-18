package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TemplateTaskPriorityEnum string

const (
	TEMPLATE_TASK_PRIORITY_ENUM_LOW    TemplateTaskPriorityEnum = "LOW"
	TEMPLATE_TASK_PRIORITY_ENUM_MEDIUM TemplateTaskPriorityEnum = "MEDIUM"
	TEMPLATE_TASK_PRIORITY_ENUM_HIGH   TemplateTaskPriorityEnum = "HIGH"
)

type TemplateTaskStatusEnum string

const (
	TEMPLATE_TASK_STATUS_ENUM_ACTIVE   TemplateTaskStatusEnum = "ACTIVE"
	TEMPLATE_TASK_STATUS_ENUM_INACTIVE TemplateTaskStatusEnum = "INACTIVE"
)

type TemplateTask struct {
	gorm.Model  `json:"-"`
	ID          uuid.UUID                `json:"id" gorm:"type:char(36);primaryKey;"`
	CoverPath   *string                  `json:"cover_path" gorm:"type:varchar(255);default:null"`
	Name        string                   `json:"name" gorm:"type:varchar(255);not null"`
	Priority    TemplateTaskPriorityEnum `json:"priority" gorm:"type:varchar(255);not null"`
	DueDuration *int                     `json:"due_duration" gorm:"type:int;default:0"`
	Status      TemplateTaskStatusEnum   `json:"status" gorm:"type:varchar(255);not null"`
	Description string                   `json:"description" gorm:"type:text;default:null"`
	Source      string                   `json:"source" gorm:"type:text;default:null"`

	TemplateTaskAttachments []TemplateTaskAttachment `json:"template_task_attachments" gorm:"foreignKey:TemplateTaskID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TemplateTaskChecklists  []TemplateTaskChecklist  `json:"template_task_checklists" gorm:"foreignKey:TemplateTaskID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Events                  []Event                  `json:"events" gorm:"foreignKey:TemplateTaskID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (t *TemplateTask) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	t.CreatedAt = time.Now().In(loc)
	t.UpdatedAt = time.Now().In(loc)
	return nil
}

func (t *TemplateTask) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	t.UpdatedAt = time.Now().In(loc)
	return nil
}

func (TemplateTask) TableName() string {
	return "template_tasks"
}
