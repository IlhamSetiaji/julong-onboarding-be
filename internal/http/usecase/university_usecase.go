package usecase

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/dto"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/repository"
	"github.com/sirupsen/logrus"
)

type IUniversityUseCase interface {
	FindAll() ([]*response.UniversityResponse, error)
}

type UniversityUseCase struct {
	Log        *logrus.Logger
	Repository repository.IUniversityRepository
	DTO        dto.IUniversityDTO
}

func NewUniversityUseCase(
	log *logrus.Logger,
	repo repository.IUniversityRepository,
	dto dto.IUniversityDTO,
) IUniversityUseCase {
	return &UniversityUseCase{
		Log:        log,
		Repository: repo,
		DTO:        dto,
	}
}

func UniversityUseCaseFactory(log *logrus.Logger) IUniversityUseCase {
	repo := repository.UniversityRepositoryFactory(log)
	dto := dto.UniversityDTOFactory(log)
	return NewUniversityUseCase(log, repo, dto)
}

func (uc *UniversityUseCase) FindAll() ([]*response.UniversityResponse, error) {
	entities, err := uc.Repository.FindAll()
	if err != nil {
		uc.Log.Error("[UniversityUseCase.FindAll] " + err.Error())
		return nil, err
	}

	var responses []*response.UniversityResponse
	for _, entity := range entities {
		responses = append(responses, uc.DTO.ConvertEntityToResponse(entity))
	}

	return responses, nil
}
