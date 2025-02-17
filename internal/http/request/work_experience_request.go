package request

import "mime/multipart"

type WorkExperience struct {
	ID              *string               `form:"id" validate:"omitempty"`
	Name            string                `form:"name" validate:"required"`
	CompanyName     string                `form:"company_name" validate:"required"`
	YearExperience  int                   `form:"year_experience" validate:"required"`
	JobDescription  string                `form:"job_description" validate:"required"`
	Certificate     *multipart.FileHeader `form:"certificate" validate:"omitempty"`
	CertificatePath string                `form:"certificate_path" validate:"omitempty"`
}
