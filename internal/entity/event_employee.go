package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventEmployee struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID  `json:"id" gorm:"type:char(36);primaryKey;"`
	EventID    uuid.UUID  `json:"event_id" gorm:"type:char(36);not null"`
	EmployeeID *uuid.UUID `json:"employee_id" gorm:"type:char(36);not null"`

	Event *Event `json:"event" gorm:"foreignKey:EventID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (e *EventEmployee) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	e.CreatedAt = time.Now().In(loc)
	e.UpdatedAt = time.Now().In(loc)
	return nil
}

func (e *EventEmployee) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	e.UpdatedAt = time.Now().In(loc)
	return nil
}

func (EventEmployee) TableName() string {
	return "event_employees"
}
