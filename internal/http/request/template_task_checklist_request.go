package request

type TemplateTaskChecklistRequest struct {
	ID   string `form:"id" validate:"omitempty"`
	Name string `form:"name" validate:"required"`
}
