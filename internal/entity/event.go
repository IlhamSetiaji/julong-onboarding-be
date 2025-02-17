package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventStatusEnum string

const (
	EVENT_STATUS_ENUM_UPCOMING EventStatusEnum = "UPCOMING"
	EVENT_STATUS_ENUM_ONGOING  EventStatusEnum = "ONGOING"
	EVENT_STATUS_ENUM_FINISHED EventStatusEnum = "FINISHED"
)

type Event struct {
	gorm.Model  `json:"-"`
	ID          uuid.UUID       `json:"id" gorm:"type:char(36);primaryKey;"`
	Name        string          `json:"name" gorm:"type:varchar(255);not null"`
	StartDate   time.Time       `json:"start_date" gorm:"type:date;not null"`
	EndDate     time.Time       `json:"end_date" gorm:"type:date;not null"`
	Description string          `json:"description" gorm:"type:text;default:null"`
	Status      EventStatusEnum `json:"status" gorm:"type:varchar(255);not null;default:'UPCOMING'"`

	EventTasks     []EventTask     `json:"event_tasks" gorm:"foreignKey:EventID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	EventEmployees []EventEmployee `json:"event_employees" gorm:"foreignKey:EventID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	e.CreatedAt = time.Now().In(loc)
	e.UpdatedAt = time.Now().In(loc)
	return nil
}

func (e *Event) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	e.UpdatedAt = time.Now().In(loc)
	return nil
}

func (Event) TableName() string {
	return "events"
}
