package usecase

import (
	"errors"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/dto"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ISurveyResponseUseCase interface {
	CreateOrUpdateSurveyResponses(req *request.SurveyResponseRequest) (*response.QuestionResponse, error)
}

type SurveyResponseUseCase struct {
	Log                      *logrus.Logger
	Viper                    *viper.Viper
	QuestionRepository       repository.IQuestionRepository
	SurveyTemplateRepository repository.ISurveyTemplateRepository
	SurveyResponseRepository repository.ISurveyResponseRepository
	EmployeeMessage          messaging.IEmployeeMessage
	QuestionDTO              dto.IQuestionDTO
	EmployeeTaskRepository   repository.IEmployeeTaskRepository
}

func NewSurveyResponseUseCase(
	Log *logrus.Logger,
	Viper *viper.Viper,
	QuestionRepository repository.IQuestionRepository,
	SurveyTemplateRepository repository.ISurveyTemplateRepository,
	SurveyResponseRepository repository.ISurveyResponseRepository,
	EmployeeMessage messaging.IEmployeeMessage,
	QuestionDTO dto.IQuestionDTO,
	EmployeeTaskRepository repository.IEmployeeTaskRepository,
) ISurveyResponseUseCase {
	return &SurveyResponseUseCase{
		Log:                      Log,
		Viper:                    Viper,
		QuestionRepository:       QuestionRepository,
		SurveyTemplateRepository: SurveyTemplateRepository,
		SurveyResponseRepository: SurveyResponseRepository,
		EmployeeMessage:          EmployeeMessage,
		QuestionDTO:              QuestionDTO,
		EmployeeTaskRepository:   EmployeeTaskRepository,
	}
}

func SurveyResponseUseCaseFactory(
	Log *logrus.Logger,
	Viper *viper.Viper,
) ISurveyResponseUseCase {
	questionRepository := repository.QuestionRepositoryFactory(Log)
	surveyTemplateRepository := repository.SurveyTemplateRepositoryFactory(Log)
	surveyResponseRepository := repository.SurveyResponseRepositoryFactory(Log)
	employeeMessage := messaging.EmployeeMessageFactory(Log)
	questionDTO := dto.QuestionDTOFactory(Log, Viper)
	employeeTaskRepository := repository.EmployeeTaskRepositoryFactory(Log)

	return NewSurveyResponseUseCase(
		Log,
		Viper,
		questionRepository,
		surveyTemplateRepository,
		surveyResponseRepository,
		employeeMessage,
		questionDTO,
		employeeTaskRepository,
	)
}

