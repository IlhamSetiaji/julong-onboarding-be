package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SurveyResponse struct {
	gorm.Model       `json:"-"`
	ID               uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	SurveyTemplateID uuid.UUID `json:"survey_template_id" gorm:"type:char(36);not null"`
	EmployeeTaskID   uuid.UUID `json:"employee_task_id" gorm:"type:char(36);not null"`
	QuestionID       uuid.UUID `json:"question_id" gorm:"type:char(36);not null"`
	Answer           string    `json:"answer" gorm:"type:text;default:null"`
	AnswerFile       string    `json:"answer_file" gorm:"type:text;default:null"`

	SurveyTemplate *SurveyTemplate `json:"survey_template" gorm:"foreignKey:SurveyTemplateID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	EmployeeTask   *EmployeeTask   `json:"employee_task" gorm:"foreignKey:EmployeeTaskID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Question       *Question       `json:"question" gorm:"foreignKey:QuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (s *SurveyResponse) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	s.CreatedAt = time.Now().In(loc)
	s.UpdatedAt = time.Now().In(loc)
	return
}

func (s *SurveyResponse) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	s.UpdatedAt = time.Now().In(loc)
	return
}

func (s *SurveyResponse) BeforeDelete(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	s.DeletedAt = gorm.DeletedAt{
		Time:  time.Now().In(loc),
		Valid: true,
	}
	return
}

func (SurveyResponse) TableName() string {
	return "survey_responses"
}
