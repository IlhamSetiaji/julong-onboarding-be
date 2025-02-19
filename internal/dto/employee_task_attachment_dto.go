package dto

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IEmployeeTaskAttachmentDTO interface {
	ConvertEntityToResponse(ent *entity.EmployeeTaskAttachment) *response.EmployeeTaskAttachmentResponse
}

type EmployeeTaskAttachmentDTO struct {
	Log   *logrus.Logger
	Viper *viper.Viper
}

func NewEmployeeTaskAttachmentDTO(log *logrus.Logger, viper *viper.Viper) IEmployeeTaskAttachmentDTO {
	return &EmployeeTaskAttachmentDTO{
		Log:   log,
		Viper: viper,
	}
}

func EmployeeTaskAttachmentDTOFactory(log *logrus.Logger, viper *viper.Viper) IEmployeeTaskAttachmentDTO {
	return NewEmployeeTaskAttachmentDTO(log, viper)
}

func (dto *EmployeeTaskAttachmentDTO) ConvertEntityToResponse(ent *entity.EmployeeTaskAttachment) *response.EmployeeTaskAttachmentResponse {
	return &response.EmployeeTaskAttachmentResponse{
		ID:             ent.ID,
		EmployeeTaskID: ent.EmployeeTaskID,
		Path: func() string {
			if ent.Path == "" {
				return ""
			}
			return dto.Viper.GetString("app.url") + ent.Path
		}(),
		PathOrigin: ent.Path,
		CreatedAt:  ent.CreatedAt,
		UpdatedAt:  ent.UpdatedAt,
	}
}
