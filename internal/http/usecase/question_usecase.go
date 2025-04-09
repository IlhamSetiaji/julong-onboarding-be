package usecase

import (
	"errors"
	"fmt"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/dto"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IQuestionUseCase interface {
	CreateOrUpdateQuestions(req *request.CreateOrUpdateQuestions) (*response.SurveyTemplateResponse, error)
	FindByIDAndUserID(questionID, userID string) (*response.QuestionResponse, error)
	FindByID(questionID string) (*entity.Question, error)
}

type QuestionUseCase struct {
	Log                      *logrus.Logger
	Viper                    *viper.Viper
	Repository               repository.IQuestionRepository
	DTO                      dto.IQuestionDTO
	QuestionOptionRepository repository.IQuestionOptionRepository
	UserProfileRepository    repository.IUserProfileRepository
	SurveyTemplateRepository repository.ISurveyTemplateRepository
	SurveyTemplateDTO        dto.ISurveyTemplateDTO
}

func NewQuestionUseCase(
	log *logrus.Logger,
	viper *viper.Viper,
	repo repository.IQuestionRepository,
	qDTO dto.IQuestionDTO,
	qoRepository repository.IQuestionOptionRepository,
	userProfileRepository repository.IUserProfileRepository,
	surveyTemplateRepository repository.ISurveyTemplateRepository,
	surveyTemplateDTO dto.ISurveyTemplateDTO,
) IQuestionUseCase {
	return &QuestionUseCase{
		Log:                      log,
		Viper:                    viper,
		Repository:               repo,
		DTO:                      qDTO,
		QuestionOptionRepository: qoRepository,
		UserProfileRepository:    userProfileRepository,
		SurveyTemplateRepository: surveyTemplateRepository,
		SurveyTemplateDTO:        surveyTemplateDTO,
	}
}

func QuestionUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IQuestionUseCase {
	repo := repository.QuestionRepositoryFactory(log)
	qDTO := dto.QuestionDTOFactory(log, viper)
	qoRepository := repository.QuestionOptionRepositoryFactory(log)
	userProfileRepository := repository.UserProfileRepositoryFactory(log)
	surveyTemplateRepository := repository.SurveyTemplateRepositoryFactory(log)
	surveyTemplateDTO := dto.SurveyTemplateDTOFactory(log, viper)
	return NewQuestionUseCase(log, viper, repo, qDTO, qoRepository, userProfileRepository, surveyTemplateRepository, surveyTemplateDTO)
}

func (u *QuestionUseCase) generateRandomSurveyNumber() (*string, error) {
	// Find the latest survey number
	latestSurvey, err := u.SurveyTemplateRepository.FindLatestSurveyNumber()
	if err != nil {
		u.Log.Error("[SurveyTemplateUseCase.generateRandomSurveyNumber] Error when finding latest survey number: ", err)
		return nil, err
	}

	// Generate the next survey number
	var nextNumber int
	if latestSurvey != nil && latestSurvey.SurveyNumber != "" {
		// Extract the numeric part of the survey number
		var currentNumber int
		_, err := fmt.Sscanf(latestSurvey.SurveyNumber, "SURVEY-%05d", &currentNumber)
		if err != nil {
			u.Log.Error("[SurveyTemplateUseCase.generateRandomSurveyNumber] Error parsing survey number: ", err)
			return nil, err
		}
		nextNumber = currentNumber + 1
	} else {
		// Start from 1 if no survey number exists
		nextNumber = 1
	}

	// Format the next survey number
	newSurveyNumber := fmt.Sprintf("SURVEY-%05d", nextNumber)

	// Check if the generated survey number already exists
	exist, err := u.SurveyTemplateRepository.FindByKeys(map[string]interface{}{
		"survey_number": newSurveyNumber,
	})
	if err != nil {
		u.Log.Error("[SurveyTemplateUseCase.generateRandomSurveyNumber] Error when finding survey number: ", err)
		return nil, err
	}
	if exist != nil {
		// Retry generating a new survey number if it already exists
		return u.generateRandomSurveyNumber()
	}

	return &newSurveyNumber, nil
}

