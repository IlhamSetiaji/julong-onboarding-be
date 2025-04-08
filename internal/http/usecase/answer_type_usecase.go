package usecase

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/dto"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/repository"
	"github.com/sirupsen/logrus"
)

type IAnswerTypeUseCase interface {
	FindAll() ([]*response.AnswerTypeResponse, error)
	FindByID(id string) (*response.AnswerTypeResponse, error)
}

type AnswerTypeUseCase struct {
	Log        *logrus.Logger
	Repository repository.IAnswerTypeRepository
	DTO        dto.IAnswerTypeDTO
}

func NewAnswerTypeUseCase(
	log *logrus.Logger,
	repository repository.IAnswerTypeRepository,
	atDTO dto.IAnswerTypeDTO,
) *AnswerTypeUseCase {
	return &AnswerTypeUseCase{
		Log:        log,
		Repository: repository,
		DTO:        atDTO,
	}
}

func AnswerTypeUseCaseFactory(
	log *logrus.Logger,
) IAnswerTypeUseCase {
	return NewAnswerTypeUseCase(
		log,
		repository.AnswerTypeRepositoryFactory(log),
		dto.AnswerTypeDTOFactory(log),
	)
}

func (uc *AnswerTypeUseCase) FindAll() ([]*response.AnswerTypeResponse, error) {
	answerTypes, err := uc.Repository.FindAll()
	if err != nil {
		return nil, err
	}

	var response []*response.AnswerTypeResponse
	for _, answerType := range answerTypes {
		atRes := uc.DTO.ConvertEntityToResponse(answerType)
		response = append(response, atRes)
	}

	return response, nil
}

func (uc *AnswerTypeUseCase) FindByID(id string) (*response.AnswerTypeResponse, error) {
	answerType, err := uc.Repository.FindByKeys(map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}

	if answerType == nil {
		return nil, nil
	}

	return uc.DTO.ConvertEntityToResponse(answerType), nil
}
