package repository

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IEventRepository interface {
	CreateEvent(event *entity.Event) (*entity.Event, error)
	UpdateEvent(event *entity.Event) (*entity.Event, error)
	DeleteEvent(event *entity.Event) error
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) ([]entity.Event, int64, error)
	FindByID(id uuid.UUID) (*entity.Event, error)
}

type EventRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewEventRepository(log *logrus.Logger, db *gorm.DB) *EventRepository {
	return &EventRepository{
		Log: log,
		DB:  db,
	}
}

func EventRepositoryFactory(log *logrus.Logger) IEventRepository {
	db := config.NewDatabase()
	return NewEventRepository(log, db)
}

func (r *EventRepository) CreateEvent(event *entity.Event) (*entity.Event, error) {
	if err := r.DB.Create(event).Error; err != nil {
		r.Log.Error("[EventRepository.CreateEvent] Error when create event: ", err)
		return nil, err
	}

	if err := r.DB.Preload("TemplateTask.TemplateTaskAttachments").
		Preload("TemplateTask.TemplateTaskChecklists").
		Preload("EventEmployees").First(event, event.ID).Error; err != nil {
		r.Log.Error("[EventRepository.CreateEvent] Error when get event: ", err)
		return nil, err
	}

	return event, nil
}

func (r *EventRepository) UpdateEvent(event *entity.Event) (*entity.Event, error) {
	if err := r.DB.Model(&entity.Event{}).Where("id = ?", event.ID).Updates(event).Error; err != nil {
		r.Log.Error("[EventRepository.UpdateEvent] Error when update event: ", err)
		return nil, err
	}

	if err := r.DB.Preload("TemplateTask.TemplateTaskAttachments").
		Preload("TemplateTask.TemplateTaskChecklists").
		Preload("EventEmployees").First(event, event.ID).Error; err != nil {
		r.Log.Error("[EventRepository.UpdateEvent] Error when get event: ", err)
		return nil, err
	}

	return event, nil
}

func (r *EventRepository) DeleteEvent(event *entity.Event) error {
	if err := r.DB.Delete(event).Error; err != nil {
		r.Log.Error("[EventRepository.DeleteEvent] Error when delete event: ", err)
		return err
	}

	return nil
}

func (r *EventRepository) FindByID(id uuid.UUID) (*entity.Event, error) {
	var event entity.Event
	if err := r.DB.Preload("TemplateTask.TemplateTaskAttachments").
		Preload("TemplateTask.TemplateTaskChecklists").
		Preload("EventEmployees").First(&event, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			r.Log.Error("[EventRepository.FindByID] Error when get event by id: ", err)
			return nil, err
		}
	}

	return &event, nil
}

func (r *EventRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) ([]entity.Event, int64, error) {
	var events []entity.Event
	var total int64

	query := r.DB.Model(&events).Preload("TemplateTask.TemplateTaskAttachments").
		Preload("TemplateTask.TemplateTaskChecklists").
		Preload("EventEmployees")

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Error("[EventRepository.FindAllPaginated] Error when count event: ", err)
		return nil, 0, err
	}

	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Preload("EventAttachments").Find(&events).Error; err != nil {
		r.Log.Error("[EventRepository.FindAllPaginated] Error when get event: ", err)
		return nil, 0, err
	}

	return events, total, nil
}
