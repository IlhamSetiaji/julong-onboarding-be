package request

import "mime/multipart"

type CreateOrUpdateQuestions struct {
	SurveyTemplateID   string            `json:"survey_template_id" validate:"omitempty,uuid"`
	Title              string            `json:"title" validate:"required"`
	Questions          []QuestionRequest `json:"questions" validate:"omitempty,dive"`
	DeletedQuestionIDs []string          `json:"deleted_question_ids" validate:"omitempty,dive,uuid"`
}

type QuestionRequest struct {
	ID              string                  `json:"id" validate:"omitempty,uuid"`
	AnswerTypeID    string                  `json:"answer_type_id" validate:"required,uuid"`
	Question        string                  `json:"question" validate:"omitempty"`
	MaxStars        int                     `json:"max_stars" validate:"omitempty"`
	Attachment      *multipart.FileHeader   `json:"attachment" validate:"omitempty"`
	AttachmentPath  string                  `json:"attachment_path" validate:"omitempty"`
	QuestionOptions []QuestionOptionRequest `json:"question_options" validate:"omitempty,dive"`
}

type QuestionOptionRequest struct {
	OptionText string `json:"option_text" validate:"required"`
}
