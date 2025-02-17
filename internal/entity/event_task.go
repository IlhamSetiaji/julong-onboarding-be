package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventTask struct {
	gorm.Model     `json:"-"`
	ID             uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	EventID        uuid.UUID `json:"event_id" gorm:"type:char(36);not null"`
	TemplateTaskID uuid.UUID `json:"template_task_id" gorm:"type:char(36);not null"`
	Name           string    `json:"name" gorm:"type:varchar(255);not null"`

	Event        *Event        `json:"event" gorm:"foreignKey:EventID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TemplateTask *TemplateTask `json:"template_task" gorm:"foreignKey:TemplateTaskID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (e *EventTask) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	e.CreatedAt = time.Now().In(loc)
	e.UpdatedAt = time.Now().In(loc)
	return nil
}

func (e *EventTask) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	e.UpdatedAt = time.Now().In(loc)
	return nil
}

func (EventTask) TableName() string {
	return "event_tasks"
}
