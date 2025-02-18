package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ICoverRepository interface {
	CreateCoverRepository(ent *entity.Cover) (*entity.Cover, error)
	UpdateCoverRepository(ent *entity.Cover) (*entity.Cover, error)
	DeleteCoverRepository(id uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.Cover, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.Cover, int64, error)
}

type CoverRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewCoverRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *CoverRepository {
	return &CoverRepository{
		Log: log,
		DB:  db,
	}
}

func CoverRepositoryFactory(
	log *logrus.Logger,
) ICoverRepository {
	db := config.NewDatabase()
	return NewCoverRepository(log, db)
}

func (r *CoverRepository) CreateCoverRepository(ent *entity.Cover) (*entity.Cover, error) {
	if err := r.DB.Create(ent).Error; err != nil {
		r.Log.Error("[CoverRepository.CreateCoverRepository] Error when create cover: ", err)
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[CoverRepository.CreateCoverRepository] Error when get cover: ", err)
		return nil, err
	}

	return ent, nil
}

func (r *CoverRepository) UpdateCoverRepository(ent *entity.Cover) (*entity.Cover, error) {
	if err := r.DB.Model(&entity.Cover{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		r.Log.Error("[CoverRepository.UpdateCoverRepository] Error when update cover: ", err)
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[CoverRepository.UpdateCoverRepository] Error when get cover: ", err)
		return nil, err
	}

	return ent, nil
}

func (r *CoverRepository) DeleteCoverRepository(id uuid.UUID) error {
	if err := r.DB.Where("id = ?", id).Delete(&entity.Cover{}).Error; err != nil {
		r.Log.Error("[CoverRepository.DeleteCoverRepository] Error when delete cover: ", err)
		return err
	}

	return nil
}

func (r *CoverRepository) FindByID(id uuid.UUID) (*entity.Cover, error) {
	var cover entity.Cover
	if err := r.DB.Where("id = ?", id).First(&cover).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[CoverRepository.FindByID] Error when get cover: ", err)
			return nil, err
		}
	}

	return &cover, nil
}

func (r *CoverRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.Cover, int64, error) {
	var covers []entity.Cover
	var total int64

	db := r.DB.Model(&entity.Cover{})

	if search != "" {
		db = db.Where("path LIKE ?", "%"+search+"%")
	}

	for key, value := range sort {
		db = db.Order(key + " " + value.(string))
	}

	if err := db.Count(&total).Error; err != nil {
		r.Log.Error("[CoverRepository.FindAllPaginated] Error when count cover: ", err)
		return nil, 0, err
	}

	if err := db.Order("created_at desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&covers).Error; err != nil {
		r.Log.Error("[CoverRepository.FindAllPaginated] Error when get covers: ", err)
		return nil, 0, err
	}

	return &covers, total, nil
}
