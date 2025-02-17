package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ISkillRepository interface {
	CreateSkill(ent *entity.Skill) (*entity.Skill, error)
	UpdateSkill(ent *entity.Skill) (*entity.Skill, error)
	DeleteSkill(id uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.Skill, error)
	DeleteByUserProfileID(userProfileID uuid.UUID) error
	DeleteNotInIDAndUserProfileID(ids []uuid.UUID, userProfileID uuid.UUID) error
}

type SkillRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewSkillRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *SkillRepository {
	return &SkillRepository{
		Log: log,
		DB:  db,
	}
}

func SkillRepositoryFactory(
	log *logrus.Logger,
) ISkillRepository {
	db := config.NewDatabase()
	return NewSkillRepository(log, db)
}

func (r *SkillRepository) CreateSkill(ent *entity.Skill) (*entity.Skill, error) {
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

func (r *SkillRepository) UpdateSkill(ent *entity.Skill) (*entity.Skill, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Model(&entity.Skill{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
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

func (r *SkillRepository) DeleteSkill(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var skill entity.Skill
	if err := tx.First(&skill, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&skill).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *SkillRepository) FindByID(id uuid.UUID) (*entity.Skill, error) {
	ent := new(entity.Skill)
	if err := r.DB.First(ent, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return ent, nil
}

func (r *SkillRepository) DeleteByUserProfileID(userProfileID uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("user_profile_id = ?", userProfileID).Delete(&entity.Skill{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *SkillRepository) DeleteNotInIDAndUserProfileID(ids []uuid.UUID, userProfileID uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("user_profile_id = ? AND id NOT IN ?", userProfileID, ids).Delete(&entity.Skill{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