func (uc *QuestionUseCase) CreateOrUpdateQuestions(req *request.CreateOrUpdateQuestions) (*response.SurveyTemplateResponse, error) {
	// check if survey template exist
	if req.SurveyTemplateID != "" {
		tq, err := uc.SurveyTemplateRepository.FindByKeys(map[string]interface{}{
			"id": req.SurveyTemplateID,
		})
		if err != nil {
			uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when finding survey template by id: %s", err.Error())
			return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when finding survey template by id: " + err.Error())
		}

		if tq == nil {
			uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] survey template with id %s not found", req.SurveyTemplateID)
			return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] survey template with id " + req.SurveyTemplateID + " not found")
		}

		_, err = uc.SurveyTemplateRepository.UpdateSurveyTemplate(&entity.SurveyTemplate{
			ID:    tq.ID,
			Title: req.Title,
		})
		if err != nil {
			uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when updating survey template: %s", err.Error())
			return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when updating survey template: " + err.Error())
		}
		req.SurveyTemplateID = tq.ID.String()
	} else {
		surveyNumber, err := uc.generateRandomSurveyNumber()
		if err != nil {
			uc.Log.Error("[SurveyTemplateUseCase.CreateSurveyTemplate] Error when generating random survey number: ", err)
			return nil, err
		}
		tq, err := uc.SurveyTemplateRepository.CreateSurveyTemplate(&entity.SurveyTemplate{
			Title:        req.Title,
			SurveyNumber: *surveyNumber,
			Status:       entity.SURVEY_TEMPLATE_STATUS_ENUM_DRAFT,
		})
		if err != nil {
			uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when creating survey template: %s", err.Error())
			return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when creating survey template: " + err.Error())
		}
		req.SurveyTemplateID = tq.ID.String()
	}

	parsedSurveyTemplateID, err := uuid.Parse(req.SurveyTemplateID)
	if err != nil {
		uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when parsing survey template id: %s", err.Error())
		return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when parsing survey template id: " + err.Error())
	}

	uc.Log.Info("Payload questions: ", req.Questions)

	var questionIDs []uuid.UUID
	for _, question := range req.Questions {
		if question.ID != "" && question.ID != uuid.Nil.String() {
			parsedQuestionID, err := uuid.Parse(question.ID)
			if err != nil {
				uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when parsing question id: %s", err.Error())
				return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when parsing question id: " + err.Error())
			}
			questionIDs = append(questionIDs, parsedQuestionID)
		}
	}

	// delete questions not in ids and survey template id
	if len(questionIDs) > 0 {
		err = uc.Repository.DeleteQuestionsNotInIDsBySurveyTemplateID(parsedSurveyTemplateID, questionIDs)
		if err != nil {
			uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when deleting questions not in ids and survey template id: %s", err.Error())
			return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when deleting questions not in ids and survey template id: " + err.Error())
		}
	}

	// create or update questions
	for i, question := range req.Questions {
		if question.ID != "" && question.ID != uuid.Nil.String() {
			exist, err := uc.Repository.FindByID(uuid.MustParse(question.ID))
			if err != nil {
				uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when finding question by id: %s", err.Error())
				return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when finding question by id: " + err.Error())
			}

			if exist == nil {
				createdQuestion, err := uc.Repository.CreateQuestion(&entity.Question{
					SurveyTemplateID: parsedSurveyTemplateID,
					AnswerTypeID:     uuid.MustParse(question.AnswerTypeID),
					Question:         question.Question,
					Number:           i + 1,
					MaxStars:         question.MaxStars,
					Attachment:       &question.AttachmentPath,
				})
				if err != nil {
					uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question: %s", err.Error())
					return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question: " + err.Error())
				}

				if len(question.QuestionOptions) > 0 {
					for _, questionOption := range question.QuestionOptions {
						_, err := uc.QuestionOptionRepository.CreateQuestionOption(&entity.QuestionOption{
							QuestionID: createdQuestion.ID,
							OptionText: questionOption.OptionText,
						})
						if err != nil {
							uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question option: %s", err.Error())
							return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question option: " + err.Error())
						}
					}
				}
			} else {
				updatedQuestion, err := uc.Repository.UpdateQuestion(&entity.Question{
					ID:               exist.ID,
					SurveyTemplateID: parsedSurveyTemplateID,
					AnswerTypeID:     uuid.MustParse(question.AnswerTypeID),
					Question:         question.Question,
					Number:           i + 1,
					MaxStars:         question.MaxStars,
					Attachment:       &question.AttachmentPath,
				})
				if err != nil {
					uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when updating question: %s", err.Error())
					return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when updating question: " + err.Error())
				}

				// delete question options
				err = uc.QuestionOptionRepository.DeleteQuestionOptionsByQuestionID(updatedQuestion.ID)
				if err != nil {
					uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when deleting question options: %s", err.Error())
					return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when deleting question options: " + err.Error())
				}

				// create question options
				if len(question.QuestionOptions) > 0 {
					for _, questionOption := range question.QuestionOptions {
						_, err := uc.QuestionOptionRepository.CreateQuestionOption(&entity.QuestionOption{
							QuestionID: updatedQuestion.ID,
							OptionText: questionOption.OptionText,
						})
						if err != nil {
							uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question option: %s", err.Error())
							return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question option: " + err.Error())
						}
					}
				}
			}
		} else {
			uc.Log.Info("Payloadku: ", question.Question)

			createdQuestion, err := uc.Repository.CreateQuestion(&entity.Question{
				SurveyTemplateID: parsedSurveyTemplateID,
				AnswerTypeID:     uuid.MustParse(question.AnswerTypeID),
				Question:         question.Question,
				Number:           i + 1,
				MaxStars:         question.MaxStars,
				Attachment:       &question.AttachmentPath,
			})
			if err != nil {
				uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question: %s", err.Error())
				return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question: " + err.Error())
			}

			if len(question.QuestionOptions) > 0 {
				for _, questionOption := range question.QuestionOptions {
					_, err := uc.QuestionOptionRepository.CreateQuestionOption(&entity.QuestionOption{
						QuestionID: createdQuestion.ID,
						OptionText: questionOption.OptionText,
					})
					if err != nil {
						uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question option: %s", err.Error())
						return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question option: " + err.Error())
					}
				}
			}
		}
	}

	// delete questions
	// if len(req.DeletedQuestionIDs) > 0 {
	// 	for _, id := range req.DeletedQuestionIDs {
	// 		err := uc.Repository.DeleteQuestion(uuid.MustParse(id))
	// 		if err != nil {
	// 			uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when deleting question: %s", err.Error())
	// 			return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when deleting question: " + err.Error())
	// 		}
	// 	}
	// }

	tQuestion, err := uc.SurveyTemplateRepository.FindByKeys(map[string]interface{}{
		"id": req.SurveyTemplateID,
	})
	if err != nil {
		uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when finding survey template by id: %s", err.Error())
		return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when finding survey template by id: " + err.Error())
	}

	return uc.SurveyTemplateDTO.ConvertEntityToResponse(tQuestion), nil
}

