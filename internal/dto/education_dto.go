package dto

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IEducationDTO interface {
	ConvertEntityToResponse(ent *entity.Education) *response.EducationResponse
}

type EducationDTO struct {
	Log   *logrus.Logger
	Viper *viper.Viper
}

func NewEducationDTO(log *logrus.Logger, viper *viper.Viper) IEducationDTO {
	return &EducationDTO{
		Log:   log,
		Viper: viper,
	}
}

func EducationDTOFactory(log *logrus.Logger, viper *viper.Viper) IEducationDTO {
	return NewEducationDTO(log, viper)
}

func (dto *EducationDTO) ConvertEntityToResponse(ent *entity.Education) *response.EducationResponse {
	return &response.EducationResponse{
		ID:             ent.ID,
		UserProfileID:  ent.UserProfileID,
		EducationLevel: ent.EducationLevel,
		Major:          ent.Major,
		SchoolName:     ent.SchoolName,
		GraduateYear:   ent.GraduateYear,
		EndDate:        ent.EndDate,
		Certificate: func() *string {
			if ent.Certificate != "" {
				certificateURL := dto.Viper.GetString("app.url") + ent.Certificate
				return &certificateURL
			}
			return nil
		}(),
		Gpa: func() *float64 {
			if ent.Gpa != nil {
				return ent.Gpa
			}
			return nil
		}(),
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,
	}
}
