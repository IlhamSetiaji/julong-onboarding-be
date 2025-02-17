package repository

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IUniversityRepository interface {
	FindAll() ([]*entity.University, error)
}

type UniversityRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewUniversityRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *UniversityRepository {
	return &UniversityRepository{
		Log: log,
		DB:  db,
	}
}

func UniversityRepositoryFactory(
	log *logrus.Logger,
) IUniversityRepository {
	db := config.NewDatabase()
	return NewUniversityRepository(log, db)
}

func (r *UniversityRepository) FindAll() ([]*entity.University, error) {
	var entities []*entity.University
	if err := r.DB.Where("alpha_two_code = ?", "ID").Find(&entities).Error; err != nil {
		r.Log.Error("[UniversityRepository.FindAll] " + err.Error())
		return nil, err
	}

	return entities, nil
}
