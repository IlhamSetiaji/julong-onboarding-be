package dto

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IQuestionOptionDTO interface {
	ConvertEntityToResponse(ent *entity.QuestionOption) *response.QuestionOptionResponse
}

type QuestionOptionDTO struct {
	Log   *logrus.Logger
	Viper *viper.Viper
}

func NewQuestionOptionDTO(log *logrus.Logger, viper *viper.Viper) IQuestionOptionDTO {
	return &QuestionOptionDTO{
		Log:   log,
		Viper: viper,
	}
}

func QuestionOptionDTOFactory(log *logrus.Logger, viper *viper.Viper) IQuestionOptionDTO {
	return NewQuestionOptionDTO(log, viper)
}

func (dto *QuestionOptionDTO) ConvertEntityToResponse(ent *entity.QuestionOption) *response.QuestionOptionResponse {
	return &response.QuestionOptionResponse{
		ID:         ent.ID,
		QuestionID: ent.QuestionID,
		OptionText: ent.OptionText,
		CreatedAt:  ent.CreatedAt,
		UpdatedAt:  ent.UpdatedAt,
	}
}
