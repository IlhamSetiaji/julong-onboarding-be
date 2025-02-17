package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
)

type UserProfileResponse struct {
	ID              uuid.UUID                 `json:"id"`
	UserID          *uuid.UUID                `json:"user_id"`
	Name            string                    `json:"name"`
	MaritalStatus   entity.MaritalStatusEnum  `json:"marital_status"`
	Gender          entity.UserGender         `json:"gender"`
	PhoneNumber     string                    `json:"phone_number"`
	Age             int                       `json:"age"`
	Address         string                    `json:"address"`
	BirthDate       time.Time                 `json:"birth_date"`
	BirthPlace      string                    `json:"birth_place"`
	Bilingual       string                    `json:"bilingual"`
	Ktp             *string                   `json:"ktp"`
	Avatar          *string                   `json:"avatar"`
	CurriculumVitae *string                   `json:"curriculum_vitae"`
	Status          entity.UserStatus         `json:"status"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
	WorkExperiences *[]WorkExperienceResponse `json:"work_experiences"`
	Educations      *[]EducationResponse      `json:"educations"`
	Skills          *[]SkillResponse          `json:"skills"`
	User            *map[string]interface{}   `json:"user"`
}
