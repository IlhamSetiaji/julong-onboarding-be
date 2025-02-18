package dto

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ICoverDTO interface {
	ConvertEntityToResponse(ent *entity.Cover) *response.CoverResponse
}

type CoverDTO struct {
	Log   *logrus.Logger
	Viper *viper.Viper
}

func NewCoverDTO(log *logrus.Logger, viper *viper.Viper) ICoverDTO {
	return &CoverDTO{
		Log:   log,
		Viper: viper,
	}
}

func CoverDTOFactory(log *logrus.Logger, viper *viper.Viper) ICoverDTO {
	return NewCoverDTO(log, viper)
}

func (dto *CoverDTO) ConvertEntityToResponse(ent *entity.Cover) *response.CoverResponse {
	return &response.CoverResponse{
		ID: ent.ID,
		Path: func() string {
			if ent.Path == "" {
				return ""
			}
			return ent.Path
		}(),
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,
	}
}
