package dto

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ITemplateTaskChecklistDTO interface {
	ConvertEntityToResponse(ent *entity.TemplateTaskChecklist) *response.TemplateTaskChecklistResponse
}

type TemplateTaskChecklistDTO struct {
	Log   *logrus.Logger
	Viper *viper.Viper
}

func NewTemplateTaskChecklistDTO(log *logrus.Logger, viper *viper.Viper) ITemplateTaskChecklistDTO {
	return &TemplateTaskChecklistDTO{
		Log:   log,
		Viper: viper,
	}
}

func TemplateTaskChecklistDTOFactory(log *logrus.Logger, viper *viper.Viper) ITemplateTaskChecklistDTO {
	return NewTemplateTaskChecklistDTO(log, viper)
}

func (dto *TemplateTaskChecklistDTO) ConvertEntityToResponse(ent *entity.TemplateTaskChecklist) *response.TemplateTaskChecklistResponse {
	return &response.TemplateTaskChecklistResponse{
		ID:             ent.ID,
		TemplateTaskID: ent.TemplateTaskID,
		Name:           ent.Name,
		CreatedAt:      ent.CreatedAt,
		UpdatedAt:      ent.UpdatedAt,
	}
}
