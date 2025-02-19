package response

import (
	"time"

	"github.com/google/uuid"
)

type EmployeeTaskChecklistResponse struct {
	ID             uuid.UUID  `json:"id"`
	EmployeeTaskID uuid.UUID  `json:"employee_task_id"`
	Name           string     `json:"name"`
	IsChecked      string     `json:"is_checked"`
	VerifiedBy     *uuid.UUID `json:"verified_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	VerifiedByName string `json:"verified_by_name"`
}
