package dto

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ISkillDTO interface {
	ConvertEntityToResponse(ent *entity.Skill) *response.SkillResponse
}

type SkillDTO struct {
	Log   *logrus.Logger
	Viper *viper.Viper
}

func NewSkillDTO(log *logrus.Logger, viper *viper.Viper) ISkillDTO {
	return &SkillDTO{
		Log:   log,
		Viper: viper,
	}
}

func SkillDTOFactory(log *logrus.Logger, viper *viper.Viper) ISkillDTO {
	return NewSkillDTO(log, viper)
}

func (dto *SkillDTO) ConvertEntityToResponse(ent *entity.Skill) *response.SkillResponse {
	return &response.SkillResponse{
		ID:            ent.ID,
		UserProfileID: ent.UserProfileID,
		Name:          ent.Name,
		Description:   ent.Description,
		Certificate: func() *string {
			if ent.Certificate != "" {
				certificateURL := dto.Viper.GetString("app.url") + ent.Certificate
				return &certificateURL
			}
			return nil
		}(),
		Level: func() *int {
			if ent.Level != nil {
				return ent.Level
			}
			return nil
		}(),
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,
	}
}
