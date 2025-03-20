package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Question struct {
	gorm.Model       `json:"-"`
	ID               uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	SurveyTemplateID uuid.UUID `json:"survey_template_id" gorm:"type:char(36);not null"`
	AnswerTypeID     uuid.UUID `json:"answer_type_id" gorm:"type:char(36);not null"`
	Question         string    `json:"question" gorm:"type:text;not null"`
	Attachment       *string   `json:"attachment" gorm:"type:varchar(255);default:null"`
	IsCompleted      string    `json:"is_completed" gorm:"type:varchar(255);not null;default:'NO'"`

	SurveyTemplate  *SurveyTemplate  `json:"survey_template" gorm:"foreignKey:SurveyTemplateID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AnswerType      *AnswerType      `json:"answer_type" gorm:"foreignKey:AnswerTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	QuestionOptions []QuestionOption `json:"question_options" gorm:"foreignKey:QuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SurveyResponses []SurveyResponse `json:"survey_responses" gorm:"foreignKey:QuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (q *Question) BeforeCreate(tx *gorm.DB) (err error) {
	q.ID = uuid.New()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	q.CreatedAt = time.Now().In(loc)
	q.UpdatedAt = time.Now().In(loc)
	return
}

func (q *Question) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	q.UpdatedAt = time.Now().In(loc)
	return
}

func (q *Question) BeforeDelete(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	q.DeletedAt = gorm.DeletedAt{
		Time:  time.Now().In(loc),
		Valid: true,
	}
	return
}

func (Question) TableName() string {
	return "questions"
}