func (uc *SurveyResponseUseCase) CreateOrUpdateSurveyResponses(req *request.SurveyResponseRequest) (*response.QuestionResponse, error) {
	// check if question is exist
	parsedQuestionID, err := uuid.Parse(req.QuestionID)
	if err != nil {
		uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateSurveyResponses] error when parsing question id: %s", err.Error())
		return nil, err
	}

	question, err := uc.QuestionRepository.FindByID(parsedQuestionID)
	if err != nil {
		uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateSurveyResponses] error when finding question by id: %s", err.Error())
		return nil, err
	}

	if question == nil {
		uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateSurveyResponses] question with id %s not found", req.QuestionID)
		return nil, err
	}

	var employeeTaskUUID uuid.UUID

	var answerIDs []uuid.UUID
	for _, ans := range req.Answers {
		if ans.ID != nil {
			parsedAnswerID, err := uuid.Parse(*ans.ID)
			if err != nil {
				uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateSurveyResponses] error when parsing answer id: %s", err.Error())
				return nil, err
			}
			answerIDs = append(answerIDs, parsedAnswerID)
		}
	}

	// delete answers by question id and not in ids
	if len(answerIDs) > 0 {
		err := uc.SurveyResponseRepository.DeleteNotInIDsAndQuestionID(parsedQuestionID, answerIDs)
		if err != nil {
			uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateSurveyResponses] error when deleting answers by question id and ids: %s", err.Error())
			return nil, err
		}
	}

	// create or update answers
	for _, ans := range req.Answers {
		parsedSurveyTemplateID, err := uuid.Parse(ans.SurveyTemplateID)
		if err != nil {
			uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateSurveyResponses] error when parsing survey template id: %s", err.Error())
			return nil, err
		}
		jp, err := uc.SurveyTemplateRepository.FindByKeys(map[string]interface{}{
			"id": parsedSurveyTemplateID,
		})
		if err != nil {
			uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateSurveyResponses] error when finding survey template by id: %s", err.Error())
			return nil, err
		}
		if jp == nil {
			uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateSurveyResponses] survey template with id %s not found", ans.SurveyTemplateID)
			return nil, errors.New("survey template not found")
		}

		parsedEmployeeTaskID, err := uuid.Parse(ans.EmployeeTaskID)
		if err != nil {
			uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateSurveyResponses] error when parsing employee task id: %s", err.Error())
			return nil, err
		}
		employeeTaskUUID = parsedEmployeeTaskID
		up, err := uc.EmployeeTaskRepository.FindByID(parsedEmployeeTaskID)
		if err != nil {
			uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateSurveyResponses] error when finding employee task by id: %s", err.Error())
			return nil, err
		}
		if up == nil {
			uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateSurveyResponses] employee task with id %s not found", ans.EmployeeTaskID)
			return nil, errors.New("employee task not found")
		}
		uc.Log.Info("Halooo")

		// check if answer is exist
		if ans.ID != nil {
			parsedAnswerID, err := uuid.Parse(*ans.ID)
			if err != nil {
				uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateSurveyResponses] error when parsing answer id: %s", err.Error())
				return nil, err
			}
			exist, err := uc.SurveyResponseRepository.FindByID(parsedAnswerID)
			if err != nil {
				uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateSurveyResponses] error when finding answer by id: %s", err.Error())
				return nil, err
			}

			if exist == nil {
				uc.Log.Infof("kontol: %+v", ans)
				_, err := uc.SurveyResponseRepository.CreateSurveyResponse(&entity.SurveyResponse{
					QuestionID:       question.ID,
					SurveyTemplateID: jp.ID,
					EmployeeTaskID:   up.ID,
					Answer:           ans.Answer,
					AnswerFile:       ans.AnswerPath,
				})
				if err != nil {
					uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateSurveyResponses] error when creating answer: %s", err.Error())
					return nil, err
				}
			} else {
				uc.Log.Infof("memek: %+v", ans)
				_, err := uc.SurveyResponseRepository.UpdateSurveyResponse(&entity.SurveyResponse{
					ID:               exist.ID,
					QuestionID:       question.ID,
					SurveyTemplateID: jp.ID,
					EmployeeTaskID:   up.ID,
					Answer:           ans.Answer,
					AnswerFile:       ans.AnswerPath,
				})

				if err != nil {
					uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateSurveyResponses] error when updating answer: %s", err.Error())
					return nil, err
				}
			}
		} else {
			uc.Log.Infof("cok: %+v", ans)
			hasil, err := uc.SurveyResponseRepository.CreateSurveyResponse(&entity.SurveyResponse{
				QuestionID:       question.ID,
				SurveyTemplateID: jp.ID,
				EmployeeTaskID:   up.ID,
				Answer:           ans.Answer,
				AnswerFile:       ans.AnswerPath,
			})
			if err != nil {
				uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateSurveyResponses] error when creating answer: %s", err.Error())
				return nil, err
			}

			uc.Log.Infof("hasil: %+v", hasil)
		}
	}

	rQuestion, err := uc.QuestionRepository.FindQuestionWithResponsesByIDAndUserProfileID(question.ID, employeeTaskUUID)
	if err != nil {
		uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateSurveyResponses] error when finding question by id: %s", err.Error())
		return nil, err
	}
	if rQuestion == nil {
		uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateSurveyResponses] question with id %s not found", req.QuestionID)
		return nil, err
	}

	// embed url to answer file
	for _, qr := range rQuestion.SurveyResponses {
		if qr.AnswerFile != "" {
			qr.AnswerFile = uc.Viper.GetString("app.url") + "/" + qr.AnswerFile
			uc.Log.Infof("[QuestionResponseUseCase.CreateOrUpdateSurveyResponses] answer file url: %s", qr.AnswerFile)
		}
	}

	return uc.QuestionDTO.ConvertEntityToResponse(rQuestion), nil
}
