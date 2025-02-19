package request

import "mime/multipart"

type EmployeeTaskAttachmentRequest struct {
	File *multipart.FileHeader `form:"file" validate:"required"`
	Path string                `form:"path" validate:"omitempty"`
}
