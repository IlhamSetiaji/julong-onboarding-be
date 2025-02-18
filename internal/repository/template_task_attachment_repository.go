package repository

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ITemplateTaskAttachmentRepository interface {
	CreateTemplateTaskAttachment(ent *entity.TemplateTaskAttachment) (*entity.TemplateTaskAttachment, error)
	DeleteByTemplateTaskID(id uuid.UUID) error
	DeleteTemplateTaskAttachment(ent *entity.TemplateTaskAttachment) error
}

type TemplateTaskAttachmentRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewTemplateTaskAttachmentRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *TemplateTaskAttachmentRepository {
	return &TemplateTaskAttachmentRepository{
		Log: log,
		DB:  db,
	}
}

func TemplateTaskAttachmentRepositoryFactory(
	log *logrus.Logger,
) ITemplateTaskAttachmentRepository {
	db := config.NewDatabase()
	return NewTemplateTaskAttachmentRepository(log, db)
}

func (r *TemplateTaskAttachmentRepository) CreateTemplateTaskAttachment(ent *entity.TemplateTaskAttachment) (*entity.TemplateTaskAttachment, error) {
	if err := r.DB.Create(ent).Error; err != nil {
		r.Log.Error("[TemplateTaskAttachmentRepository.CreateTemplateTaskAttachment] Error when create template task attachment: ", err)
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[TemplateTaskAttachmentRepository.CreateTemplateTaskAttachment] Error when get template task attachment: ", err)
		return nil, err
	}

	return ent, nil
}

func (r *TemplateTaskAttachmentRepository) DeleteByTemplateTaskID(id uuid.UUID) error {
	if err := r.DB.Where("template_task_id = ?", id).Delete(&entity.TemplateTaskAttachment{}).Error; err != nil {
		r.Log.Error("[TemplateTaskAttachmentRepository.DeleteByTemplateTaskID] Error when delete template task attachment: ", err)
		return err
	}

	return nil
}

func (r *TemplateTaskAttachmentRepository) DeleteTemplateTaskAttachment(ent *entity.TemplateTaskAttachment) error {
	if err := r.DB.Delete(ent).Error; err != nil {
		r.Log.Error("[TemplateTaskAttachmentRepository.DeleteTemplateTaskAttachment] Error when delete template task attachment: ", err)
		return err
	}

	return nil
}
