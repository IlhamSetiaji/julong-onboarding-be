package usecase

import (
	"errors"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/dto"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ITemplateTaskUseCase interface {
	CreateTemplateTask(req *request.CreateTemplateTaskRequest) (*response.TemplateTaskResponse, error)
	UpdateTemplateTask(req *request.UpdateTemplateTaskRequest) (*response.TemplateTaskResponse, error)
	DeleteTemplateTask(id uuid.UUID) error
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, status entity.TemplateTaskStatusEnum) (*[]response.TemplateTaskResponse, int64, error)
	FindByID(id uuid.UUID) (*response.TemplateTaskResponse, error)
}

type TemplateTaskUseCase struct {
	Log                              *logrus.Logger
	DTO                              dto.ITemplateTaskDTO
	Repository                       repository.ITemplateTaskRepository
	TemplateTaskAttachmentRepository repository.ITemplateTaskAttachmentRepository
	TemplateTaskChecklistRepository  repository.ITemplateTaskChecklistRepository
	Viper                            *viper.Viper
	SurveyTemplateRepository         repository.ISurveyTemplateRepository
}

func NewTemplateTaskUseCase(
	log *logrus.Logger,
	dto dto.ITemplateTaskDTO,
	repo repository.ITemplateTaskRepository,
	attachmentRepo repository.ITemplateTaskAttachmentRepository,
	checklistRepo repository.ITemplateTaskChecklistRepository,
	viper *viper.Viper,
	surveyTemplateRepo repository.ISurveyTemplateRepository,
) ITemplateTaskUseCase {
	return &TemplateTaskUseCase{
		Log:                              log,
		DTO:                              dto,
		Repository:                       repo,
		TemplateTaskAttachmentRepository: attachmentRepo,
		TemplateTaskChecklistRepository:  checklistRepo,
		Viper:                            viper,
		SurveyTemplateRepository:         surveyTemplateRepo,
	}
}

func TemplateTaskUseCaseFactory(log *logrus.Logger, viper *viper.Viper) ITemplateTaskUseCase {
	dto := dto.TemplateTaskDTOFactory(log, viper)
	repo := repository.TemplateTaskRepositoryFactory(log)
	attachmentRepo := repository.TemplateTaskAttachmentRepositoryFactory(log)
	checklistRepo := repository.TemplateTaskChecklistRepositoryFactory(log)
	surveyTemplateRepo := repository.SurveyTemplateRepositoryFactory(log)
	return NewTemplateTaskUseCase(log, dto, repo, attachmentRepo, checklistRepo, viper, surveyTemplateRepo)
}

