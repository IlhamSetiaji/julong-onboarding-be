package usecase

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/dto"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ITemplateTaskAttachmentUseCase interface {
	FindByID(id uuid.UUID) (*response.TemplateTaskAttachmentResponse, error)
	DeleteTemplateTaskAttachment(id uuid.UUID) error
}

type TemplateTaskAttachmentUseCase struct {
	Log                       *logrus.Logger
	Repository                repository.ITemplateTaskAttachmentRepository
	TemplateTaskAttachmentDTO dto.ITemplateTaskAttachmentDTO
	Viper                     *viper.Viper
}

func NewTemplateTaskAttachmentUseCase(
	log *logrus.Logger,
	repo repository.ITemplateTaskAttachmentRepository,
	dto dto.ITemplateTaskAttachmentDTO,
	viper *viper.Viper,
) ITemplateTaskAttachmentUseCase {
	return &TemplateTaskAttachmentUseCase{
		Log:                       log,
		Repository:                repo,
		TemplateTaskAttachmentDTO: dto,
		Viper:                     viper,
	}
}

func TemplateTaskAttachmentUseCaseFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ITemplateTaskAttachmentUseCase {
	repo := repository.TemplateTaskAttachmentRepositoryFactory(log)
	dto := dto.TemplateTaskAttachmentDTOFactory(log, viper)
	return NewTemplateTaskAttachmentUseCase(log, repo, dto, viper)
}

func (uc *TemplateTaskAttachmentUseCase) FindByID(id uuid.UUID) (*response.TemplateTaskAttachmentResponse, error) {
	attachment, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[TemplateTaskAttachmentUseCase.FindByID] " + err.Error())
		return nil, err
	}

	return uc.TemplateTaskAttachmentDTO.ConvertEntityToResponse(attachment), nil
}

func (uc *TemplateTaskAttachmentUseCase) DeleteTemplateTaskAttachment(id uuid.UUID) error {
	err := uc.Repository.DeleteByTemplateTaskID(id)
	if err != nil {
		uc.Log.Error("[TemplateTaskAttachmentUseCase.DeleteTemplateTaskAttachment] " + err.Error())
		return err
	}

	return nil
}
