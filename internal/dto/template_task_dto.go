package dto

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ITemplateTaskDTO interface {
	ConvertEntityToResponse(ent *entity.TemplateTask) *response.TemplateTaskResponse
}

type TemplateTaskDTO struct {
	Log                       *logrus.Logger
	Viper                     *viper.Viper
	TemplateTaskAttachmentDTO ITemplateTaskAttachmentDTO
	TemplateTaskChecklistDTO  ITemplateTaskChecklistDTO
	SurveyTemplateDTO         ISurveyTemplateDTO
}

func NewTemplateTaskDTO(log *logrus.Logger, viper *viper.Viper, templateTaskAttachmentDTO ITemplateTaskAttachmentDTO, templateTaskChecklistDTO ITemplateTaskChecklistDTO, surveyTemplateDTO ISurveyTemplateDTO) ITemplateTaskDTO {
	return &TemplateTaskDTO{
		Log:                       log,
		Viper:                     viper,
		TemplateTaskAttachmentDTO: templateTaskAttachmentDTO,
		TemplateTaskChecklistDTO:  templateTaskChecklistDTO,
		SurveyTemplateDTO:         surveyTemplateDTO,
	}
}

func TemplateTaskDTOFactory(log *logrus.Logger, viper *viper.Viper) ITemplateTaskDTO {
	templateTaskAttachmentDTO := TemplateTaskAttachmentDTOFactory(log, viper)
	templateTaskChecklistDTO := TemplateTaskChecklistDTOFactory(log, viper)
	surveyTemplateDTO := SurveyTemplateDTOFactory(log, viper)
	return NewTemplateTaskDTO(log, viper, templateTaskAttachmentDTO, templateTaskChecklistDTO, surveyTemplateDTO)
}

func (dto *TemplateTaskDTO) ConvertEntityToResponse(ent *entity.TemplateTask) *response.TemplateTaskResponse {
	return &response.TemplateTaskResponse{
		ID:   ent.ID,
		Name: ent.Name,
		SurveyTemplateID: func() *uuid.UUID {
			if ent.SurveyTemplateID == nil {
				return nil
			}

			return ent.SurveyTemplateID
		}(),
		Description: ent.Description,
		CoverPath: func() *string {
			if ent.CoverPath == nil {
				return nil
			}
			combinedUrl := dto.Viper.GetString("app.url") + *ent.CoverPath
			return &combinedUrl
		}(),
		CoverPathOrigin: ent.CoverPath,
		Priority:        ent.Priority,
		DueDuration: func() *int {
			if ent.DueDuration == nil {
				return nil
			}
			return ent.DueDuration
		}(),
		Status:           ent.Status,
		Source:           ent.Source,
		OrganizationType: ent.OrganizationType,
		CreatedAt:        ent.CreatedAt,
		UpdatedAt:        ent.UpdatedAt,
		TemplateTaskAttachments: func() []response.TemplateTaskAttachmentResponse {
			var res []response.TemplateTaskAttachmentResponse
			if ent.TemplateTaskAttachments == nil || len(ent.TemplateTaskAttachments) == 0 {
				return nil
			}
			for _, attachment := range ent.TemplateTaskAttachments {
				resp := dto.TemplateTaskAttachmentDTO.ConvertEntityToResponse(&attachment)
				res = append(res, *resp)
			}
			return res
		}(),
		TemplateTaskChecklists: func() []response.TemplateTaskChecklistResponse {
			var res []response.TemplateTaskChecklistResponse
			if ent.TemplateTaskChecklists == nil || len(ent.TemplateTaskChecklists) == 0 {
				return nil
			}
			for _, checklist := range ent.TemplateTaskChecklists {
				resp := dto.TemplateTaskChecklistDTO.ConvertEntityToResponse(&checklist)
				res = append(res, *resp)
			}
			return res
		}(),
		SurveyTemplate: func() *response.SurveyTemplateResponse {
			if ent.SurveyTemplate == nil {
				return nil
			}
			resp := dto.SurveyTemplateDTO.ConvertEntityToResponse(ent.SurveyTemplate)
			return resp
		}(),
	}
}
