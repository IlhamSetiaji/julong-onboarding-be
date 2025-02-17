package request

import "mime/multipart"

type Education struct {
	ID              *string               `form:"id" validate:"omitempty,uuid"`
	EducationLevel  string                `form:"education_level" validate:"required,education_level_validation"`
	Major           string                `form:"major" validate:"required"`
	SchoolName      string                `form:"school_name" validate:"required"`
	GraduateYear    int                   `form:"graduate_year" validate:"required"`
	EndDate         string                `form:"end_date" validate:"required,datetime=2006-01-02"`
	Certificate     *multipart.FileHeader `form:"certificate" validate:"omitempty"`
	CertificatePath string                `form:"certificate_path" validate:"omitempty"`
	Gpa             *float64              `form:"gpa" validate:"omitempty,gte=0"`
}
