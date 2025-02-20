package repository

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IEventEmployeeRepository interface {
	CreateEventEmployee(ent *entity.EventEmployee) (*entity.EventEmployee, error)
	DeleteByEventID(eventID string) error
}

type EventEmployeeRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewEventEmployeeRepository(log *logrus.Logger, db *gorm.DB) *EventEmployeeRepository {
	return &EventEmployeeRepository{
		Log: log,
		DB:  db,
	}
}

func EventEmployeeRepositoryFactory(log *logrus.Logger) IEventEmployeeRepository {
	db := config.NewDatabase()
	return NewEventEmployeeRepository(log, db)
}

func (r *EventEmployeeRepository) CreateEventEmployee(ent *entity.EventEmployee) (*entity.EventEmployee, error) {
	if err := r.DB.Create(ent).Error; err != nil {
		r.Log.Error("[EventEmployeeRepository.CreateEventEmployee] Error when create event employee: ", err)
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[EventEmployeeRepository.CreateEventEmployee] Error when get event employee: ", err)
		return nil, err
	}

	return ent, nil
}

func (r *EventEmployeeRepository) DeleteByEventID(eventID string) error {
	if err := r.DB.Where("event_id = ?", eventID).Delete(&entity.EventEmployee{}).Error; err != nil {
		r.Log.Error("[EventEmployeeRepository.DeleteByEventID] Error when delete event employee by event id: ", err)
		return err
	}

	return nil
}
