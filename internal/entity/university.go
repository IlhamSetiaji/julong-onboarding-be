package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type University struct {
	gorm.Model   `json:"-"`
	ID           uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	Name         string    `json:"name" gorm:"type:varchar(255);default:null"`
	Country      string    `json:"country" gorm:"type:varchar(255);default:null"`
	AlphaTwoCode string    `json:"alpha_two_code" gorm:"type:char(2);default:null"`
}

func (u *University) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *University) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return nil
}

func (University) TableName() string {
	return "universities"
}
