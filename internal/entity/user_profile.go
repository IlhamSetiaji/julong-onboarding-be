package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MaritalStatusEnum string

const (
	MARITAL_STATUS_ENUM_SINGLE   MaritalStatusEnum = "single"
	MARITAL_STATUS_ENUM_MARRIED  MaritalStatusEnum = "married"
	MARITAL_STATUS_ENUM_DIVORCED MaritalStatusEnum = "divorced"
	MARITAL_STATUS_ENUM_WIDOWED  MaritalStatusEnum = "widowed"
	MARITAL_STATUS_ENUM_ANY      MaritalStatusEnum = "any"
)

type UserStatus string
type UserGender string

const (
	USER_ACTIVE   UserStatus = "ACTIVE"
	USER_INACTIVE UserStatus = "INACTIVE"
	USER_PENDING  UserStatus = "PENDING"
)

const (
	MALE   UserGender = "MALE"
	FEMALE UserGender = "FEMALE"
)

type UserProfile struct {
	gorm.Model      `json:"-"`
	ID              uuid.UUID         `json:"id" gorm:"type:char(36);primaryKey;"`
	UserID          *uuid.UUID        `json:"user_id" gorm:"type:char(36);not null"`
	Name            string            `json:"name" gorm:"type:varchar(255);default:null"`
	MaritalStatus   MaritalStatusEnum `json:"marital_status" gorm:"type:varchar(255);default:'single'"`
	Gender          UserGender        `json:"gender" gorm:"type:varchar(255);default:null"`
	Status          UserStatus        `json:"status" gorm:"type:varchar(255);default:null"`
	PhoneNumber     string            `json:"phone_number" gorm:"type:varchar(255);default:null"`
	Age             int               `json:"age" gorm:"type:int;default:null"`
	BirthDate       time.Time         `json:"birth_date" gorm:"type:date;default:null"`
	BirthPlace      string            `json:"birth_place" gorm:"type:varchar(255);default:null"`
	Address         string            `json:"address" gorm:"type:text;default:null"`
	Ktp             string            `json:"ktp" gorm:"type:varchar(255);default:null"`
	CurriculumVitae string            `json:"curriculum_vitae" gorm:"type:text;default:null"`
	Avatar          string            `json:"avatar" gorm:"type:text;default:null"`
	Bilingual       string            `json:"bilingual" gorm:"type:varchar(255);default:null"`

	WorkExperiences []WorkExperience `json:"work_experiences" gorm:"foreignKey:UserProfileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Educations      []Education      `json:"educations" gorm:"foreignKey:UserProfileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Skills          []Skill          `json:"skills" gorm:"foreignKey:UserProfileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (up *UserProfile) BeforeCreate(tx *gorm.DB) (err error) {
	up.ID = uuid.New()
	up.CreatedAt = time.Now()
	up.UpdatedAt = time.Now()
	return
}

func (up *UserProfile) BeforeUpdate(tx *gorm.DB) (err error) {
	up.UpdatedAt = time.Now()
	return
}

func (UserProfile) TableName() string {
	return "user_profiles"
}
