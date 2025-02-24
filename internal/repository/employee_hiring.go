package repository

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IEmployeeHiringRepository interface {
	CreateEmployeeHiring(ent *entity.EmployeeHiring) (*entity.EmployeeHiring, error)
}

type EmployeeHiringRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewEmployeeHiringRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *EmployeeHiringRepository {
	return &EmployeeHiringRepository{
		Log: log,
		DB:  db,
	}
}

func EmployeeHiringRepositoryFactory(
	log *logrus.Logger,
) IEmployeeHiringRepository {
	db := config.NewDatabase()
	return NewEmployeeHiringRepository(log, db)
}

func (r *EmployeeHiringRepository) CreateEmployeeHiring(ent *entity.EmployeeHiring) (*entity.EmployeeHiring, error) {
	if err := r.DB.Create(ent).Error; err != nil {
		r.Log.Error("[EmployeeHiringRepository.CreateEmployeeHiring] Error when create employee hiring: ", err)
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[EmployeeHiringRepository.CreateEmployeeHiring] Error when get employee hiring: ", err)
		return nil, err
	}

	return ent, nil
}
