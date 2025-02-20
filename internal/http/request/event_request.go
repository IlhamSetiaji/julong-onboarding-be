package request

type CreateEventRequest struct {
	TemplateTaskID string `json:"template_task_id" validate:"required,uuid"`
	Name           string `json:"name" validate:"required"`
	StartDate      string `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate        string `json:"end_date" validate:"required,datetime=2006-01-02"`
	Description    string `json:"description" validate:"omitempty"`
	Status         string `json:"status" validate:"omitempty,event_status_validation"`
	EventEmployees []struct {
		EmployeeID string `json:"employee_id" validate:"required,uuid"`
	} `json:"event_employees" validate:"required,dive"`
}

type UpdateEventRequest struct {
	ID             string `json:"id" validate:"required,uuid"`
	TemplateTaskID string `json:"template_task_id" validate:"required,uuid"`
	Name           string `json:"name" validate:"required"`
	StartDate      string `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate        string `json:"end_date" validate:"required,datetime=2006-01-02"`
	Description    string `json:"description" validate:"omitempty"`
	Status         string `json:"status" validate:"omitempty,event_status_validation"`
	EventEmployees []struct {
		EmployeeID string `json:"employee_id" validate:"required,uuid"`
	} `json:"event_employees" validate:"required,dive"`
}
