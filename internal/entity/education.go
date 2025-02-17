package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EducationLevelEnum string

const (
	EDUCATION_LEVEL_ENUM_DOCTORAL EducationLevelEnum = "1 - Doctoral / Professor"
	EDUCATION_LEVEL_ENUM_MASTER   EducationLevelEnum = "2 - Master Degree"
	EDUCATION_LEVEL_ENUM_BACHELOR EducationLevelEnum = "3 - Bachelor"
	EDUCATION_LEVEL_ENUM_D1       EducationLevelEnum = "4 - Diploma 1"
	EDUCATION_LEVEL_ENUM_D2       EducationLevelEnum = "5 - Diploma 2"
	EDUCATION_LEVEL_ENUM_D3       EducationLevelEnum = "6 - Diploma 3"
	EDUCATION_LEVEL_ENUM_D4       EducationLevelEnum = "7 - Diploma 4"
	EDUCATION_LEVEL_ENUM_SD       EducationLevelEnum = "8 - Elementary School"
	EDUCATION_LEVEL_ENUM_SMA      EducationLevelEnum = "9 - Senior High School"
	EDUCATION_LEVEL_ENUM_SMP      EducationLevelEnum = "10 - Junior High School"
	EDUCATION_LEVEL_ENUM_TK       EducationLevelEnum = "11 - Unschooled"
)

type Education struct {
	gorm.Model     `json:"-"`
	ID             uuid.UUID          `json:"id" gorm:"type:char(36);primaryKey;"`
	UserProfileID  uuid.UUID          `json:"user_profile_id" gorm:"type:char(36);not null"`
	EducationLevel EducationLevelEnum `json:"education_level" gorm:"type:varchar(255);not null"`
	Major          string             `json:"major" gorm:"type:varchar(255);not null"`
	SchoolName     string             `json:"school_name" gorm:"type:varchar(255);not null"`
	GraduateYear   int                `json:"graduate_year" gorm:"type:int;not null"`
	EndDate        time.Time          `json:"end_date" gorm:"type:date;not null"`
	Certificate    string             `json:"certificate" gorm:"type:text;not null"`
	Gpa            *float64           `json:"gpa" gorm:"type:float;not null"`

	// Applicant *Applicant `json:"applicant" gorm:"foreignKey:ApplicantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserProfile *UserProfile `json:"user_profile" gorm:"foreignKey:UserProfileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (e *Education) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New()
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Education) BeforeUpdate(tx *gorm.DB) (err error) {
	e.UpdatedAt = time.Now()
	return nil
}

func (Education) TableName() string {
	return "educations"
}
