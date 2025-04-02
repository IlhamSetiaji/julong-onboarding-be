package usecase

import (
	"fmt"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/dto"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ISurveyTemplateUseCase interface {
	CreateSurveyTemplate(req *request.CreateSurveyTemplateRequest) (*response.SurveyTemplateResponse, error)
	UpdateSurveyTemplate(req *request.UpdateSurveyTemplateRequest) (*response.SurveyTemplateResponse, error)
	DeleteSurveyTemplate(id string) error
	FindSurveyTemplateByID(id string) (*response.SurveyTemplateResponse, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.SurveyTemplateResponse, int64, error)
}

type SurveyTemplateUseCase struct {
	Log                      *logrus.Logger
	Viper                    *viper.Viper
	SurveyTemplateRepository repository.ISurveyTemplateRepository
	SurveyTemplateDTO        dto.ISurveyTemplateDTO
}

func NewSurveyTemplateUseCase(
	log *logrus.Logger,
	viper *viper.Viper,
	surveyTemplateRepository repository.ISurveyTemplateRepository,
	surveyTemplateDTO dto.ISurveyTemplateDTO,
) *SurveyTemplateUseCase {
	return &SurveyTemplateUseCase{
		Log:                      log,
		Viper:                    viper,
		SurveyTemplateRepository: surveyTemplateRepository,
		SurveyTemplateDTO:        surveyTemplateDTO,
	}
}

func SurveyTemplateUseCaseFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ISurveyTemplateUseCase {
	surveyTemplateRepository := repository.SurveyTemplateRepositoryFactory(log)
	surveyTemplateDTO := dto.SurveyTemplateDTOFactory(log, viper)
	return NewSurveyTemplateUseCase(log, viper, surveyTemplateRepository, surveyTemplateDTO)
}

func (u *SurveyTemplateUseCase) CreateSurveyTemplate(req *request.CreateSurveyTemplateRequest) (*response.SurveyTemplateResponse, error) {
	surveyNumber, err := u.generateRandomSurveyNumber()
	if err != nil {
		u.Log.Error("[SurveyTemplateUseCase.CreateSurveyTemplate] Error when generating random survey number: ", err)
		return nil, err
	}

	ent, err := u.SurveyTemplateRepository.CreateSurveyTemplate(&entity.SurveyTemplate{
		Title:        req.Title,
		SurveyNumber: *surveyNumber,
		Status:       entity.SURVEY_TEMPLATE_STATUS_ENUM_DRAFT,
	})

	if err != nil {
		u.Log.Error("[SurveyTemplateUseCase.CreateSurveyTemplate] Error when creating survey template: ", err)
		return nil, err
	}

	resp := u.SurveyTemplateDTO.ConvertEntityToResponse(ent)
	return resp, nil
}

func (u *SurveyTemplateUseCase) UpdateSurveyTemplate(req *request.UpdateSurveyTemplateRequest) (*response.SurveyTemplateResponse, error) {
	ent, err := u.SurveyTemplateRepository.FindByKeys(map[string]interface{}{
		"id": req.ID,
	})
	if err != nil {
		u.Log.Error("[SurveyTemplateUseCase.UpdateSurveyTemplate] Error when finding survey template: ", err)
		return nil, err
	}

	if ent == nil {
		return nil, fmt.Errorf("survey template not found")
	}

	ent.Title = req.Title
	ent.Status = entity.SurveyTemplateStatusEnum(req.Status)

	ent, err = u.SurveyTemplateRepository.UpdateSurveyTemplate(ent)
	if err != nil {
		u.Log.Error("[SurveyTemplateUseCase.UpdateSurveyTemplate] Error when updating survey template: ", err)
		return nil, err
	}

	resp := u.SurveyTemplateDTO.ConvertEntityToResponse(ent)
	return resp, nil
}

func (u *SurveyTemplateUseCase) DeleteSurveyTemplate(id string) error {
	ent, err := u.SurveyTemplateRepository.FindByKeys(map[string]interface{}{
		"id": id,
	})
	if err != nil {
		u.Log.Error("[SurveyTemplateUseCase.DeleteSurveyTemplate] Error when finding survey template: ", err)
		return err
	}

	if ent == nil {
		return fmt.Errorf("survey template not found")
	}

	err = u.SurveyTemplateRepository.DeleteSurveyTemplate(ent)
	if err != nil {
		u.Log.Error("[SurveyTemplateUseCase.DeleteSurveyTemplate] Error when deleting survey template: ", err)
		return err
	}

	return nil
}

func (u *SurveyTemplateUseCase) FindSurveyTemplateByID(id string) (*response.SurveyTemplateResponse, error) {
	ent, err := u.SurveyTemplateRepository.FindByKeys(map[string]interface{}{
		"id": id,
	})
	if err != nil {
		u.Log.Error("[SurveyTemplateUseCase.FindSurveyTemplateByID] Error when finding survey template: ", err)
		return nil, err
	}

	if ent == nil {
		return nil, fmt.Errorf("survey template not found")
	}

	resp := u.SurveyTemplateDTO.ConvertEntityToResponse(ent)
	return resp, nil
}

func (u *SurveyTemplateUseCase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.SurveyTemplateResponse, int64, error) {
	ent, total, err := u.SurveyTemplateRepository.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		u.Log.Error("[SurveyTemplateUseCase.FindAllPaginated] Error when finding survey template: ", err)
		return nil, 0, err
	}

	var resp []response.SurveyTemplateResponse
	for _, v := range *ent {
		resp = append(resp, *u.SurveyTemplateDTO.ConvertEntityToResponse(&v))
	}

	return &resp, total, nil
}

func (u *SurveyTemplateUseCase) generateRandomSurveyNumber() (*string, error) {
	// Find the latest survey number
	latestSurvey, err := u.SurveyTemplateRepository.FindLatestSurveyNumber()
	if err != nil {
		u.Log.Error("[SurveyTemplateUseCase.generateRandomSurveyNumber] Error when finding latest survey number: ", err)
		return nil, err
	}

	// Generate the next survey number
	var nextNumber int
	if latestSurvey != nil && latestSurvey.SurveyNumber != "" {
		// Extract the numeric part of the survey number
		var currentNumber int
		_, err := fmt.Sscanf(latestSurvey.SurveyNumber, "SURVEY-%05d", &currentNumber)
		if err != nil {
			u.Log.Error("[SurveyTemplateUseCase.generateRandomSurveyNumber] Error parsing survey number: ", err)
			return nil, err
		}
		nextNumber = currentNumber + 1
	} else {
		// Start from 1 if no survey number exists
		nextNumber = 1
	}

	// Format the next survey number
	newSurveyNumber := fmt.Sprintf("SURVEY-%05d", nextNumber)

	// Check if the generated survey number already exists
	exist, err := u.SurveyTemplateRepository.FindByKeys(map[string]interface{}{
		"survey_number": newSurveyNumber,
	})
	if err != nil {
		u.Log.Error("[SurveyTemplateUseCase.generateRandomSurveyNumber] Error when finding survey number: ", err)
		return nil, err
	}
	if exist != nil {
		// Retry generating a new survey number if it already exists
		return u.generateRandomSurveyNumber()
	}

	return &newSurveyNumber, nil
}
