package repository

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ISurveyTemplateRepository interface {
	CreateSurveyTemplate(ent *entity.SurveyTemplate) (*entity.SurveyTemplate, error)
	UpdateSurveyTemplate(entsur *entity.SurveyTemplate) (*entity.SurveyTemplate, error)
	DeleteSurveyTemplate(ent *entity.SurveyTemplate) error
	FindByKeys(keys map[string]interface{}) (*entity.SurveyTemplate, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.SurveyTemplate, int64, error)
	FindLatestSurveyNumber() (*entity.SurveyTemplate, error)
}

type SurveyTemplateRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewSurveyTemplateRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *SurveyTemplateRepository {
	return &SurveyTemplateRepository{
		Log: log,
		DB:  db,
	}
}

func SurveyTemplateRepositoryFactory(
	log *logrus.Logger,
) ISurveyTemplateRepository {
	db := config.NewDatabase()
	return NewSurveyTemplateRepository(log, db)
}

func (r *SurveyTemplateRepository) CreateSurveyTemplate(ent *entity.SurveyTemplate) (*entity.SurveyTemplate, error) {
	if err := r.DB.Create(ent).Error; err != nil {
		r.Log.Error("[SurveyTemplateRepository.CreateSurveyTemplate] Error when create survey template: ", err)
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[SurveyTemplateRepository.CreateSurveyTemplate] Error when get survey template: ", err)
		return nil, err
	}

	return ent, nil
}

func (r *SurveyTemplateRepository) UpdateSurveyTemplate(ent *entity.SurveyTemplate) (*entity.SurveyTemplate, error) {
	if err := r.DB.Model(&entity.SurveyTemplate{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		r.Log.Error("[SurveyTemplateRepository.UpdateSurveyTemplate] Error when update survey template: ", err)
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[SurveyTemplateRepository.UpdateSurveyTemplate] Error when get survey template: ", err)
		return nil, err
	}

	return ent, nil
}

func (r *SurveyTemplateRepository) DeleteSurveyTemplate(ent *entity.SurveyTemplate) error {
	if err := r.DB.Delete(ent).Error; err != nil {
		r.Log.Error("[SurveyTemplateRepository.DeleteSurveyTemplate] Error when delete survey template: ", err)
		return err
	}

	return nil
}

func (r *SurveyTemplateRepository) FindByKeys(keys map[string]interface{}) (*entity.SurveyTemplate, error) {
	var ent entity.SurveyTemplate
	if err := r.DB.Preload("Questions.QuestionOptions").Where(keys).First(&ent).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.Log.Error("[SurveyTemplateRepository.FindByKeys] Error when find survey template by keys: ", err)
		return nil, err
	}

	return &ent, nil
}

func (r *SurveyTemplateRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.SurveyTemplate, int64, error) {
	var ents []entity.SurveyTemplate
	var total int64

	query := r.DB.Preload("EmployeeTaskAttachments").Preload("EmployeeTaskChecklists").Where("name LIKE ?", "%"+search+"%")
	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&ents).Error; err != nil {
		r.Log.Error("[SurveyTemplateRepository.FindAllPaginated] Error when get template tasks: ", err)
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Error("[SurveyTemplateRepository.FindAllPaginated] Error when count template tasks: ", err)
		return nil, 0, err
	}

	return &ents, total, nil
}

func (r *SurveyTemplateRepository) FindLatestSurveyNumber() (*entity.SurveyTemplate, error) {
	var latestSurvey entity.SurveyTemplate
	if err := r.DB.Order("created_at desc").First(&latestSurvey).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.Log.Error("[SurveyTemplateRepository.FindLatestSurveyNumber] Error when find latest survey number: ", err)
		return nil, err
	}

	return &latestSurvey, nil
}
