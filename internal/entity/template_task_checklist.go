package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TemplateTaskChecklist struct {
	gorm.Model     `json:"-"`
	ID             uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	TemplateTaskID uuid.UUID `json:"template_task_id" gorm:"type:char(36);not null"`
	Name           string    `json:"name" gorm:"type:varchar(255);not null"`

	TemplateTask *TemplateTask `json:"template_task" gorm:"foreignKey:TemplateTaskID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (t *TemplateTaskChecklist) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	t.CreatedAt = time.Now().In(loc)
	t.UpdatedAt = time.Now().In(loc)
	return nil
}

func (t *TemplateTaskChecklist) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	t.UpdatedAt = time.Now().In(loc)
	return nil
}

func (TemplateTaskChecklist) TableName() string {
	return "template_task_checklists"
}
