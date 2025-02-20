package dto

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IEventDTO interface {
	ConvertEntityToResponse(ent *entity.Event) *response.EventResponse
}

type EventDTO struct {
	Log              *logrus.Logger
	Viper            *viper.Viper
	EventEmployeeDTO IEventEmployeeDTO
	TemplateTaskDTO  ITemplateTaskDTO
}

func NewEventDTO(log *logrus.Logger, viper *viper.Viper, eventEmployeeDTO IEventEmployeeDTO, templateTaskDTO ITemplateTaskDTO) IEventDTO {
	return &EventDTO{
		Log:              log,
		Viper:            viper,
		EventEmployeeDTO: eventEmployeeDTO,
		TemplateTaskDTO:  templateTaskDTO,
	}
}

func EventDTOFactory(log *logrus.Logger, viper *viper.Viper) IEventDTO {
	eventEmployeeDTO := EventEmployeeDTOFactory(log, viper)
	templateTaskDTO := TemplateTaskDTOFactory(log, viper)
	return NewEventDTO(log, viper, eventEmployeeDTO, templateTaskDTO)
}

func (dto *EventDTO) ConvertEntityToResponse(ent *entity.Event) *response.EventResponse {
	return &response.EventResponse{
		ID:             ent.ID,
		TemplateTaskID: ent.TemplateTaskID,
		Name:           ent.Name,
		StartDate:      ent.StartDate,
		EndDate:        ent.EndDate,
		Description:    ent.Description,
		Status:         ent.Status,
		CreatedAt:      ent.CreatedAt,
		UpdatedAt:      ent.UpdatedAt,

		TemplateTask: func() *response.TemplateTaskResponse {
			templateTask := ent.TemplateTask
			if templateTask == nil {
				return nil
			}
			return dto.TemplateTaskDTO.ConvertEntityToResponse(templateTask)
		}(),
		EventEmployees: func() []response.EventEmployeeResponse {
			eventEmployees := ent.EventEmployees
			if eventEmployees == nil {
				return nil
			}

			var responses []response.EventEmployeeResponse
			for _, eventEmployee := range eventEmployees {
				resp := dto.EventEmployeeDTO.ConvertEntityToResponse(&eventEmployee)
				responses = append(responses, *resp)
			}
			return responses
		}(),
	}
}
