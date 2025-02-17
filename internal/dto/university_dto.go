package dto

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IUniversityDTO interface {
	ConvertEntityToResponse(ent *entity.University) *response.UniversityResponse
}

type UniversityDTO struct {
	Log *logrus.Logger
}

func NewUniversityDTO(log *logrus.Logger) IUniversityDTO {
	return &UniversityDTO{
		Log: log,
	}
}

func UniversityDTOFactory(log *logrus.Logger) IUniversityDTO {
	return NewUniversityDTO(log)
}

func (dto *UniversityDTO) ConvertEntityToResponse(ent *entity.University) *response.UniversityResponse {
	return &response.UniversityResponse{
		ID:           ent.ID,
		Name:         ent.Name,
		Country:      ent.Country,
		AlphaTwoCode: ent.AlphaTwoCode,
	}
}
