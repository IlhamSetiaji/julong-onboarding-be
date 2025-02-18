package request

import "mime/multipart"

type CreateCoverRequest struct {
	File *multipart.FileHeader `form:"file" validate:"required"`
	Path string                `form:"path" validate:"omitempty"`
}

type UpdateCoverRequest struct {
	ID   string                `form:"id" validate:"required,uuid"`
	File *multipart.FileHeader `form:"file" validate:"required"`
	Path string                `form:"path" validate:"omitempty"`
}
