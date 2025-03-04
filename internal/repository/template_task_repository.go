package repository

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ITemplateTaskRepository interface {
	CreateTemplateTask(ent *entity.TemplateTask) (*entity.TemplateTask, error)
	UpdateTemplateTask(ent *entity.TemplateTask) (*entity.TemplateTask, error)
	DeleteTemplateTask(ent *entity.TemplateTask) error
	FindByID(id uuid.UUID) (*entity.TemplateTask, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, status entity.TemplateTaskStatusEnum) (*[]entity.TemplateTask, int64, error)
	FindAll() (*[]entity.TemplateTask, error)
	CountKanbanProgressByEmployeeID(employeeID uuid.UUID, kanban entity.EmployeeTaskKanbanEnum) (int, error)
	FindAllByKeys(keys map[string]interface{}) (*[]entity.TemplateTask, error)
}

type TemplateTaskRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewTemplateTaskRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *TemplateTaskRepository {
	return &TemplateTaskRepository{
		Log: log,
		DB:  db,
	}
}

func TemplateTaskRepositoryFactory(
	log *logrus.Logger,
) ITemplateTaskRepository {
	db := config.NewDatabase()
	return NewTemplateTaskRepository(log, db)
}

func (r *TemplateTaskRepository) CreateTemplateTask(ent *entity.TemplateTask) (*entity.TemplateTask, error) {
	if err := r.DB.Create(ent).Error; err != nil {
		r.Log.Error("[TemplateTaskRepository.CreateTemplateTask] Error when create template task: ", err)
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[TemplateTaskRepository.CreateTemplateTask] Error when get template task: ", err)
		return nil, err
	}

	return ent, nil
}

func (r *TemplateTaskRepository) UpdateTemplateTask(ent *entity.TemplateTask) (*entity.TemplateTask, error) {
	if err := r.DB.Model(&entity.TemplateTask{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		r.Log.Error("[TemplateTaskRepository.UpdateTemplateTask] Error when update template task: ", err)
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[TemplateTaskRepository.UpdateTemplateTask] Error when get template task: ", err)
		return nil, err
	}

	return ent, nil
}

func (r *TemplateTaskRepository) DeleteTemplateTask(ent *entity.TemplateTask) error {
	if err := r.DB.Delete(ent).Error; err != nil {
		r.Log.Error("[TemplateTaskRepository.DeleteTemplateTask] Error when delete template task: ", err)
		return err
	}

	return nil
}

func (r *TemplateTaskRepository) FindByID(id uuid.UUID) (*entity.TemplateTask, error) {
	var templateTask entity.TemplateTask
	if err := r.DB.Preload("TemplateTaskAttachments").Preload("TemplateTaskChecklists").First(&templateTask, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			r.Log.Error("[TemplateTaskRepository.FindByID] Error when get template task: ", err)
			return nil, err
		}
	}

	return &templateTask, nil
}

func (r *TemplateTaskRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, status entity.TemplateTaskStatusEnum) (*[]entity.TemplateTask, int64, error) {
	var templateTasks []entity.TemplateTask
	var total int64

	query := r.DB.Preload("TemplateTaskAttachments").Preload("TemplateTaskChecklists").Where("name LIKE ?", "%"+search+"%")
	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&templateTasks).Error; err != nil {
		r.Log.Error("[TemplateTaskRepository.FindAllPaginated] Error when get template tasks: ", err)
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Error("[TemplateTaskRepository.FindAllPaginated] Error when count template tasks: ", err)
		return nil, 0, err
	}

	return &templateTasks, total, nil
}

func (r *TemplateTaskRepository) FindAll() (*[]entity.TemplateTask, error) {
	var templateTasks []entity.TemplateTask

	if err := r.DB.Where("status = ?", entity.TEMPLATE_TASK_STATUS_ENUM_ACTIVE).Find(&templateTasks).Error; err != nil {
		r.Log.Error("[TemplateTaskRepository.FindAll] Error when get template tasks: ", err)
		return nil, err
	}

	return &templateTasks, nil
}

func (r *TemplateTaskRepository) CountKanbanProgressByEmployeeID(employeeID uuid.UUID, kanban entity.EmployeeTaskKanbanEnum) (int, error) {
	var total int64
	if err := r.DB.Model(&entity.EmployeeTask{}).Where("employee_id = ? AND kanban = ?", employeeID, kanban).Count(&total).Error; err != nil {
		r.Log.Error("[TemplateTaskRepository.CountKanbanProgressByEmployeeID] Error when count kanban progress by employee id: ", err)
		return 0, err
	}

	return int(total), nil
}

func (r *TemplateTaskRepository) FindAllByKeys(keys map[string]interface{}) (*[]entity.TemplateTask, error) {
	var templateTasks []entity.TemplateTask

	if err := r.DB.Where(keys).Find(&templateTasks).Error; err != nil {
		r.Log.Error("[TemplateTaskRepository.FindAllByKeys] Error when get template tasks by keys: ", err)
		return nil, err
	}

	return &templateTasks, nil
}
