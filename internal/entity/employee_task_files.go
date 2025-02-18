package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmployeeTaskFiles struct {
	gorm.Model     `json:"-"`
	ID             uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	EmployeeTaskID uuid.UUID `json:"employee_task_id" gorm:"type:char(36);not null"`
	Path           string    `json:"path" gorm:"type:varchar(255);not null"`

	EmployeeTask *EmployeeTask `json:"employee_task" gorm:"foreignKey:EmployeeTaskID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (e *EmployeeTaskFiles) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	e.CreatedAt = time.Now().In(loc)
	e.UpdatedAt = time.Now().In(loc)
	return nil
}

func (e *EmployeeTaskFiles) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	e.UpdatedAt = time.Now().In(loc)
	return nil
}

func (EmployeeTaskFiles) TableName() string {
	return "employee_task_files"
}
