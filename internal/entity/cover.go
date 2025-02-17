package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cover struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	Path       string    `json:"path" gorm:"type:varchar(255);not null"`
}

func (c *Cover) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	c.CreatedAt = time.Now().In(loc)
	c.UpdatedAt = time.Now().In(loc)
	return nil
}

func (c *Cover) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	c.UpdatedAt = time.Now().In(loc)
	return nil
}

func (Cover) TableName() string {
	return "covers"
}