func (uc *TemplateTaskUseCase) CreateTemplateTask(req *request.CreateTemplateTaskRequest) (*response.TemplateTaskResponse, error) {
	var duration *int
	if req.DueDuration != nil {
		duration = req.DueDuration
	}

	var surveyTemplateUUID *uuid.UUID
	if req.SurveyTemplateID != nil && *req.SurveyTemplateID != "" {
		parsedSurveyTemplateID, err := uuid.Parse(*req.SurveyTemplateID)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error parsing survey template id: ", err)
			return nil, err
		}
		surveyTemplate, err := uc.SurveyTemplateRepository.FindByKeys(map[string]interface{}{
			"id": parsedSurveyTemplateID,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error finding survey template by id: ", err)
			return nil, err
		}
		if surveyTemplate == nil {
			return nil, errors.New("survey template not found")
		}

		surveyTemplateUUID = &parsedSurveyTemplateID
	}

	templateTask, err := uc.Repository.CreateTemplateTask(&entity.TemplateTask{
		Name:             req.Name,
		CoverPath:        &req.CoverPath,
		Priority:         entity.TemplateTaskPriorityEnum(req.Priority),
		DueDuration:      duration,
		Status:           entity.TemplateTaskStatusEnum(req.Status),
		Description:      req.Description,
		Source:           "ONBOARDING",
		OrganizationType: req.OrganizationType,
		SurveyTemplateID: surveyTemplateUUID,
	})
	if err != nil {
		uc.Log.Error("[TemplateTaskUseCase.CreateTemplateTask] " + err.Error())
		return nil, err
	}

	// create template task attachments
	for _, attachment := range req.TemplateTaskAttachments {
		_, err := uc.TemplateTaskAttachmentRepository.CreateTemplateTaskAttachment(&entity.TemplateTaskAttachment{
			TemplateTaskID: templateTask.ID,
			Path:           attachment.Path,
		})
		if err != nil {
			uc.Log.Error("[TemplateTaskUseCase.CreateTemplateTask] " + err.Error())
			return nil, err
		}
	}

	// delete template task checklists
	err = uc.TemplateTaskChecklistRepository.DeleteByTemplateTaskID(templateTask.ID.String())
	if err != nil {
		uc.Log.Error("[TemplateTaskUseCase.CreateTemplateTask] " + err.Error())
		return nil, err
	}
	// create template task checklists
	for _, checklist := range req.TemplateTaskChecklists {
		if checklist.ID != nil {
			parsedChecklistID, err := uuid.Parse(*checklist.ID)
			if err != nil {
				uc.Log.Error("[TemplateTaskUseCase.CreateTemplateTask] " + err.Error())
				return nil, err
			}
			exist, err := uc.TemplateTaskChecklistRepository.FindByKeys(map[string]interface{}{
				"id": parsedChecklistID,
			})
			if err != nil {
				uc.Log.Error("[TemplateTaskUseCase.CreateTemplateTask] " + err.Error())
				return nil, err
			}
			if exist == nil {
				_, err := uc.TemplateTaskChecklistRepository.CreateTaskChecklistRepository(&entity.TemplateTaskChecklist{
					TemplateTaskID: templateTask.ID,
					Name:           checklist.Name,
				})
				if err != nil {
					uc.Log.Error("[TemplateTaskUseCase.CreateTemplateTask] " + err.Error())
					return nil, err
				}
			} else {
				_, err := uc.TemplateTaskChecklistRepository.UpdateTaskChecklistRepository(&entity.TemplateTaskChecklist{
					ID:             parsedChecklistID,
					TemplateTaskID: templateTask.ID,
					Name:           checklist.Name,
				})
				if err != nil {
					uc.Log.Error("[TemplateTaskUseCase.CreateTemplateTask] " + err.Error())
					return nil, err
				}
			}
		} else {
			_, err := uc.TemplateTaskChecklistRepository.CreateTaskChecklistRepository(&entity.TemplateTaskChecklist{
				TemplateTaskID: templateTask.ID,
				Name:           checklist.Name,
			})
			if err != nil {
				uc.Log.Error("[TemplateTaskUseCase.CreateTemplateTask] " + err.Error())
				return nil, err
			}
		}
	}

	findById, err := uc.Repository.FindByID(templateTask.ID)
	if err != nil {
		uc.Log.Error("[TemplateTaskUseCase.CreateTemplateTask] " + err.Error())
		return nil, err
	}
	if findById == nil {
		return nil, errors.New("Template task not found")
	}
	return uc.DTO.ConvertEntityToResponse(findById), nil
}