func (uc *QuestionUseCase) FindByIDAndUserID(questionID, userID string) (*response.QuestionResponse, error) {
	q, err := uc.Repository.FindByID(uuid.MustParse(questionID))
	if err != nil {
		uc.Log.Errorf("[QuestionUseCase.FindByIDAndUserID] error when finding question by id: %s", err.Error())
		return nil, errors.New("[QuestionUseCase.FindByIDAndUserID] error when finding question by id: " + err.Error())
	}

	if q == nil {
		uc.Log.Errorf("[QuestionUseCase.FindByIDAndUserID] question with id %s not found", questionID)
		return nil, errors.New("[QuestionUseCase.FindByIDAndUserID] question with id " + questionID + " not found")
	}

	up, err := uc.UserProfileRepository.FindByUserID(uuid.MustParse(userID))
	if err != nil {
		uc.Log.Errorf("[QuestionUseCase.FindByIDAndUserID] error when finding user profile by user id: %s", err.Error())
		return nil, errors.New("[QuestionUseCase.FindByIDAndUserID] error when finding user profile by user id: " + err.Error())
	}

	if up == nil {
		uc.Log.Errorf("[QuestionUseCase.FindByIDAndUserID] user profile with user id %s not found", userID)
		return nil, errors.New("[QuestionUseCase.FindByIDAndUserID] user profile with user id " + userID + " not found")
	}

	qr, err := uc.Repository.FindQuestionWithResponsesByIDAndUserProfileID(q.ID, up.ID)
	if err != nil {
		uc.Log.Errorf("[QuestionUseCase.FindByIDAndUserID] error when finding question with responses by id and user profile id: %s", err.Error())
		return nil, errors.New("[QuestionUseCase.FindByIDAndUserID] error when finding question with responses by id and user profile id: " + err.Error())
	}

	if qr == nil {
		uc.Log.Errorf("[QuestionUseCase.FindByIDAndUserID] question response not found")
		return nil, errors.New("[QuestionUseCase.FindByIDAndUserID] question response not found")
	}

	return uc.DTO.ConvertEntityToResponse(qr), nil
}

func (uc *QuestionUseCase) FindByID(questionID string) (*entity.Question, error) {
	parsedID, err := uuid.Parse(questionID)
	if err != nil {
		uc.Log.Errorf("[QuestionUseCase.FindByID] error when parsing question id: %s", err.Error())
		return nil, errors.New("[QuestionUseCase.FindByID] error when parsing question id: " + err.Error())
	}

	q, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Errorf("[QuestionUseCase.FindByID] error when finding question by id: %s", err.Error())
		return nil, errors.New("[QuestionUseCase.FindByID] error when finding question by id: " + err.Error())
	}

	return q, nil
}
