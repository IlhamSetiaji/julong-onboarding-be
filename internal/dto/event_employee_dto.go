package dto

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IEventEmployeeDTO interface {
	ConvertEntityToResponse(ent *entity.EventEmployee) *response.EventEmployeeResponse
}

type EventEmployeeDTO struct {
	Log             *logrus.Logger
	Viper           *viper.Viper
	EmployeeMessage messaging.IEmployeeMessage
}

func NewEventEmployeeDTO(log *logrus.Logger, viper *viper.Viper, employeeMessage messaging.IEmployeeMessage) IEventEmployeeDTO {
	return &EventEmployeeDTO{
		Log:             log,
		Viper:           viper,
		EmployeeMessage: employeeMessage,
	}
}

func EventEmployeeDTOFactory(log *logrus.Logger, viper *viper.Viper) IEventEmployeeDTO {
	employeeMessage := messaging.EmployeeMessageFactory(log)
	return NewEventEmployeeDTO(log, viper, employeeMessage)
}

func (dto *EventEmployeeDTO) ConvertEntityToResponse(ent *entity.EventEmployee) *response.EventEmployeeResponse {
	var employeeName string
	if ent.EmployeeID != nil {
		employee, err := dto.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
			ID: ent.EmployeeID.String(),
		})
		if err != nil {
			dto.Log.Errorf("[EventEmployeeDTO.ConvertEntityToResponse] " + err.Error())
			employeeName = ""
		} else {
			employeeName = employee.Name
		}
	}

	return &response.EventEmployeeResponse{
		ID:           ent.ID,
		EventID:      ent.EventID,
		EmployeeID:   ent.EmployeeID,
		EmployeeName: employeeName,
		CreatedAt:    ent.CreatedAt,
		UpdatedAt:    ent.UpdatedAt,
	}
}
