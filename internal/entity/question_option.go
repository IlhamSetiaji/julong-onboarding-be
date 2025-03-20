package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionOption struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	QuestionID uuid.UUID `json:"question_id" gorm:"type:char(36);not null"`
	OptionText string    `json:"option_text" gorm:"type:text;not null"`

	Question *Question `json:"question" gorm:"foreignKey:QuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (q *QuestionOption) BeforeCreate(tx *gorm.DB) (err error) {
	q.ID = uuid.New()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	q.CreatedAt = time.Now().In(loc)
	q.UpdatedAt = time.Now().In(loc)
	return
}

func (q *QuestionOption) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	q.UpdatedAt = time.Now().In(loc)
	return
}

func (q *QuestionOption) BeforeDelete(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	q.DeletedAt = gorm.DeletedAt{
		Time:  time.Now().In(loc),
		Valid: true,
	}
	return
}

func (QuestionOption) TableName() string {
	return "question_options"
}
