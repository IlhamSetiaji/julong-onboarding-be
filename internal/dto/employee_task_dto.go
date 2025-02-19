package dto

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IEmployeeTaskDTO interface {
	ConvertEntityToResponse(ent *entity.EmployeeTask) *response.EmployeeTaskResponse
}

type EmployeeTaskDTO struct {
	Log                       *logrus.Logger
	Viper                     *viper.Viper
	EmployeeTaskAttachmentDTO IEmployeeTaskAttachmentDTO
	EmployeeTaskChecklistDTO  IEmployeeTaskChecklistDTO
	EmployeeMessage           messaging.IEmployeeMessage
}

func NewEmployeeTaskDTO(log *logrus.Logger, viper *viper.Viper, employeeTaskAttachmentDTO IEmployeeTaskAttachmentDTO, employeeTaskChecklistDTO IEmployeeTaskChecklistDTO, employeeMessage messaging.IEmployeeMessage) IEmployeeTaskDTO {
	return &EmployeeTaskDTO{
		Log:                       log,
		Viper:                     viper,
		EmployeeTaskAttachmentDTO: employeeTaskAttachmentDTO,
		EmployeeTaskChecklistDTO:  employeeTaskChecklistDTO,
		EmployeeMessage:           employeeMessage,
	}
}

func EmployeeTaskDTOFactory(log *logrus.Logger, viper *viper.Viper) IEmployeeTaskDTO {
	employeeTaskAttachmentDTO := EmployeeTaskAttachmentDTOFactory(log, viper)
	employeeTaskChecklistDTO := EmployeeTaskChecklistDTOFactory(log, viper)
	employeeMessage := messaging.EmployeeMessageFactory(log)
	return NewEmployeeTaskDTO(log, viper, employeeTaskAttachmentDTO, employeeTaskChecklistDTO, employeeMessage)
}

func (dto *EmployeeTaskDTO) ConvertEntityToResponse(ent *entity.EmployeeTask) *response.EmployeeTaskResponse {
	var verifiedByName string
	if ent.VerifiedBy != nil {
		employee, err := dto.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
			ID: ent.VerifiedBy.String(),
		})
		if err != nil {
			dto.Log.Errorf("[ProjectPicDTO.ConvertEntityToResponse] " + err.Error())
			verifiedByName = ""
		} else {
			verifiedByName = employee.Name
		}
	}

	var employeeName string
	if ent.EmployeeID != nil {
		employee, err := dto.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
			ID: ent.EmployeeID.String(),
		})
		if err != nil {
			dto.Log.Errorf("[ProjectPicDTO.ConvertEntityToResponse] " + err.Error())
			employeeName = ""
		} else {
			employeeName = employee.Name
		}
	}

	return &response.EmployeeTaskResponse{
		ID: ent.ID,
		CoverPath: func() *string {
			if ent.CoverPath == nil {
				return nil
			}
			path := dto.Viper.GetString("app.url") + *ent.CoverPath
			return &path
		}(),
		EmployeeID:     ent.EmployeeID,
		TemplateTaskID: ent.TemplateTaskID,
		VerifiedBy:     ent.VerifiedBy,
		Name:           ent.Name,
		Priority:       ent.Priority,
		Description:    ent.Description,
		StartDate:      ent.StartDate,
		EndDate:        ent.EndDate,
		IsDone:         ent.IsDone,
		Proof: func() *string {
			if ent.Proof == nil {
				return nil
			}
			path := dto.Viper.GetString("app.url") + *ent.Proof
			return &path
		}(),
		Status:    ent.Status,
		Kanban:    ent.Kanban,
		Notes:     ent.Notes,
		Source:    ent.Source,
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,

		VerifiedByName: verifiedByName,
		EmployeeName:   employeeName,
		EmployeeTaskChecklists: func() []response.EmployeeTaskChecklistResponse {
			var checklists []response.EmployeeTaskChecklistResponse
			for _, checklist := range ent.EmployeeTaskChecklists {
				response := dto.EmployeeTaskChecklistDTO.ConvertEntityToResponse(&checklist)
				checklists = append(checklists, *response)
			}
			return checklists
		}(),
		EmployeeTaskAttachments: func() []response.EmployeeTaskAttachmentResponse {
			var attachments []response.EmployeeTaskAttachmentResponse
			for _, attachment := range ent.EmployeeTaskAttachments {
				response := dto.EmployeeTaskAttachmentDTO.ConvertEntityToResponse(&attachment)
				attachments = append(attachments, *response)
			}
			return attachments
		}(),
	}
}
