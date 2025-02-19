package repository

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IEmployeeTaskChecklistRepository interface {
	CreateEmployeeTaskChecklist(ent *entity.EmployeeTaskChecklist) (*entity.EmployeeTaskChecklist, error)
	UpdateEmployeeTaskChecklist(ent *entity.EmployeeTaskChecklist) (*entity.EmployeeTaskChecklist, error)
	DeleteByEmployeeTaskID(id uuid.UUID) error
	FindByKeys(keys map[string]interface{}) (*entity.EmployeeTaskChecklist, error)
	DeleteByEmployeeTaskIDAndNotInChecklistIDs(employeeTaskID uuid.UUID, checklistIDs []uuid.UUID) error
}

type EmployeeTaskChecklistRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewEmployeeTaskChecklistRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *EmployeeTaskChecklistRepository {
	return &EmployeeTaskChecklistRepository{
		Log: log,
		DB:  db,
	}
}

func EmployeeTaskChecklistRepositoryFactory(
	log *logrus.Logger,
) IEmployeeTaskChecklistRepository {
	db := config.NewDatabase()
	return NewEmployeeTaskChecklistRepository(log, db)
}

func (r *EmployeeTaskChecklistRepository) CreateEmployeeTaskChecklist(ent *entity.EmployeeTaskChecklist) (*entity.EmployeeTaskChecklist, error) {
	if err := r.DB.Create(ent).Error; err != nil {
		r.Log.Error("[EmployeeTaskChecklistRepository.CreateEmployeeTaskChecklist] Error when create employee task checklist: ", err)
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[EmployeeTaskChecklistRepository.CreateEmployeeTaskChecklist] Error when get employee task checklist: ", err)
		return nil, err
	}

	return ent, nil
}

func (r *EmployeeTaskChecklistRepository) UpdateEmployeeTaskChecklist(ent *entity.EmployeeTaskChecklist) (*entity.EmployeeTaskChecklist, error) {
	if err := r.DB.Model(&entity.EmployeeTaskChecklist{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		r.Log.Error("[EmployeeTaskChecklistRepository.UpdateEmployeeTaskChecklist] Error when update employee task checklist: ", err)
		return nil, err
	}

	return ent, nil
}

func (r *EmployeeTaskChecklistRepository) DeleteByEmployeeTaskID(id uuid.UUID) error {
	if err := r.DB.Where("employee_task_id = ?", id).Delete(&entity.EmployeeTaskChecklist{}).Error; err != nil {
		r.Log.Error("[EmployeeTaskChecklistRepository.DeleteByEmployeeTaskID] Error when delete employee task checklist: ", err)
		return err
	}

	return nil
}

func (r *EmployeeTaskChecklistRepository) FindByKeys(keys map[string]interface{}) (*entity.EmployeeTaskChecklist, error) {
	var ent entity.EmployeeTaskChecklist
	if err := r.DB.Where(keys).First(&ent).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			r.Log.Error("[EmployeeTaskChecklistRepository.FindByKeys] Error when find employee task checklist: ", err)
			return nil, err
		}
	}

	return &ent, nil
}

func (r *EmployeeTaskChecklistRepository) DeleteByEmployeeTaskIDAndNotInChecklistIDs(employeeTaskID uuid.UUID, checklistIDs []uuid.UUID) error {
	if err := r.DB.Where("employee_task_id = ? AND id NOT IN ?", employeeTaskID, checklistIDs).Delete(&entity.EmployeeTaskChecklist{}).Error; err != nil {
		r.Log.Error("[EmployeeTaskChecklistRepository.DeleteByEmployeeTaskIDAndNotInChecklistIDs] Error when delete employee task checklist: ", err)
		return err
	}

	return nil
}
