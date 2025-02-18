package request

import "mime/multipart"

type TemplateTaskAttachmentRequest struct {
	File *multipart.FileHeader `form:"file" validate:"required"`
	Path string                `form:"path" validate:"omitempty"`
}
