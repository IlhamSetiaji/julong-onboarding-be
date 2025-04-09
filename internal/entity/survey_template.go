package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SurveyTemplateStatusEnum string

const (
	SURVEY_TEMPLATE_STATUS_ENUM_DRAFT     SurveyTemplateStatusEnum = "DRAFT"
	SURVEY_TEMPLATE_STATUS_ENUM_SUBMITTED SurveyTemplateStatusEnum = "SUBMITTED"
)

type SurveyTemplate struct {
	gorm.Model   `json:"-"`
	ID           uuid.UUID                `json:"id" gorm:"type:char(36);primaryKey;"`
	SurveyNumber string                   `json:"survey_number" gorm:"type:varchar(255);not null"`
	Title        string                   `json:"title" gorm:"type:varchar(255);not null"`
	Status       SurveyTemplateStatusEnum `json:"status" gorm:"type:varchar(255);not null;default:'DRAFT'"`

	Questions       []Question       `json:"questions" gorm:"foreignKey:SurveyTemplateID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SurveyResponses []SurveyResponse `json:"survey_responses" gorm:"foreignKey:SurveyTemplateID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	EmployeeTasks   []EmployeeTask   `json:"employee_tasks" gorm:"foreignKey:SurveyTemplateID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (s *SurveyTemplate) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	s.CreatedAt = time.Now().In(loc)
	s.UpdatedAt = time.Now().In(loc)
	return
}

func (s *SurveyTemplate) BeforeUpdate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	s.UpdatedAt = time.Now().In(loc)
	return
}

func (s *SurveyTemplate) BeforeDelete(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	s.DeletedAt = gorm.DeletedAt{
		Time:  time.Now().In(loc),
		Valid: true,
	}

	if err := tx.Model(&TemplateTask{}).Where("survey_template_id = ?", s.ID).Update("survey_template_id", nil).Error; err != nil {
		return err
	}

	if err := tx.Model(&EmployeeTask{}).Where("survey_template_id = ?", s.ID).Update("survey_template_id", nil).Error; err != nil {
		return err
	}

	if err := tx.Model(&SurveyResponse{}).Where("survey_template_id = ?", s.ID).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}

	if err := tx.Model(&Question{}).Where("survey_template_id = ?", s.ID).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}

	return
}

func (SurveyTemplate) TableName() string {
	return "survey_templates"
}
