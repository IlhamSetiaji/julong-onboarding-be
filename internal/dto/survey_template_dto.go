package dto

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ISurveyTemplateDTO interface {
	ConvertEntityToResponse(ent *entity.SurveyTemplate) *response.SurveyTemplateResponse
}

type SurveyTemplateDTO struct {
	Log         *logrus.Logger
	Viper       *viper.Viper
	QuestionDTO IQuestionDTO
}

func NewSurveyTemplateDTO(log *logrus.Logger, viper *viper.Viper, questionDTO IQuestionDTO) ISurveyTemplateDTO {
	return &SurveyTemplateDTO{
		Log:         log,
		Viper:       viper,
		QuestionDTO: questionDTO,
	}
}

func SurveyTemplateDTOFactory(log *logrus.Logger, viper *viper.Viper) ISurveyTemplateDTO {
	questionDTO := QuestionDTOFactory(log, viper)
	return NewSurveyTemplateDTO(log, viper, questionDTO)
}

func (dto *SurveyTemplateDTO) ConvertEntityToResponse(ent *entity.SurveyTemplate) *response.SurveyTemplateResponse {
	return &response.SurveyTemplateResponse{
		ID:           ent.ID,
		SurveyNumber: ent.SurveyNumber,
		Title:        ent.Title,
		Status:       ent.Status,
		CreatedAt:    ent.CreatedAt,
		UpdatedAt:    ent.UpdatedAt,

		Questions: func() []response.QuestionResponse {
			if ent.Questions == nil {
				return nil
			}
			var questions []response.QuestionResponse
			for _, question := range ent.Questions {
				questions = append(questions, *dto.QuestionDTO.ConvertEntityToResponse(&question))
			}

			return questions
		}(),
	}
}
