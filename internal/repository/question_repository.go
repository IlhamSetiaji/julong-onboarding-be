package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IQuestionRepository interface {
	CreateQuestion(ent *entity.Question) (*entity.Question, error)
	UpdateQuestion(ent *entity.Question) (*entity.Question, error)
	DeleteQuestion(id uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.Question, error)
	FindQuestionWithResponsesByIDAndUserProfileID(questionID, userProfileID uuid.UUID) (*entity.Question, error)
	FindAllByTemplateQuestionIDsAndJobPostingID(templateQuestionIDs []uuid.UUID, jobPostingID uuid.UUID) (*[]entity.Question, error)
	DeleteQuestionsNotInIDsBySurveyTemplateID(surveyTemplateID uuid.UUID, questionIDs []uuid.UUID) error
}

type QuestionRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewQuestionRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *QuestionRepository {
	return &QuestionRepository{
		Log: log,
		DB:  db,
	}
}

func QuestionRepositoryFactory(
	log *logrus.Logger,
) IQuestionRepository {
	db := config.NewDatabase()
	return NewQuestionRepository(log, db)
}

func (r *QuestionRepository) CreateQuestion(ent *entity.Question) (*entity.Question, error) {
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

	if err := r.DB.Preload("AnswerType").First(ent, ent.ID).Error; err != nil {
		return nil, errors.New("[QuestionRepository.Create] error when preloading associations " + err.Error())
	}

	r.Log.Infof("[QuestionRepository.Create] question created: %v", ent)

	return ent, nil
}

func (r *QuestionRepository) UpdateQuestion(ent *entity.Question) (*entity.Question, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Model(&entity.Question{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if err := r.DB.Preload("AnswerType").First(ent, ent.ID).Error; err != nil {
		return nil, errors.New("[QuestionRepository.Update] error when preloading associations " + err.Error())
	}

	r.Log.Infof("[QuestionRepository.Update] question updated: %v", ent)

	return ent, nil
}

func (r *QuestionRepository) DeleteQuestion(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var question entity.Question

	if err := tx.First(&question, id).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[QuestionRepository.DeleteQuestion] error when query question: %v", err)
		return errors.New("[QuestionRepository.DeleteQuestion] error when query question " + err.Error())
	}

	if err := tx.Where("id = ?", id).Delete(&question).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *QuestionRepository) FindByID(id uuid.UUID) (*entity.Question, error) {
	var q entity.Question

	if err := r.DB.
		Where("id = ?", id).
		Preload("QuestionOptions").Preload("AnswerType").
		First(&q).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &q, nil
}

func (r *QuestionRepository) FindQuestionWithResponsesByIDAndUserProfileID(questionID, userProfileID uuid.UUID) (*entity.Question, error) {
	var q entity.Question

	if err := r.DB.
		Preload("QuestionResponses", "user_profile_id = ?", userProfileID).Preload("QuestionResponses.UserProfile").
		Where("id = ?", questionID).
		First(&q).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &q, nil
}

func (r *QuestionRepository) FindAllByTemplateQuestionIDsAndJobPostingID(templateQuestionIDs []uuid.UUID, jobPostingID uuid.UUID) (*[]entity.Question, error) {
	var questions []entity.Question

	if err := r.DB.
		Where("template_question_id IN (?)", templateQuestionIDs).
		Preload("QuestionResponses", "job_posting_id = ?", jobPostingID).
		Preload("QuestionResponses.UserProfile").
		Preload("QuestionOptions").
		Preload("AnswerType").
		Find(&questions).
		Error; err != nil {
		return nil, err
	}

	return &questions, nil
}

func (r *QuestionRepository) DeleteQuestionsNotInIDsBySurveyTemplateID(surveyTemplateID uuid.UUID, questionIDs []uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("survey_template_id = ? AND id NOT IN (?)", surveyTemplateID, questionIDs).Delete(&entity.Question{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
