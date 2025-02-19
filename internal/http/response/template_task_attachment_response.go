package response

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TemplateTaskAttachmentResponse struct {
	gorm.Model     `json:"-"`
	ID             uuid.UUID `json:"id"`
	TemplateTaskID uuid.UUID `json:"template_task_id"`
	Path           string    `json:"path"`
	PathOrigin     string    `json:"path_origin"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
