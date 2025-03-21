package response

import (
	"time"

	"github.com/google/uuid"
)

type QuestionOptionResponse struct {
	ID         uuid.UUID `json:"id"`
	QuestionID uuid.UUID `json:"question_id"`
	OptionText string    `json:"option_text"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
