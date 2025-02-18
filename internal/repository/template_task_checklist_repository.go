package repository

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ITemplateTaskChecklistRepository interface {
	CreateTaskChecklistRepository(ent *entity.TemplateTaskChecklist) (*entity.TemplateTaskChecklist, error)
	UpdateTaskChecklistRepository(ent *entity.TemplateTaskChecklist) (*entity.TemplateTaskChecklist, error)
	DeleteByTemplateTaskID(id string) error
	FindByKeys(keys map[string]interface{}) (*entity.TemplateTaskChecklist, error)
}

type TemplateTaskChecklistRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewTemplateTaskChecklistRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *TemplateTaskChecklistRepository {
	return &TemplateTaskChecklistRepository{
		Log: log,
		DB:  db,
	}
}

func TemplateTaskChecklistRepositoryFactory(
	log *logrus.Logger,
) ITemplateTaskChecklistRepository {
	db := config.NewDatabase()
	return NewTemplateTaskChecklistRepository(log, db)
}

func (r *TemplateTaskChecklistRepository) CreateTaskChecklistRepository(ent *entity.TemplateTaskChecklist) (*entity.TemplateTaskChecklist, error) {
	if err := r.DB.Create(ent).Error; err != nil {
		r.Log.Error("[TemplateTaskChecklistRepository.CreateTaskChecklistRepository] Error when create template task checklist: ", err)
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[TemplateTaskChecklistRepository.CreateTaskChecklistRepository] Error when get template task checklist: ", err)
		return nil, err
	}

	return ent, nil
}

func (r *TemplateTaskChecklistRepository) UpdateTaskChecklistRepository(ent *entity.TemplateTaskChecklist) (*entity.TemplateTaskChecklist, error) {
	if err := r.DB.Model(&entity.TemplateTaskChecklist{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		r.Log.Error("[TemplateTaskChecklistRepository.UpdateTaskChecklistRepository] Error when update template task checklist: ", err)
		return nil, err
	}

	return ent, nil
}

func (r *TemplateTaskChecklistRepository) DeleteByTemplateTaskID(id string) error {
	if err := r.DB.Where("template_task_id = ?", id).Delete(&entity.TemplateTaskChecklist{}).Error; err != nil {
		r.Log.Error("[TemplateTaskChecklistRepository.DeleteByTemplateTaskID] Error when delete template task checklist: ", err)
		return err
	}

	return nil
}

func (r *TemplateTaskChecklistRepository) FindByKeys(keys map[string]interface{}) (*entity.TemplateTaskChecklist, error) {
	var taskChecklist entity.TemplateTaskChecklist
	if err := r.DB.Where(keys).First(&taskChecklist).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			r.Log.Error("[TemplateTaskChecklistRepository.FindByKeys] Error when get template task checklist: ", err)
			return nil, err
		}
	}

	return &taskChecklist, nil
}
