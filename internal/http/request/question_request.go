package request

import "mime/multipart"

type CreateOrUpdateQuestions struct {
	SurveyTemplateID   string            `form:"survey_template_id" validate:"required,uuid"`
	Questions          []QuestionRequest `json:"questions" validate:"required,dive"`
	DeletedQuestionIDs []string          `json:"deleted_question_ids" validate:"omitempty,dive,uuid"`
}

type QuestionRequest struct {
	ID              string                  `form:"id" validate:"required,uuid"`
	AnswerTypeID    string                  `form:"answer_type_id" validate:"required,uuid"`
	Question        string                  `form:"question" validate:"required"`
	Attachment      *multipart.FileHeader   `form:"attachment" validate:"omitempty"`
	AttachmentPath  string                  `form:"attachment_path" validate:"omitempty"`
	QuestionOptions []QuestionOptionRequest `form:"question_options" validate:"omitempty,dive"`
}

type QuestionOptionRequest struct {
	OptionText string `form:"option_text" validate:"required"`
}
