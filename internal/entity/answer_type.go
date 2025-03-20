package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnswerType struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	Name       string    `json:"name" gorm:"type:varchar(255);not null"`

	Questions []Question `json:"questions" gorm:"foreignKey:AnswerTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (a *AnswerType) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	a.CreatedAt = time.Now().In(loc)
	a.UpdatedAt = time.Now().In(loc)
	return
}

func (a *AnswerType) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	a.UpdatedAt = time.Now().In(loc)
	return
}

func (a *AnswerType) BeforeDelete(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	a.DeletedAt = gorm.DeletedAt{
		Time:  time.Now().In(loc),
		Valid: true,
	}
	return
}

func (AnswerType) TableName() string {
	return "answer_types"
}
