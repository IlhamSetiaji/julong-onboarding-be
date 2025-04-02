package request

type CreateSurveyTemplateRequest struct {
	Title string `json:"title" validate:"required"`
}

type UpdateSurveyTemplateRequest struct {
	ID     string `json:"id" validate:"required,uuid"`
	Title  string `json:"title" validate:"required"`
	Status string `json:"status" validate:"required,survey_template_status_validation"`
}
