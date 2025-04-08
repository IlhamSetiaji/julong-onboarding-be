package response

import (
	"github.com/google/uuid"
)

type CheckJobExistMessageResponse struct {
	JobID uuid.UUID `json:"job_id"`
	Exist bool      `json:"exist"`
}

type SendFindJobByIDMessageResponse struct {
	JobID     uuid.UUID `json:"job_id"`
	Name      string    `json:"name"`
	MidsuitID string    `json:"midsuit_id"`
}

type SendFindJobLevelByIDMessageResponse struct {
	JobLevelID uuid.UUID `json:"job_level_id"`
	Name       string    `json:"name"`
	Level      float64   `json:"level"`
	MidsuitID  string    `json:"midsuit_id"`
}
