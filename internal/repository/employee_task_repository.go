package repository

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IEmployeeTaskRepository interface {
	CreateEmployeeTask(ent *entity.EmployeeTask) (*entity.EmployeeTask, error)
	UpdateEmployeeTask(ent *entity.EmployeeTask) (*entity.EmployeeTask, error)
	DeleteEmployeeTask(ent *entity.EmployeeTask) error
	FindByID(id uuid.UUID) (*entity.EmployeeTask, error)
	FindAllByEmployeeID(employeeID uuid.UUID) (*[]entity.EmployeeTask, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.EmployeeTask, int64, error)
	CountByKeys(keys map[string]interface{}) (int64, error)
	FindAllByEmployeeIDAndKanbanPaginated(employeeID uuid.UUID, kanban entity.EmployeeTaskKanbanEnum, page, pageSize int, search string, sort map[string]interface{}) (*[]entity.EmployeeTask, int64, error)
	FindByKeys(keys map[string]interface{}) (*entity.EmployeeTask, error)
}

type EmployeeTaskRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewEmployeeTaskRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *EmployeeTaskRepository {
	return &EmployeeTaskRepository{
		Log: log,
		DB:  db,
	}
}

func EmployeeTaskRepositoryFactory(
	log *logrus.Logger,
) IEmployeeTaskRepository {
	db := config.NewDatabase()
	return NewEmployeeTaskRepository(log, db)
}

func (r *EmployeeTaskRepository) CreateEmployeeTask(ent *entity.EmployeeTask) (*entity.EmployeeTask, error) {
	if err := r.DB.Create(ent).Error; err != nil {
		r.Log.Error("[EmployeeTaskRepository.CreateEmployeeTask] Error when create employee task: ", err)
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[EmployeeTaskRepository.CreateEmployeeTask] Error when get employee task: ", err)
		return nil, err
	}

	return ent, nil
}

func (r *EmployeeTaskRepository) UpdateEmployeeTask(ent *entity.EmployeeTask) (*entity.EmployeeTask, error) {
	if err := r.DB.Model(&entity.EmployeeTask{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		r.Log.Error("[EmployeeTaskRepository.UpdateEmployeeTask] Error when update employee task: ", err)
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[EmployeeTaskRepository.UpdateEmployeeTask] Error when get employee task: ", err)
		return nil, err
	}

	return ent, nil
}

func (r *EmployeeTaskRepository) DeleteEmployeeTask(ent *entity.EmployeeTask) error {
	if err := r.DB.Delete(ent).Error; err != nil {
		r.Log.Error("[EmployeeTaskRepository.DeleteEmployeeTask] Error when delete employee task: ", err)
		return err
	}

	return nil
}

func (r *EmployeeTaskRepository) FindByID(id uuid.UUID) (*entity.EmployeeTask, error) {
	var ent entity.EmployeeTask
	if err := r.DB.Preload("EmployeeTaskAttachments").Preload("EmployeeTaskChecklists").First(&ent, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			r.Log.Error("[EmployeeTaskRepository.FindByID] Error when get employee task by id: ", err)
			return nil, err
		}
	}

	return &ent, nil
}

func (r *EmployeeTaskRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.EmployeeTask, int64, error) {
	var employeeTasks []entity.EmployeeTask
	var total int64

	query := r.DB.Preload("EmployeeTaskAttachments").Preload("EmployeeTaskChecklists").Where("name LIKE ?", "%"+search+"%")
	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&employeeTasks).Error; err != nil {
		r.Log.Error("[EmployeeTaskRepository.FindAllPaginated] Error when get template tasks: ", err)
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Error("[EmployeeTaskRepository.FindAllPaginated] Error when count template tasks: ", err)
		return nil, 0, err
	}

	return &employeeTasks, total, nil
}

func (r *EmployeeTaskRepository) CountByKeys(keys map[string]interface{}) (int64, error) {
	var total int64

	if err := r.DB.Model(&entity.EmployeeTask{}).Where(keys).Count(&total).Error; err != nil {
		r.Log.Error("[EmployeeTaskRepository.CountByKeys] Error when count employee tasks: ", err)
		return 0, err
	}

	return total, nil
}

func (r *EmployeeTaskRepository) FindAllByEmployeeID(employeeID uuid.UUID) (*[]entity.EmployeeTask, error) {
	var employeeTasks []entity.EmployeeTask

	if err := r.DB.Preload("EmployeeTaskAttachments").Preload("EmployeeTaskChecklists").Where("employee_id = ?", employeeID).Find(&employeeTasks).Error; err != nil {
		r.Log.Error("[EmployeeTaskRepository.FindAllByEmployeeID] Error when get employee tasks by employee id: ", err)
		return nil, err
	}

	return &employeeTasks, nil
}

func (r *EmployeeTaskRepository) FindAllByEmployeeIDAndKanbanPaginated(employeeID uuid.UUID, kanban entity.EmployeeTaskKanbanEnum, page, pageSize int, search string, sort map[string]interface{}) (*[]entity.EmployeeTask, int64, error) {
	var employeeTasks []entity.EmployeeTask
	var total int64

	query := r.DB.Preload("EmployeeTaskAttachments").Preload("EmployeeTaskChecklists").Where("employee_id = ?", employeeID).Where("kanban = ?", kanban)
	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&employeeTasks).Error; err != nil {
		r.Log.Error("[EmployeeTaskRepository.FindAllByEmployeeIDAndKanbanPaginated] Error when get employee tasks by employee id and kanban: ", err)
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Error("[EmployeeTaskRepository.FindAllByEmployeeIDAndKanbanPaginated] Error when count employee tasks by employee id and kanban: ", err)
		return nil, 0, err
	}

	return &employeeTasks, total, nil
}

func (r *EmployeeTaskRepository) FindByKeys(keys map[string]interface{}) (*entity.EmployeeTask, error) {
	var ent entity.EmployeeTask

	if err := r.DB.Preload("EmployeeTaskAttachments").Preload("EmployeeTaskChecklists").Where(keys).First(&ent).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			r.Log.Error("[EmployeeTaskRepository.FindByKeys] Error when get employee task by keys: ", err)
			return nil, err
		}
	}

	return &ent, nil
}