func (uc *TemplateTaskUseCase) UpdateTemplateTask(req *request.UpdateTemplateTaskRequest) (*response.TemplateTaskResponse, error) {
	parsedId, err := uuid.Parse(req.ID)
	if err != nil {
		uc.Log.Error("[TemplateTaskUseCase.UpdateTemplateTask] " + err.Error())
		return nil, err
	}
	ttExist, err := uc.Repository.FindByID(parsedId)
	if err != nil {
		uc.Log.Error("[TemplateTaskUseCase.UpdateTemplateTask] " + err.Error())
		return nil, err
	}
	if ttExist == nil {
		return nil, errors.New("Template task not found")
	}
	var duration *int
	if req.DueDuration != nil {
		duration = req.DueDuration
	}

	var surveyTemplateUUID *uuid.UUID
	if req.SurveyTemplateID != nil && *req.SurveyTemplateID != "" {
		parsedSurveyTemplateID, err := uuid.Parse(*req.SurveyTemplateID)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error parsing survey template id: ", err)
			return nil, err
		}
		surveyTemplate, err := uc.SurveyTemplateRepository.FindByKeys(map[string]interface{}{
			"id": parsedSurveyTemplateID,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error finding survey template by id: ", err)
			return nil, err
		}
		if surveyTemplate == nil {
			return nil, errors.New("survey template not found")
		}

		surveyTemplateUUID = &parsedSurveyTemplateID
	}

	templateTask, err := uc.Repository.UpdateTemplateTask(&entity.TemplateTask{
		ID:               parsedId,
		Name:             req.Name,
		SurveyTemplateID: surveyTemplateUUID,
		CoverPath:        &req.CoverPath,
		Priority:         entity.TemplateTaskPriorityEnum(req.Priority),
		DueDuration:      duration,
		Status:           entity.TemplateTaskStatusEnum(req.Status),
		Description:      req.Description,
		Source:           "ONBOARDING",
		OrganizationType: req.OrganizationType,
	})
	if err != nil {
		uc.Log.Error("[TemplateTaskUseCase.CreateTemplateTask] " + err.Error())
		return nil, err
	}

	// create template task attachments
	for _, attachment := range req.TemplateTaskAttachments {
		_, err := uc.TemplateTaskAttachmentRepository.CreateTemplateTaskAttachment(&entity.TemplateTaskAttachment{
			TemplateTaskID: templateTask.ID,
			Path:           attachment.Path,
		})
		if err != nil {
			uc.Log.Error("[TemplateTaskUseCase.CreateTemplateTask] " + err.Error())
			return nil, err
		}
	}

	// delete template task checklists
	err = uc.TemplateTaskChecklistRepository.DeleteByTemplateTaskID(templateTask.ID.String())
	if err != nil {
		uc.Log.Error("[TemplateTaskUseCase.CreateTemplateTask] " + err.Error())
		return nil, err
	}
	// create template task checklists
	for _, checklist := range req.TemplateTaskChecklists {
		if checklist.ID != nil {
			parsedChecklistID, err := uuid.Parse(*checklist.ID)
			if err != nil {
				uc.Log.Error("[TemplateTaskUseCase.CreateTemplateTask] " + err.Error())
				return nil, err
			}
			exist, err := uc.TemplateTaskChecklistRepository.FindByKeys(map[string]interface{}{
				"id": parsedChecklistID,
			})
			if err != nil {
				uc.Log.Error("[TemplateTaskUseCase.CreateTemplateTask] " + err.Error())
				return nil, err
			}
			if exist == nil {
				_, err := uc.TemplateTaskChecklistRepository.CreateTaskChecklistRepository(&entity.TemplateTaskChecklist{
					TemplateTaskID: templateTask.ID,
					Name:           checklist.Name,
				})
				if err != nil {
					uc.Log.Error("[TemplateTaskUseCase.CreateTemplateTask] " + err.Error())
					return nil, err
				}
			} else {
				_, err := uc.TemplateTaskChecklistRepository.UpdateTaskChecklistRepository(&entity.TemplateTaskChecklist{
					ID:             parsedChecklistID,
					TemplateTaskID: templateTask.ID,
					Name:           checklist.Name,
				})
				if err != nil {
					uc.Log.Error("[TemplateTaskUseCase.CreateTemplateTask] " + err.Error())
					return nil, err
				}
			}
		} else {
			_, err := uc.TemplateTaskChecklistRepository.CreateTaskChecklistRepository(&entity.TemplateTaskChecklist{
				TemplateTaskID: templateTask.ID,
				Name:           checklist.Name,
			})
			if err != nil {
				uc.Log.Error("[TemplateTaskUseCase.CreateTemplateTask] " + err.Error())
				return nil, err
			}
		}
	}

	findById, err := uc.Repository.FindByID(templateTask.ID)
	if err != nil {
		uc.Log.Error("[TemplateTaskUseCase.CreateTemplateTask] " + err.Error())
		return nil, err
	}
	if findById == nil {
		return nil, errors.New("Template task not found")
	}
	return uc.DTO.ConvertEntityToResponse(findById), nil
}

func (uc *TemplateTaskUseCase) DeleteTemplateTask(id uuid.UUID) error {
	templateTask, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[TemplateTaskUseCase.DeleteTemplateTask] " + err.Error())
		return err
	}

	if templateTask == nil {
		return errors.New("Template task not found")
	}

	err = uc.Repository.DeleteTemplateTask(templateTask)
	if err != nil {
		uc.Log.Error("[TemplateTaskUseCase.DeleteTemplateTask] " + err.Error())
		return err
	}

	return nil
}

func (uc *TemplateTaskUseCase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, status entity.TemplateTaskStatusEnum) (*[]response.TemplateTaskResponse, int64, error) {
	entities, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort, status)
	if err != nil {
		uc.Log.Error("[TemplateTaskUseCase.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	var responses []response.TemplateTaskResponse
	for _, entity := range *entities {
		res := uc.DTO.ConvertEntityToResponse(&entity)
		responses = append(responses, *res)
	}

	return &responses, total, nil
}

func (uc *TemplateTaskUseCase) FindByID(id uuid.UUID) (*response.TemplateTaskResponse, error) {
	templateTask, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[TemplateTaskUseCase.FindByID] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(templateTask), nil
}
