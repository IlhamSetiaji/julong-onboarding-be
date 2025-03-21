package dto

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IAnswerTypeDTO interface {
	ConvertEntityToResponse(ent *entity.AnswerType) *response.AnswerTypeResponse
}

type AnswerTypeDTO struct {
	Log   *logrus.Logger
	Viper *viper.Viper
}

func NewAnswerTypeDTO(log *logrus.Logger, viper *viper.Viper) IAnswerTypeDTO {
	return &AnswerTypeDTO{
		Log:   log,
		Viper: viper,
	}
}

func AnswerTypeDTOFactory(log *logrus.Logger, viper *viper.Viper) IAnswerTypeDTO {
	return NewAnswerTypeDTO(log, viper)
}

func (dto *AnswerTypeDTO) ConvertEntityToResponse(ent *entity.AnswerType) *response.AnswerTypeResponse {
	return &response.AnswerTypeResponse{
		ID:        ent.ID,
		Name:      ent.Name,
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,
	}
}
