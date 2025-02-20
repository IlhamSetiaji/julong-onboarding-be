package response

import (
	"time"

	"github.com/google/uuid"
)

type EventEmployeeResponse struct {
	ID         uuid.UUID  `json:"id"`
	EventID    uuid.UUID  `json:"event_id"`
	EmployeeID *uuid.UUID `json:"employee_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`

	EmployeeName string `json:"employee_name"`
}
