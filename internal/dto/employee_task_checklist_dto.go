package dto

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IEmployeeTaskChecklistDTO interface {
	ConvertEntityToResponse(ent *entity.EmployeeTaskChecklist) *response.EmployeeTaskChecklistResponse
}

type EmployeeTaskChecklistDTO struct {
	Log   *logrus.Logger
	Viper *viper.Viper
}

func NewEmployeeTaskChecklistDTO(log *logrus.Logger, viper *viper.Viper) IEmployeeTaskChecklistDTO {
	return &EmployeeTaskChecklistDTO{
		Log:   log,
		Viper: viper,
	}
}

func EmployeeTaskChecklistDTOFactory(log *logrus.Logger, viper *viper.Viper) IEmployeeTaskChecklistDTO {
	return NewEmployeeTaskChecklistDTO(log, viper)
}

func (dto *EmployeeTaskChecklistDTO) ConvertEntityToResponse(ent *entity.EmployeeTaskChecklist) *response.EmployeeTaskChecklistResponse {
	return &response.EmployeeTaskChecklistResponse{
		ID:             ent.ID,
		EmployeeTaskID: ent.EmployeeTaskID,
		Name:           ent.Name,
		IsChecked:      ent.IsChecked,
		VerifiedBy:     ent.VerifiedBy,
		MidsuitID:      ent.MidsuitID,
		CreatedAt:      ent.CreatedAt,
		UpdatedAt:      ent.UpdatedAt,
	}
}
