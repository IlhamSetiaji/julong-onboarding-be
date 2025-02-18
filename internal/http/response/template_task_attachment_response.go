package response

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TemplateTaskAttachmentResponse struct {
	gorm.Model     `json:"-"`
	ID             uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	TemplateTaskID uuid.UUID `json:"template_task_id" gorm:"type:char(36);not null"`
	Path           string    `json:"path" gorm:"type:varchar(255);not null"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
