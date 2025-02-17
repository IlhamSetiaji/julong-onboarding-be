package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IUserProfileRepository interface {
	CreateUserProfile(ent *entity.UserProfile) (*entity.UserProfile, error)
	UpdateUserProfile(ent *entity.UserProfile) (*entity.UserProfile, error)
	FindByID(id uuid.UUID) (*entity.UserProfile, error)
	FindByUserID(userID uuid.UUID) (*entity.UserProfile, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]entity.UserProfile, int64, error)
	DeleteUserProfile(id uuid.UUID) error
}

type UserProfileRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewUserProfileRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *UserProfileRepository {
	return &UserProfileRepository{
		Log: log,
		DB:  db,
	}
}

func UserProfileRepositoryFactory(
	log *logrus.Logger,
) IUserProfileRepository {
	db := config.NewDatabase()
	return NewUserProfileRepository(log, db)
}

func (r *UserProfileRepository) CreateUserProfile(ent *entity.UserProfile) (*entity.UserProfile, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[UserProfileRepository.CreateUserProfile] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserProfileRepository.CreateUserProfile] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserProfileRepository.CreateUserProfile] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("Applicants").Preload("WorkExperiences").Preload("Educations").Preload("Skills").First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[UserProfileRepository.CreateUserProfile] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *UserProfileRepository) FindByID(id uuid.UUID) (*entity.UserProfile, error) {
	var ent entity.UserProfile
	if err := r.DB.Preload("Applicants").Preload("WorkExperiences").Preload("Educations").Preload("Skills").First(&ent, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[UserProfileRepository.FindByID] " + err.Error())
			return nil, err
		}
	}

	return &ent, nil
}

func (r *UserProfileRepository) UpdateUserProfile(ent *entity.UserProfile) (*entity.UserProfile, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[UserProfileRepository.UpdateUserProfile] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Model(&entity.UserProfile{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserProfileRepository.UpdateUserProfile] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserProfileRepository.UpdateUserProfile] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("Applicants").Preload("WorkExperiences").Preload("Educations").Preload("Skills").First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[UserProfileRepository.UpdateUserProfile] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *UserProfileRepository) FindByUserID(userID uuid.UUID) (*entity.UserProfile, error) {
	var ent entity.UserProfile
	if err := r.DB.Preload("Applicants").Preload("WorkExperiences").Preload("Educations").Preload("Skills").Where("user_id = ?", userID).First(&ent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[UserProfileRepository.FindByUserID] " + err.Error())
			return nil, err
		}
	}

	return &ent, nil
}

func (r *UserProfileRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]entity.UserProfile, int64, error) {
	var userProfiles []entity.UserProfile
	var total int64

	query := r.DB.Preload("Applicants").Preload("WorkExperiences").Preload("Educations").Preload("Skills")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if filter["status"] != nil {
		query = query.Where("status = ?", filter["status"])
	}

	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&userProfiles).Error; err != nil {
		r.Log.Error("[UserProfileRepository.FindAllPaginated] " + err.Error())
		return nil, 0, errors.New("[UserProfileRepository.FindAllPaginated] " + err.Error())
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Error("[UserProfileRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	return &userProfiles, total, nil
}

func (r *UserProfileRepository) DeleteUserProfile(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[UserProfileRepository.DeleteUserProfile] " + tx.Error.Error())
		return tx.Error
	}

	var userProfile entity.UserProfile
	if err := tx.Where("id = ?", id).First(&userProfile).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserProfileRepository.DeleteUserProfile] " + err.Error())
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&userProfile).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserProfileRepository.DeleteUserProfile] " + err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserProfileRepository.DeleteUserProfile] " + err.Error())
		return err
	}

	return nil
}
