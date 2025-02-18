package usecase

import (
	"errors"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/dto"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ICoverUseCase interface {
	CreateCover(req *request.CreateCoverRequest) (*response.CoverResponse, error)
	FindByID(id uuid.UUID) (*response.CoverResponse, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.CoverResponse, int64, error)
	UpdateCover(id uuid.UUID, req *request.UpdateCoverRequest) (*response.CoverResponse, error)
	DeleteCover(id uuid.UUID) error
}

type CoverUseCase struct {
	Log        *logrus.Logger
	DTO        dto.ICoverDTO
	Repository repository.ICoverRepository
	Viper      *viper.Viper
}

func NewCoverUseCase(
	log *logrus.Logger,
	dto dto.ICoverDTO,
	repo repository.ICoverRepository,
	viper *viper.Viper,
) ICoverUseCase {
	return &CoverUseCase{
		Log:        log,
		DTO:        dto,
		Repository: repo,
		Viper:      viper,
	}
}

func CoverUseCaseFactory(log *logrus.Logger, viper *viper.Viper) ICoverUseCase {
	dto := dto.CoverDTOFactory(log, viper)
	repo := repository.CoverRepositoryFactory(log)
	return NewCoverUseCase(log, dto, repo, viper)
}

func (uc *CoverUseCase) CreateCover(req *request.CreateCoverRequest) (*response.CoverResponse, error) {
	cover, err := uc.Repository.CreateCoverRepository(&entity.Cover{
		Path: req.Path,
	})
	if err != nil {
		uc.Log.Error("[CoverUseCase.CreateCover] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(cover), nil
}

func (uc *CoverUseCase) FindByID(id uuid.UUID) (*response.CoverResponse, error) {
	cover, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[CoverUseCase.FindByID] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(cover), nil
}

func (uc *CoverUseCase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.CoverResponse, int64, error) {
	entities, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		uc.Log.Error("[CoverUseCase.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	var responses []response.CoverResponse
	for _, entity := range *entities {
		res := uc.DTO.ConvertEntityToResponse(&entity)
		responses = append(responses, *res)
	}

	return &responses, total, nil
}

func (uc *CoverUseCase) UpdateCover(id uuid.UUID, req *request.UpdateCoverRequest) (*response.CoverResponse, error) {
	cover, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[CoverUseCase.UpdateCover] " + err.Error())
		return nil, err
	}
	if cover == nil {
		return nil, errors.New("Cover not found")
	}

	cov, err := uc.Repository.UpdateCoverRepository(&entity.Cover{
		ID:   id,
		Path: req.Path,
	})
	if err != nil {
		uc.Log.Error("[CoverUseCase.UpdateCover] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(cov), nil
}

func (uc *CoverUseCase) DeleteCover(id uuid.UUID) error {
	cover, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[CoverUseCase.DeleteCover] " + err.Error())
		return err
	}
	if cover == nil {
		return errors.New("Cover not found")
	}

	err = uc.Repository.DeleteCoverRepository(id)
	if err != nil {
		uc.Log.Error("[CoverUseCase.DeleteCover] " + err.Error())
		return err
	}

	return nil
}
