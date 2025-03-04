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
	ConvertEntitiesToKanbanResponse(entities []entity.EmployeeTask) *response.EmployeeTaskKanbanResponse
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

	var progress, progressVerified int
	if ent.EmployeeTaskChecklists != nil && len(ent.EmployeeTaskChecklists) > 0 {
		for _, checklist := range ent.EmployeeTaskChecklists {
			if checklist.IsChecked == "YES" {
				progress++
			}
			if checklist.IsChecked == "YES" && checklist.VerifiedBy != nil {
				progressVerified++
			}
		}

		progress = (progress * 100) / len(ent.EmployeeTaskChecklists)
		progressVerified = (progressVerified * 100) / len(ent.EmployeeTaskChecklists)
		if progress > 100 {
			progress = 100
		}
		if progressVerified > 100 {
			progressVerified = 100
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
		IsChecklist: func() string {
			if ent.EmployeeTaskChecklists != nil && len(ent.EmployeeTaskChecklists) > 0 {
				return "YES"
			}
			return "NO"
		}(),
		Status:           ent.Status,
		Kanban:           ent.Kanban,
		Notes:            ent.Notes,
		Source:           ent.Source,
		Progress:         progress,
		ProgressVerified: progressVerified,
		CreatedAt:        ent.CreatedAt,
		UpdatedAt:        ent.UpdatedAt,

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

func (dto *EmployeeTaskDTO) ConvertEntitiesToKanbanResponse(entities []entity.EmployeeTask) *response.EmployeeTaskKanbanResponse {
	var toDo []response.EmployeeTaskResponse
	var inProgress []response.EmployeeTaskResponse
	var needReview []response.EmployeeTaskResponse
	var completed []response.EmployeeTaskResponse

	for _, ent := range entities {
		response := dto.ConvertEntityToResponse(&ent)
		switch ent.Kanban {
		case entity.EMPLOYEE_TASK_KANBAN_ENUM_TODO:
			toDo = append(toDo, *response)
		case entity.EPMLOYEE_TASK_KANBAN_ENUM_IN_PROGRESS:
			inProgress = append(inProgress, *response)
		case entity.EMPLOYEE_TASK_KANBAN_ENUM_NEED_REVIEW:
			needReview = append(needReview, *response)
		case entity.EMPLOYEE_TASK_KANBAN_ENUM_COMPLETED:
			completed = append(completed, *response)
		}
	}

	return &response.EmployeeTaskKanbanResponse{
		ToDo:       toDo,
		InProgress: inProgress,
		NeedReview: needReview,
		Completed:  completed,
	}
}
