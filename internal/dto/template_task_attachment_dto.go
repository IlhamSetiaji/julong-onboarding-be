package dto

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ITemplateTaskAttachmentDTO interface {
	ConvertEntityToResponse(ent *entity.TemplateTaskAttachment) *response.TemplateTaskAttachmentResponse
}

type TemplateTaskAttachmentDTO struct {
	Log   *logrus.Logger
	Viper *viper.Viper
}

func NewTemplateTaskAttachmentDTO(log *logrus.Logger, viper *viper.Viper) ITemplateTaskAttachmentDTO {
	return &TemplateTaskAttachmentDTO{
		Log:   log,
		Viper: viper,
	}
}

func TemplateTaskAttachmentDTOFactory(log *logrus.Logger, viper *viper.Viper) ITemplateTaskAttachmentDTO {
	return NewTemplateTaskAttachmentDTO(log, viper)
}

func (dto *TemplateTaskAttachmentDTO) ConvertEntityToResponse(ent *entity.TemplateTaskAttachment) *response.TemplateTaskAttachmentResponse {
	return &response.TemplateTaskAttachmentResponse{
		ID:             ent.ID,
		TemplateTaskID: ent.TemplateTaskID,
		Path: func() string {
			if ent.Path == "" {
				return ""
			}
			return dto.Viper.GetString("app.url") + ent.Path
		}(),
		PathOrigin: ent.Path,
		CreatedAt:  ent.CreatedAt,
		UpdatedAt:  ent.UpdatedAt,
	}
}
