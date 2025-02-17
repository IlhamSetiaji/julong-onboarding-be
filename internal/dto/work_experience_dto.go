package dto

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IWorkExperienceDTO interface {
	ConvertEntityToResponse(ent *entity.WorkExperience) *response.WorkExperienceResponse
}

type WorkExperienceDTO struct {
	Log   *logrus.Logger
	Viper *viper.Viper
}

func NewWorkExperienceDTO(log *logrus.Logger, viper *viper.Viper) IWorkExperienceDTO {
	return &WorkExperienceDTO{
		Log:   log,
		Viper: viper,
	}
}

func WorkExperienceDTOFactory(log *logrus.Logger, viper *viper.Viper) IWorkExperienceDTO {
	return NewWorkExperienceDTO(log, viper)
}

func (dto *WorkExperienceDTO) ConvertEntityToResponse(ent *entity.WorkExperience) *response.WorkExperienceResponse {
	return &response.WorkExperienceResponse{
		ID:             ent.ID,
		UserProfileID:  ent.UserProfileID,
		Name:           ent.Name,
		CompanyName:    ent.CompanyName,
		YearExperience: ent.YearExperience,
		JobDescription: ent.JobDescription,
		Certificate: func() *string {
			if ent.Certificate != "" {
				certificateURL := dto.Viper.GetString("app.url") + ent.Certificate
				return &certificateURL
			}
			return nil
		}(),
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,
	}
}
