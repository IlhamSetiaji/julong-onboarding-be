package usecase

import (
	"errors"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/dto"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IEmployeeTaskAttachmentUseCase interface {
	FindByID(id uuid.UUID) (*response.EmployeeTaskAttachmentResponse, error)
	DeleteEmployeeTaskAttachment(id uuid.UUID) error
}

type EmployeeTaskAttachmentUseCase struct {
	Log        *logrus.Logger
	Repository repository.IEmployeeTaskAttachmentRepository
	DTO        dto.IEmployeeTaskAttachmentDTO
	Viper      *viper.Viper
}

func NewEmployeeTaskAttachmentUseCase(
	log *logrus.Logger,
	repo repository.IEmployeeTaskAttachmentRepository,
	dto dto.IEmployeeTaskAttachmentDTO,
	viper *viper.Viper,
) IEmployeeTaskAttachmentUseCase {
	return &EmployeeTaskAttachmentUseCase{
		Log:        log,
		Repository: repo,
		DTO:        dto,
		Viper:      viper,
	}
}

func EmployeeTaskAttachmentUseCaseFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IEmployeeTaskAttachmentUseCase {
	repo := repository.EmployeeTaskAttachmentRepositoryFactory(log)
	etaDTO := dto.EmployeeTaskAttachmentDTOFactory(log, viper)
	return NewEmployeeTaskAttachmentUseCase(log, repo, etaDTO, viper)
}

func (uc *EmployeeTaskAttachmentUseCase) FindByID(id uuid.UUID) (*response.EmployeeTaskAttachmentResponse, error) {
	attachment, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[EmployeeTaskAttachmentUseCase.FindByID] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(attachment), nil
}

func (uc *EmployeeTaskAttachmentUseCase) DeleteEmployeeTaskAttachment(id uuid.UUID) error {
	ent, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[EmployeeTaskAttachmentUseCase.DeleteEmployeeTaskAttachment] " + err.Error())
		return err
	}
	if ent == nil {
		return errors.New("Employee Task Attachment not found")
	}

	err = uc.Repository.DeleteEmployeeTaskAttachment(ent)
	if err != nil {
		uc.Log.Error("[EmployeeTaskAttachmentUseCase.DeleteEmployeeTaskAttachment] " + err.Error())
		return err
	}

	return nil
}
