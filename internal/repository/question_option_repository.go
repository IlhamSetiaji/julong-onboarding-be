package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IQuestionOptionRepository interface {
	CreateQuestionOption(ent *entity.QuestionOption) (*entity.QuestionOption, error)
	DeleteQuestionOption(id uuid.UUID) error
	DeleteQuestionOptionsByQuestionID(questionID uuid.UUID) error
}

type QuestionOptionRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewQuestionOptionRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *QuestionOptionRepository {
	return &QuestionOptionRepository{
		Log: log,
		DB:  db,
	}
}

func QuestionOptionRepositoryFactory(
	log *logrus.Logger,
) IQuestionOptionRepository {
	db := config.NewDatabase()
	return NewQuestionOptionRepository(log, db)
}

func (r *QuestionOptionRepository) CreateQuestionOption(ent *entity.QuestionOption) (*entity.QuestionOption, error) {
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

	if err := r.DB.Preload("Question").First(ent, ent.ID).Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *QuestionOptionRepository) DeleteQuestionOption(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var questionOption entity.QuestionOption
	if err := tx.Where("id = ?", id).First(&questionOption).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[QuestionOptionRepository.DeleteQuestionOption] error when finding question option: %v", err)
		return errors.New("[QuestionOptionRepository.DeleteQuestionOption] error when finding question option")
	}

	if err := tx.Delete(&questionOption).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[QuestionOptionRepository.DeleteQuestionOption] error when deleting question option: %v", err)
		return errors.New("[QuestionOptionRepository.DeleteQuestionOption] error when deleting question option")
	}

	if err := tx.Commit().Error; err != nil {
		r.Log.Errorf("[QuestionOptionRepository.DeleteQuestionOption] error when committing transaction: %v", err)
		return errors.New("[QuestionOptionRepository.DeleteQuestionOption] error when committing transaction")
	}

	return nil
}

func (r *QuestionOptionRepository) DeleteQuestionOptionsByQuestionID(questionID uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("question_id = ?", questionID).Delete(&entity.QuestionOption{}).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[QuestionOptionRepository.DeleteQuestionOptionsByQuestionID] error when deleting question options: %v", err)
		return errors.New("[QuestionOptionRepository.DeleteQuestionOptionsByQuestionID] error when deleting question options")
	}

	if err := tx.Commit().Error; err != nil {
		r.Log.Errorf("[QuestionOptionRepository.DeleteQuestionOptionsByQuestionID] error when committing transaction: %v", err)
		return errors.New("[QuestionOptionRepository.DeleteQuestionOptionsByQuestionID] error when committing transaction")
	}

	return nil
}
