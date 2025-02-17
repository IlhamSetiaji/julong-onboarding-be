package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmployeeTaskChecklist struct {
	gorm.Model     `json:"-"`
	ID             uuid.UUID  `json:"id" gorm:"type:char(36);primaryKey;"`
	EmployeeTaskID uuid.UUID  `json:"employee_task_id" gorm:"type:char(36);not null"`
	Name           string     `json:"name" gorm:"type:varchar(255);not null"`
	IsChecked      string     `json:"is_checked" gorm:"type:varchar(255);not null;default:'NO'"`
	VerifiedBy     *uuid.UUID `json:"verified_by" gorm:"type:char(36);default:null"`

	EmployeeTask *EmployeeTask `json:"employee_task" gorm:"foreignKey:EmployeeTaskID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (e *EmployeeTaskChecklist) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	e.CreatedAt = time.Now().In(loc)
	e.UpdatedAt = time.Now().In(loc)
	return nil
}

func (e *EmployeeTaskChecklist) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	e.UpdatedAt = time.Now().In(loc)
	return nil
}

func (EmployeeTaskChecklist) TableName() string {
	return "employee_task_checklists"
}
