package request

type EmployeeTaskChecklistRequest struct {
	ID         *string `form:"id" validate:"omitempty"`
	Name       string  `form:"name" validate:"required"`
	IsChecked  *string `form:"is_checked" validate:"omitempty"`
	VerifiedBy *string `form:"verified_by" validate:"omitempty,uuid"`
}
