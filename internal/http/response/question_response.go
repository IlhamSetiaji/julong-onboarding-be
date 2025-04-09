package response

import (
	"time"

	"github.com/google/uuid"
)

type QuestionResponse struct {
	ID               uuid.UUID `json:"id"`
	SurveyTemplateID uuid.UUID `json:"survey_template_id"`
	AnswerTypeID     uuid.UUID `json:"answer_type_id"`
	Question         string    `json:"question"`
	Attachment       *string   `json:"attachment"`
	IsCompleted      string    `json:"is_completed"`
	Number           int       `json:"number"`
	MaxStars         int       `json:"max_stars"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	AnswerType      *AnswerTypeResponse      `json:"answer_type"`
	QuestionOptions []QuestionOptionResponse `json:"question_options"`
	SurveyResponses []SurveyResponseResponse `json:"survey_responses"`
}
