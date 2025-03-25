package dto

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IQuestionDTO interface {
	ConvertEntityToResponse(ent *entity.Question) *response.QuestionResponse
}

type QuestionDTO struct {
	Log               *logrus.Logger
	Viper             *viper.Viper
	AnswerTypeDTO     IAnswerTypeDTO
	QuestionOptionDTO IQuestionOptionDTO
	SurveyResponseDTO ISurveyResponseDTO
}

func NewQuestionDTO(
	log *logrus.Logger,
	viper *viper.Viper,
	answerTypeDTO IAnswerTypeDTO,
	questionOptionDTO IQuestionOptionDTO,
	surveyResponseDTO ISurveyResponseDTO,
) IQuestionDTO {
	return &QuestionDTO{
		Log:               log,
		Viper:             viper,
		AnswerTypeDTO:     answerTypeDTO,
		QuestionOptionDTO: questionOptionDTO,
		SurveyResponseDTO: surveyResponseDTO,
	}
}

func QuestionDTOFactory(log *logrus.Logger, viper *viper.Viper) IQuestionDTO {
	answerTypeDTO := AnswerTypeDTOFactory(log)
	questionOptionDTO := QuestionOptionDTOFactory(log, viper)
	surveyResponseDTO := SurveyResponseDTOFactory(log, viper)
	return NewQuestionDTO(log, viper, answerTypeDTO, questionOptionDTO, surveyResponseDTO)
}

func (dto *QuestionDTO) ConvertEntityToResponse(ent *entity.Question) *response.QuestionResponse {
	return &response.QuestionResponse{
		ID:               ent.ID,
		SurveyTemplateID: ent.SurveyTemplateID,
		AnswerTypeID:     ent.AnswerTypeID,
		Question:         ent.Question,
		Attachment: func() *string {
			if ent.Attachment == nil {
				return ent.Attachment
			}
			path := dto.Viper.GetString("app.url") + *ent.Attachment
			return &path
		}(),
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,

		AnswerType: func() *response.AnswerTypeResponse {
			if ent.AnswerType == nil {
				return nil
			}
			return dto.AnswerTypeDTO.ConvertEntityToResponse(ent.AnswerType)
		}(),
		QuestionOptions: func() []response.QuestionOptionResponse {
			if len(ent.QuestionOptions) == 0 {
				return nil
			}
			var responses []response.QuestionOptionResponse
			for _, questionOption := range ent.QuestionOptions {
				responses = append(responses, *dto.QuestionOptionDTO.ConvertEntityToResponse(&questionOption))
			}
			return responses
		}(),
		SurveyResponses: func() []response.SurveyResponseResponse {
			if len(ent.SurveyResponses) == 0 {
				return nil
			}
			var responses []response.SurveyResponseResponse
			for _, surveyResponse := range ent.SurveyResponses {
				responses = append(responses, *dto.SurveyResponseDTO.ConvertEntityToResponse(&surveyResponse))
			}
			return responses
		}(),
	}
}
