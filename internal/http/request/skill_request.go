package request

import "mime/multipart"

type Skill struct {
	ID              *string               `form:"id" validate:"omitempty"`
	Name            string                `form:"name" validate:"required"`
	Description     string                `form:"description" validate:"required"`
	Level           *int                  `form:"level" validate:"omitempty,gte=0"`
	Certificate     *multipart.FileHeader `form:"certificate" validate:"omitempty"`
	CertificatePath string                `form:"certificate_path" validate:"omitempty"`
}
