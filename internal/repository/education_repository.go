package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IEducationRepository interface {
	CreateEducation(ent *entity.Education) (*entity.Education, error)
	UpdateEducation(ent *entity.Education) (*entity.Education, error)
	DeleteEducation(id uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.Education, error)
	DeleteByUserProfileID(userProfileID uuid.UUID) error
	DeleteNotInIDAndUserProfileID(ids []uuid.UUID, userProfileID uuid.UUID) error
}

type EducationRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewEducationRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *EducationRepository {
	return &EducationRepository{
		Log: log,
		DB:  db,
	}
}

func EducationRepositoryFactory(
	log *logrus.Logger,
) IEducationRepository {
	db := config.NewDatabase()
	return NewEducationRepository(log, db)
}

func (r *EducationRepository) CreateEducation(ent *entity.Education) (*entity.Education, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *EducationRepository) UpdateEducation(ent *entity.Education) (*entity.Education, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Model(&entity.Education{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *EducationRepository) DeleteEducation(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var education entity.Education
	if err := tx.First(&education, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&education).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *EducationRepository) FindByID(id uuid.UUID) (*entity.Education, error) {
	ent := new(entity.Education)
	if err := r.DB.First(ent, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return ent, nil
}

func (r *EducationRepository) DeleteByUserProfileID(userProfileID uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("user_profile_id = ?", userProfileID).Delete(&entity.Education{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *EducationRepository) DeleteNotInIDAndUserProfileID(ids []uuid.UUID, userProfileID uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("user_profile_id = ? AND id NOT IN ?", userProfileID, ids).Delete(&entity.Education{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
