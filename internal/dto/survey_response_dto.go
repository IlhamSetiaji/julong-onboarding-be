package dto

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ISurveyResponseDTO interface {
	ConvertEntityToResponse(ent *entity.SurveyResponse) *response.SurveyResponseResponse
}

type SurveyResponseDTO struct {
	Log   *logrus.Logger
	Viper *viper.Viper
}

func NewSurveyResponseDTO(log *logrus.Logger, viper *viper.Viper) ISurveyResponseDTO {
	return &SurveyResponseDTO{
		Log:   log,
		Viper: viper,
	}
}

func SurveyResponseDTOFactory(log *logrus.Logger, viper *viper.Viper) ISurveyResponseDTO {
	return NewSurveyResponseDTO(log, viper)
}

func (dto *SurveyResponseDTO) ConvertEntityToResponse(ent *entity.SurveyResponse) *response.SurveyResponseResponse {
	return &response.SurveyResponseResponse{
		ID:               ent.ID,
		SurveyTemplateID: ent.SurveyTemplateID,
		EmployeeTaskID:   ent.EmployeeTaskID,
		QuestionID:       ent.QuestionID,
		Answer:           ent.Answer,
		AnswerFile: func() string {
			if ent.AnswerFile != "" {
				return dto.Viper.GetString("app.url") + ent.AnswerFile
			}
			return ""
		}(),
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,
	}
}
