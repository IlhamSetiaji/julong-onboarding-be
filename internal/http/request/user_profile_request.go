package request

import "mime/multipart"

type FillUserProfileRequest struct {
	ID              string                `form:"id" validate:"omitempty,uuid"`
	Name            string                `form:"name" validate:"omitempty"`
	MaritalStatus   string                `form:"marital_status" validate:"required,marital_status_validation"`
	Gender          string                `form:"gender" validate:"required,user_gender_validation"`
	PhoneNumber     string                `form:"phone_number" validate:"required"`
	Age             int                   `form:"age" validate:"required"`
	BirthDate       string                `form:"birth_date" validate:"required,datetime=2006-01-02"`
	BirthPlace      string                `form:"birth_place" validate:"required"`
	Address         string                `form:"address" validate:"omitempty"`
	Bilingual       string                `form:"bilingual" validate:"omitempty"`
	Ktp             *multipart.FileHeader `form:"ktp" validate:"omitempty"`
	CurriculumVitae *multipart.FileHeader `form:"curriculum_vitae" validate:"omitempty"`
	KtpPath         string                `form:"ktp_path" validate:"omitempty"`
	CvPath          string                `form:"cv_path" validate:"omitempty"`
	WorkExperiences []WorkExperience      `form:"work_experiences" validate:"omitempty,dive"`
	Educations      []Education           `form:"educations" validate:"omitempty,dive"`
	Skills          []Skill               `form:"skills" validate:"omitempty,dive"`
}

type UpdateStatusUserProfileRequest struct {
	ID     string `json:"id" validate:"required,uuid"`
	Status string `json:"status" validate:"required,user_status_validation"`
}
