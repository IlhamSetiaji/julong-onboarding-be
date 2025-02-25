package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmployeeHiring struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	EmployeeID uuid.UUID `json:"employee_id" gorm:"type:char(36);not null"`
	HiringDate time.Time `json:"hiring_date" gorm:"type:date;not null"`
}

func (e *EmployeeHiring) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	e.CreatedAt = time.Now().In(loc)
	e.UpdatedAt = time.Now().In(loc)
	return nil
}

func (e *EmployeeHiring) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	e.UpdatedAt = time.Now().In(loc)
	return nil
}
