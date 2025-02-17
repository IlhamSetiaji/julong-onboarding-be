package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkExperience struct {
	gorm.Model     `json:"-"`
	ID             uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	UserProfileID  uuid.UUID `json:"user_profile_id" gorm:"type:char(36);not null"`
	Name           string    `json:"name" gorm:"type:varchar(255);not null"`
	CompanyName    string    `json:"company_name" gorm:"type:varchar(255);not null"`
	YearExperience int       `json:"year_experienced" gorm:"type:int;not null"`
	JobDescription string    `json:"job_description" gorm:"type:text;not null"`
	Certificate    string    `json:"certificate" gorm:"type:text;not null"`

	// Applicant *Applicant `json:"applicant" gorm:"foreignKey:ApplicantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserProfile *UserProfile `json:"user_profile" gorm:"foreignKey:UserProfileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (we *WorkExperience) BeforeCreate(tx *gorm.DB) (err error) {
	we.ID = uuid.New()
	we.CreatedAt = time.Now()
	we.UpdatedAt = time.Now()
	return
}

func (we *WorkExperience) BeforeUpdate(tx *gorm.DB) (err error) {
	we.UpdatedAt = time.Now()
	return
}

func (WorkExperience) TableName() string {
	return "work_experiences"
}
