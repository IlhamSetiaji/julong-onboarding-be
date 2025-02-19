package repository

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IEmployeeTaskAttachmentRepository interface {
	CreateEmployeeTaskAttachment(ent *entity.EmployeeTaskAttachment) (*entity.EmployeeTaskAttachment, error)
	DeleteByTemplateTaskID(id uuid.UUID) error
	DeleteEmployeeTaskAttachment(ent *entity.EmployeeTaskAttachment) error
	FindByID(id uuid.UUID) (*entity.EmployeeTaskAttachment, error)
}

type EmployeeTaskAttachmentRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewEmployeeTaskAttachmentRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *EmployeeTaskAttachmentRepository {
	return &EmployeeTaskAttachmentRepository{
		Log: log,
		DB:  db,
	}
}

func EmployeeTaskAttachmentRepositoryFactory(
	log *logrus.Logger,
) IEmployeeTaskAttachmentRepository {
	db := config.NewDatabase()
	return NewEmployeeTaskAttachmentRepository(log, db)
}

func (r *EmployeeTaskAttachmentRepository) CreateEmployeeTaskAttachment(ent *entity.EmployeeTaskAttachment) (*entity.EmployeeTaskAttachment, error) {
	if err := r.DB.Create(ent).Error; err != nil {
		r.Log.Error("[EmployeeTaskAttachmentRepository.CreateEmployeeTaskAttachment] Error when create employee task attachment: ", err)
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[EmployeeTaskAttachmentRepository.CreateEmployeeTaskAttachment] Error when get employee task attachment: ", err)
		return nil, err
	}

	return ent, nil
}

func (r *EmployeeTaskAttachmentRepository) DeleteByTemplateTaskID(id uuid.UUID) error {
	if err := r.DB.Where("employee_task_id = ?", id).Delete(&entity.EmployeeTaskAttachment{}).Error; err != nil {
		r.Log.Error("[EmployeeTaskAttachmentRepository.DeleteByTemplateTaskID] Error when delete employee task attachment: ", err)
		return err
	}

	return nil
}

func (r *EmployeeTaskAttachmentRepository) DeleteEmployeeTaskAttachment(ent *entity.EmployeeTaskAttachment) error {
	if err := r.DB.Delete(ent).Error; err != nil {
		r.Log.Error("[EmployeeTaskAttachmentRepository.DeleteEmployeeTaskAttachment] Error when delete employee task attachment: ", err)
		return err
	}

	return nil
}

func (r *EmployeeTaskAttachmentRepository) FindByID(id uuid.UUID) (*entity.EmployeeTaskAttachment, error) {
	ent := new(entity.EmployeeTaskAttachment)
	if err := r.DB.Where("id = ?", id).First(ent).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			r.Log.Error("[EmployeeTaskAttachmentRepository.FindByID] Error when get employee task attachment: ", err)
			return nil, err
		}
	}

	return ent, nil
}
