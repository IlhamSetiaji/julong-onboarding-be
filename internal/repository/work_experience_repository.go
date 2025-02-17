package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IWorkExperienceRepository interface {
	CreateWorkExperience(ent *entity.WorkExperience) (*entity.WorkExperience, error)
	UpdateWorkExperience(ent *entity.WorkExperience) (*entity.WorkExperience, error)
	DeleteWorkExperience(id uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.WorkExperience, error)
	DeleteByUserProfileID(userProfileID uuid.UUID) error
	DeleteNotInIDAndUserProfileID(ids []uuid.UUID, userProfileID uuid.UUID) error
}

type WorkExperienceRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewWorkExperienceRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *WorkExperienceRepository {
	return &WorkExperienceRepository{
		Log: log,
		DB:  db,
	}
}

func WorkExperienceRepositoryFactory(
	log *logrus.Logger,
) IWorkExperienceRepository {
	db := config.NewDatabase()
	return NewWorkExperienceRepository(log, db)
}

func (r *WorkExperienceRepository) CreateWorkExperience(ent *entity.WorkExperience) (*entity.WorkExperience, error) {
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

func (r *WorkExperienceRepository) UpdateWorkExperience(ent *entity.WorkExperience) (*entity.WorkExperience, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Model(&entity.WorkExperience{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
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

func (r *WorkExperienceRepository) DeleteWorkExperience(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var workExperience entity.WorkExperience
	if err := tx.First(&workExperience, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&workExperience).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *WorkExperienceRepository) FindByID(id uuid.UUID) (*entity.WorkExperience, error) {
	ent := new(entity.WorkExperience)
	if err := r.DB.First(ent, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return ent, nil
}

func (r *WorkExperienceRepository) DeleteByUserProfileID(userProfileID uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("user_profile_id = ?", userProfileID).Delete(&entity.WorkExperience{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *WorkExperienceRepository) DeleteNotInIDAndUserProfileID(ids []uuid.UUID, userProfileID uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("id NOT IN ?", ids).Where("user_profile_id = ?", userProfileID).Delete(&entity.WorkExperience{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
