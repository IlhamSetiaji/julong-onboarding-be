package repository

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ISurveyResponseRepository interface {
	CreateSurveyResponse(ent *entity.SurveyResponse) (*entity.SurveyResponse, error)
	UpdateSurveyResponse(ent *entity.SurveyResponse) (*entity.SurveyResponse, error)
	DeleteSurveyResponse(id uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.SurveyResponse, error)
	FindAllByQuestionID(questionID uuid.UUID) ([]entity.SurveyResponse, error)
	GetAllByKeys(keys map[string]interface{}) ([]entity.SurveyResponse, error)
	DeleteByQuestionID(questionID uuid.UUID) error
	DeleteByQuestionIDs(questionIDs []uuid.UUID) error
	DeleteNotInIDsAndQuestionID(questionID uuid.UUID, ids []uuid.UUID) error
}

type SurveyResponseRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewSurveyResponseRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *SurveyResponseRepository {
	return &SurveyResponseRepository{
		Log: log,
		DB:  db,
	}
}

func SurveyResponseRepositoryFactory(
	log *logrus.Logger,
) ISurveyResponseRepository {
	db := config.NewDatabase()
	return NewSurveyResponseRepository(log, db)
}

func (r *SurveyResponseRepository) CreateSurveyResponse(ent *entity.SurveyResponse) (*entity.SurveyResponse, error) {
	if err := r.DB.Create(ent).Error; err != nil {
		r.Log.Error("[SurveyResponseRepository.CreateSurveyResponse] Error when create survey response: ", err)
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[SurveyResponseRepository.CreateSurveyResponse] Error when get survey response: ", err)
		return nil, err
	}

	return ent, nil
}

func (r *SurveyResponseRepository) UpdateSurveyResponse(ent *entity.SurveyResponse) (*entity.SurveyResponse, error) {
	if err := r.DB.Model(&entity.SurveyResponse{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		r.Log.Error("[SurveyResponseRepository.UpdateSurveyResponse] Error when update survey response: ", err)
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[SurveyResponseRepository.UpdateSurveyResponse] Error when get survey response: ", err)
		return nil, err
	}

	return ent, nil
}

func (r *SurveyResponseRepository) DeleteSurveyResponse(id uuid.UUID) error {
	var surveyResponse entity.SurveyResponse
	if err := r.DB.Where("id = ?", id).Delete(&surveyResponse).Error; err != nil {
		r.Log.Error("[SurveyResponseRepository.DeleteSurveyResponse] Error when delete survey response: ", err)
		return err
	}

	return nil
}

func (r *SurveyResponseRepository) FindByID(id uuid.UUID) (*entity.SurveyResponse, error) {
	var surveyResponse entity.SurveyResponse
	if err := r.DB.Where("id = ?", id).First(&surveyResponse).Error; err != nil {
		r.Log.Error("[SurveyResponseRepository.FindByID] Error when get survey response: ", err)
		return nil, err
	}

	return &surveyResponse, nil
}

func (r *SurveyResponseRepository) FindAllByQuestionID(questionID uuid.UUID) ([]entity.SurveyResponse, error) {
	var surveyResponses []entity.SurveyResponse
	if err := r.DB.Where("question_id = ?", questionID).Find(&surveyResponses).Error; err != nil {
		r.Log.Error("[SurveyResponseRepository.FindAllByQuestionID] Error when get survey responses: ", err)
		return nil, err
	}

	return surveyResponses, nil
}

func (r *SurveyResponseRepository) GetAllByKeys(keys map[string]interface{}) ([]entity.SurveyResponse, error) {
	var surveyResponses []entity.SurveyResponse
	if err := r.DB.Where(keys).Find(&surveyResponses).Error; err != nil {
		r.Log.Error("[SurveyResponseRepository.GetAllByKeys] Error when get survey responses: ", err)
		return nil, err
	}

	return surveyResponses, nil
}

func (r *SurveyResponseRepository) DeleteByQuestionID(questionID uuid.UUID) error {
	var surveyResponse entity.SurveyResponse
	if err := r.DB.Where("question_id = ?", questionID).Delete(&surveyResponse).Error; err != nil {
		r.Log.Error("[SurveyResponseRepository.DeleteByQuestionID] Error when delete survey response: ", err)
		return err
	}

	return nil
}

func (r *SurveyResponseRepository) DeleteByQuestionIDs(questionIDs []uuid.UUID) error {
	var surveyResponse entity.SurveyResponse
	if err := r.DB.Where("question_id IN ?", questionIDs).Delete(&surveyResponse).Error; err != nil {
		r.Log.Error("[SurveyResponseRepository.DeleteByQuestionIDs] Error when delete survey response: ", err)
		return err
	}

	return nil
}

func (r *SurveyResponseRepository) DeleteNotInIDsAndQuestionID(questionID uuid.UUID, ids []uuid.UUID) error {
	var surveyResponse entity.SurveyResponse
	if err := r.DB.Where("question_id = ? AND id NOT IN ?", questionID, ids).Delete(&surveyResponse).Error; err != nil {
		r.Log.Error("[SurveyResponseRepository.DeleteNotInIDsAndQuestionID] Error when delete survey response: ", err)
		return err
	}

	return nil
}
