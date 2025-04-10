package request

import "mime/multipart"

type SurveyResponseRequest struct {
	QuestionID string          `form:"question_id" validate:"required,uuid"`
	Answers    []AnswerRequest `form:"answers" validate:"omitempty,dive"`
	// DeletedAnswerIDs []string        `form:"deleted_answer_ids" validate:"omitempty,dive,uuid"`
}

type AnswerRequest struct {
	ID               *string               `form:"id" validate:"omitempty,uuid"`
	SurveyTemplateID string                `form:"survey_template_id" validate:"required,uuid"`
	EmployeeTaskID   string                `form:"employee_task_id" validate:"required,uuid"`
	Answer           string                `form:"answer" validate:"omitempty"`
	AnswerFile       *multipart.FileHeader `form:"answer_file" validate:"omitempty"`
	AnswerPath       string                `form:"answer_path" validate:"omitempty"`
}

type SurveyResponseBulkRequest struct {
	SurveyTemplateID string              `form:"survey_template_id" validate:"required,uuid"`
	EmployeeTaskID   string              `form:"employee_task_id" validate:"required,uuid"`
	Answers          []AnswerBulkRequest `form:"answers" validate:"omitempty,dive"`
}

type AnswerBulkRequest struct {
	ID         *string               `form:"id" validate:"omitempty,uuid"`
	QuestionID string                `form:"question_id" validate:"required,uuid"`
	Answer     string                `form:"answer" validate:"omitempty"`
	AnswerFile *multipart.FileHeader `form:"answer_file" validate:"omitempty"`
	AnswerPath string                `form:"answer_path" validate:"omitempty"`
}
