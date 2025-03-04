package usecase

import (
	"errors"
	"time"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/dto"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IEmployeeTaskUseCase interface {
	CreateEmployeeTask(req *request.CreateEmployeeTaskRequest) (*response.EmployeeTaskResponse, error)
	UpdateEmployeeTask(req *request.UpdateEmployeeTaskRequest) (*response.EmployeeTaskResponse, error)
	DeleteEmployeeTask(id uuid.UUID) error
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.EmployeeTaskResponse, int64, error)
	FindByID(id uuid.UUID) (*response.EmployeeTaskResponse, error)
	CountByKanbanAndEmployeeID(kanban entity.EmployeeTaskKanbanEnum, employeeID uuid.UUID) (int64, error)
	FindAllByEmployeeID(employeeID uuid.UUID) (*response.EmployeeTaskKanbanResponse, error)
	FindAllByEmployeeIDAndKanbanPaginated(employeeID uuid.UUID, kanban entity.EmployeeTaskKanbanEnum, page, pageSize int, search string, sort map[string]interface{}) (*[]response.EmployeeTaskResponse, int64, error)
	UpdateEmployeeTaskOnly(req *request.UpdateEmployeeTaskOnlyRequest) (*response.EmployeeTaskResponse, error)
	CreateEmployeeTasksForRecruitment(req *request.CreateEmployeeTasksForRecruitment) error
	CountKanbanProgressByEmployeeID(employeeID uuid.UUID) (*response.EmployeeTaskProgressResponse, error)
}

type EmployeeTaskUseCase struct {
	Log                              *logrus.Logger
	DTO                              dto.IEmployeeTaskDTO
	Repository                       repository.IEmployeeTaskRepository
	Viper                            *viper.Viper
	TemplateTaskRepository           repository.ITemplateTaskRepository
	EmployeeTaskAttachmentRepository repository.IEmployeeTaskAttachmentRepository
	EmployeeTaskChecklistRepository  repository.IEmployeeTaskChecklistRepository
	EmployeeHiringRepository         repository.IEmployeeHiringRepository
}

func NewEmployeeTaskUseCase(
	log *logrus.Logger,
	dto dto.IEmployeeTaskDTO,
	repository repository.IEmployeeTaskRepository,
	viper *viper.Viper,
	templateTaskRepository repository.ITemplateTaskRepository,
	etaRepo repository.IEmployeeTaskAttachmentRepository,
	etcRepo repository.IEmployeeTaskChecklistRepository,
	ehRepo repository.IEmployeeHiringRepository,
) IEmployeeTaskUseCase {
	return &EmployeeTaskUseCase{
		Log:                              log,
		DTO:                              dto,
		Repository:                       repository,
		Viper:                            viper,
		TemplateTaskRepository:           templateTaskRepository,
		EmployeeTaskAttachmentRepository: etaRepo,
		EmployeeTaskChecklistRepository:  etcRepo,
		EmployeeHiringRepository:         ehRepo,
	}
}

func EmployeeTaskUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IEmployeeTaskUseCase {
	etDTO := dto.EmployeeTaskDTOFactory(log, viper)
	repo := repository.EmployeeTaskRepositoryFactory(log)
	ttRepository := repository.TemplateTaskRepositoryFactory(log)
	etaRepo := repository.EmployeeTaskAttachmentRepositoryFactory(log)
	etcRepo := repository.EmployeeTaskChecklistRepositoryFactory(log)
	ehRepo := repository.EmployeeHiringRepositoryFactory(log)
	return NewEmployeeTaskUseCase(log, etDTO, repo, viper, ttRepository, etaRepo, etcRepo, ehRepo)
}

func (uc *EmployeeTaskUseCase) CreateEmployeeTask(req *request.CreateEmployeeTaskRequest) (*response.EmployeeTaskResponse, error) {
	var templateTaskUUID *uuid.UUID
	if req.TemplateTaskID != nil && *req.TemplateTaskID != "" {
		parsedTemplateTaskID, err := uuid.Parse(*req.TemplateTaskID)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error parsing template task id: ", err)
			return nil, err
		}
		templateTask, err := uc.TemplateTaskRepository.FindByID(parsedTemplateTaskID)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error finding template task by id: ", err)
			return nil, err
		}
		if templateTask == nil {
			return nil, errors.New("template task not found")
		}

		templateTaskUUID = &parsedTemplateTaskID
	}

	parsedEmployeeID, err := uuid.Parse(*req.EmployeeID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error parsing employee id: ", err)
		return nil, err
	}

	parsedStartDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error parsing start date: ", err)
		return nil, err
	}

	parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error parsing end date: ", err)
		return nil, err
	}

	employeeTask, err := uc.Repository.CreateEmployeeTask(&entity.EmployeeTask{
		CoverPath:      req.CoverPath,
		EmployeeID:     &parsedEmployeeID,
		TemplateTaskID: templateTaskUUID,
		Name:           req.Name,
		Priority:       entity.EmployeeTaskPriorityEnum(req.Priority),
		Description:    req.Description,
		StartDate:      parsedStartDate,
		EndDate:        parsedEndDate,
		Source:         "ONBOARDING",
	})
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error creating employee task: ", err)
		return nil, err
	}

	// create employee task attachments
	for _, attachmentReq := range req.EmployeeTaskAttachments {
		_, err := uc.EmployeeTaskAttachmentRepository.CreateEmployeeTaskAttachment(&entity.EmployeeTaskAttachment{
			EmployeeTaskID: employeeTask.ID,
			Path:           attachmentReq.Path,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error creating employee task attachment: ", err)
			return nil, err
		}
	}

	// delete employee task checklists
	err = uc.EmployeeTaskChecklistRepository.DeleteByEmployeeTaskID(employeeTask.ID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error deleting employee task checklists: ", err)
		return nil, err
	}

	// create employee task checklists
	for _, checklistReq := range req.EmployeeTaskChecklists {
		if checklistReq.ID != nil {
			parsedChecklistID, err := uuid.Parse(*checklistReq.ID)
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error parsing checklist id: ", err)
				return nil, err
			}
			exist, err := uc.EmployeeTaskChecklistRepository.FindByKeys(map[string]interface{}{
				"id": parsedChecklistID,
			})
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error finding checklist by id: ", err)
				return nil, err
			}
			if exist == nil {
				_, err := uc.EmployeeTaskChecklistRepository.CreateEmployeeTaskChecklist(&entity.EmployeeTaskChecklist{
					EmployeeTaskID: employeeTask.ID,
					Name:           checklistReq.Name,
				})
				if err != nil {
					uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error creating employee task checklist: ", err)
					return nil, err
				}
			} else {
				var verifiedBy *uuid.UUID
				if checklistReq.VerifiedBy != nil {
					parsedVerifiedBy, err := uuid.Parse(*checklistReq.VerifiedBy)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error parsing verified by: ", err)
						return nil, err
					}
					verifiedBy = &parsedVerifiedBy
				}
				var isChecked string
				if checklistReq.IsChecked != nil {
					isChecked = *checklistReq.IsChecked
				} else {
					isChecked = "NO"
				}
				_, err := uc.EmployeeTaskChecklistRepository.UpdateEmployeeTaskChecklist(&entity.EmployeeTaskChecklist{
					ID:             parsedChecklistID,
					EmployeeTaskID: employeeTask.ID,
					Name:           checklistReq.Name,
					IsChecked:      isChecked,
					VerifiedBy:     verifiedBy,
				})
				if err != nil {
					uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error updating employee task checklist: ", err)
					return nil, err
				}
			}
		} else {
			_, err := uc.EmployeeTaskChecklistRepository.CreateEmployeeTaskChecklist(&entity.EmployeeTaskChecklist{
				EmployeeTaskID: employeeTask.ID,
				Name:           checklistReq.Name,
			})
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error creating employee task checklist: ", err)
				return nil, err
			}
		}
	}

	findById, err := uc.Repository.FindByID(employeeTask.ID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error finding employee task by id: ", err)
		return nil, err
	}
	if findById == nil {
		return nil, errors.New("employee task not found")
	}

	return uc.DTO.ConvertEntityToResponse(findById), nil
}

func (uc *EmployeeTaskUseCase) UpdateEmployeeTask(req *request.UpdateEmployeeTaskRequest) (*response.EmployeeTaskResponse, error) {
	parsedID, err := uuid.Parse(*req.ID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error parsing id: ", err)
		return nil, err
	}
	empTask, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error finding employee task by id: ", err)
		return nil, err
	}
	if empTask == nil {
		return nil, errors.New("employee task not found")
	}

	var templateTaskUUID *uuid.UUID
	if req.TemplateTaskID != nil && *req.TemplateTaskID != "" {
		parsedTemplateTaskID, err := uuid.Parse(*req.TemplateTaskID)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error parsing template task id: ", err)
			return nil, err
		}
		templateTask, err := uc.TemplateTaskRepository.FindByID(parsedTemplateTaskID)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error finding template task by id: ", err)
			return nil, err
		}
		if templateTask == nil {
			return nil, errors.New("template task not found")
		}

		templateTaskUUID = &parsedTemplateTaskID
	}

	var verifiedBy *uuid.UUID
	if req.VerifiedBy != nil && *req.VerifiedBy != "" {
		parsedVerifiedBy, err := uuid.Parse(*req.VerifiedBy)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error parsing verified by: ", err)
			return nil, err
		}
		verifiedBy = &parsedVerifiedBy
	}

	parsedEmployeeID, err := uuid.Parse(*req.EmployeeID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error parsing employee id: ", err)
		return nil, err
	}

	parsedStartDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error parsing start date: ", err)
		return nil, err
	}

	parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error parsing end date: ", err)
		return nil, err
	}

	employeeTask, err := uc.Repository.UpdateEmployeeTask(&entity.EmployeeTask{
		ID:             parsedID,
		CoverPath:      req.CoverPath,
		EmployeeID:     &parsedEmployeeID,
		TemplateTaskID: templateTaskUUID,
		Name:           req.Name,
		Priority:       entity.EmployeeTaskPriorityEnum(req.Priority),
		Description:    req.Description,
		StartDate:      parsedStartDate,
		EndDate:        parsedEndDate,
		VerifiedBy:     verifiedBy,
		IsDone:         req.IsDone,
		Proof:          req.ProofPath,
		Status:         entity.EmployeeTaskStatusEnum(req.Status),
		Kanban:         entity.EmployeeTaskKanbanEnum(req.Kanban),
		Notes:          req.Notes,
	})
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error creating employee task: ", err)
		return nil, err
	}

	// create employee task attachments
	for _, attachmentReq := range req.EmployeeTaskAttachments {
		_, err := uc.EmployeeTaskAttachmentRepository.CreateEmployeeTaskAttachment(&entity.EmployeeTaskAttachment{
			EmployeeTaskID: employeeTask.ID,
			Path:           attachmentReq.Path,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error creating employee task attachment: ", err)
			return nil, err
		}
	}

	// create employee task checklists
	var checklistIds []uuid.UUID
	for _, checklistReq := range req.EmployeeTaskChecklists {
		if checklistReq.ID != nil {
			parsedChecklistID, err := uuid.Parse(*checklistReq.ID)
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error parsing checklist id: ", err)
				return nil, err
			}
			exist, err := uc.EmployeeTaskChecklistRepository.FindByKeys(map[string]interface{}{
				"id": parsedChecklistID,
			})
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error finding checklist by id: ", err)
				return nil, err
			}
			if exist == nil {
				_, err := uc.EmployeeTaskChecklistRepository.CreateEmployeeTaskChecklist(&entity.EmployeeTaskChecklist{
					EmployeeTaskID: employeeTask.ID,
					Name:           checklistReq.Name,
				})
				if err != nil {
					uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error creating employee task checklist: ", err)
					return nil, err
				}
			} else {
				checklistIds = append(checklistIds, parsedChecklistID)
				var verifiedBy *uuid.UUID
				if checklistReq.VerifiedBy != nil {
					parsedVerifiedBy, err := uuid.Parse(*checklistReq.VerifiedBy)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error parsing verified by: ", err)
						return nil, err
					}
					verifiedBy = &parsedVerifiedBy
				}
				var isChecked string
				if checklistReq.IsChecked != nil {
					isChecked = *checklistReq.IsChecked
				} else {
					isChecked = "NO"
				}
				_, err := uc.EmployeeTaskChecklistRepository.UpdateEmployeeTaskChecklist(&entity.EmployeeTaskChecklist{
					ID:             parsedChecklistID,
					EmployeeTaskID: employeeTask.ID,
					Name:           checklistReq.Name,
					IsChecked:      isChecked,
					VerifiedBy:     verifiedBy,
				})
				if err != nil {
					uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error updating employee task checklist: ", err)
					return nil, err
				}
			}
		} else {
			_, err := uc.EmployeeTaskChecklistRepository.CreateEmployeeTaskChecklist(&entity.EmployeeTaskChecklist{
				EmployeeTaskID: employeeTask.ID,
				Name:           checklistReq.Name,
			})
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error creating employee task checklist: ", err)
				return nil, err
			}
		}
	}

	// delete employee task checklists
	err = uc.EmployeeTaskChecklistRepository.DeleteByEmployeeTaskIDAndNotInChecklistIDs(employeeTask.ID, checklistIds)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error deleting employee task checklists: ", err)
		return nil, err
	}

	if req.Kanban == string(entity.EMPLOYEE_TASK_KANBAN_ENUM_TODO) {
		etData, err := uc.Repository.FindByID(employeeTask.ID)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error finding employee task by id: ", err)
			return nil, err
		}
		if etData == nil {
			return nil, errors.New("employee task not found")
		}

		if len(etData.EmployeeTaskChecklists) > 0 {
			for _, checklist := range etData.EmployeeTaskChecklists {
				_, err := uc.EmployeeTaskChecklistRepository.UpdateEmployeeTaskChecklist(&entity.EmployeeTaskChecklist{
					ID:         checklist.ID,
					IsChecked:  "NO",
					VerifiedBy: nil,
				})
				if err != nil {
					uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error updating employee task checklist: ", err)
					return nil, err
				}
			}
		}
	}

	findById, err := uc.Repository.FindByID(employeeTask.ID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error finding employee task by id: ", err)
		return nil, err
	}
	if findById == nil {
		return nil, errors.New("employee task not found")
	}

	return uc.DTO.ConvertEntityToResponse(findById), nil
}

func (uc *EmployeeTaskUseCase) DeleteEmployeeTask(id uuid.UUID) error {
	parsedId, err := uuid.Parse(id.String())
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.DeleteEmployeeTaskUseCase] error parsing id: ", err)
		return err
	}
	exist, err := uc.Repository.FindByID(parsedId)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.DeleteEmployeeTaskUseCase] error finding employee task by id: ", err)
		return err
	}
	if exist == nil {
		return errors.New("employee task not found")
	}

	err = uc.Repository.DeleteEmployeeTask(exist)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.DeleteEmployeeTaskUseCase] error deleting employee task: ", err)
		return err
	}

	return nil
}

func (uc *EmployeeTaskUseCase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.EmployeeTaskResponse, int64, error) {
	employeeTasks, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.FindAllPaginated] error finding all employee tasks: ", err)
		return nil, 0, err
	}

	var responses []response.EmployeeTaskResponse
	for _, employeeTask := range *employeeTasks {
		responses = append(responses, *uc.DTO.ConvertEntityToResponse(&employeeTask))
	}

	return &responses, total, nil
}

func (uc *EmployeeTaskUseCase) FindByID(id uuid.UUID) (*response.EmployeeTaskResponse, error) {
	parsedId, err := uuid.Parse(id.String())
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.FindByID] error parsing id: ", err)
		return nil, err
	}
	employeeTask, err := uc.Repository.FindByID(parsedId)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.FindByID] error finding employee task by id: ", err)
		return nil, err
	}
	if employeeTask == nil {
		return nil, errors.New("employee task not found")
	}

	return uc.DTO.ConvertEntityToResponse(employeeTask), nil
}

func (uc *EmployeeTaskUseCase) CountByKanbanAndEmployeeID(kanban entity.EmployeeTaskKanbanEnum, employeeID uuid.UUID) (int64, error) {
	count, err := uc.Repository.CountByKeys(map[string]interface{}{
		"kanban":      kanban,
		"employee_id": employeeID,
	})
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CountByKanbanAndEmployeeID] error counting by kanban and employee id: ", err)
		return 0, err
	}

	return count, nil
}

func (uc *EmployeeTaskUseCase) FindAllByEmployeeID(employeeID uuid.UUID) (*response.EmployeeTaskKanbanResponse, error) {
	employeeTasks, err := uc.Repository.FindAllByEmployeeID(employeeID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.FindAllByEmployeeID] error finding all by employee id: ", err)
		return nil, err
	}

	formattedResponse := uc.DTO.ConvertEntitiesToKanbanResponse(*employeeTasks)

	return formattedResponse, nil
}

func (uc *EmployeeTaskUseCase) FindAllByEmployeeIDAndKanbanPaginated(employeeID uuid.UUID, kanban entity.EmployeeTaskKanbanEnum, page, pageSize int, search string, sort map[string]interface{}) (*[]response.EmployeeTaskResponse, int64, error) {
	employeeTasks, total, err := uc.Repository.FindAllByEmployeeIDAndKanbanPaginated(employeeID, kanban, page, pageSize, search, sort)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.FindAllByEmployeeIDAndKanbanPaginated] error finding all by employee id and kanban: ", err)
		return nil, 0, err
	}

	var responses []response.EmployeeTaskResponse
	for _, employeeTask := range *employeeTasks {
		responses = append(responses, *uc.DTO.ConvertEntityToResponse(&employeeTask))
	}

	return &responses, total, nil
}

func (uc *EmployeeTaskUseCase) UpdateEmployeeTaskOnly(req *request.UpdateEmployeeTaskOnlyRequest) (*response.EmployeeTaskResponse, error) {
	parsedID, err := uuid.Parse(*req.ID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskOnly] error parsing id: ", err)
		return nil, err
	}
	empTask, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskOnly] error finding employee task by id: ", err)
		return nil, err
	}
	if empTask == nil {
		return nil, errors.New("employee task not found")
	}

	parsedStartDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskOnly] error parsing start date: ", err)
		return nil, err
	}

	parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskOnly] error parsing end date: ", err)
		return nil, err
	}

	employeeTask, err := uc.Repository.UpdateEmployeeTask(&entity.EmployeeTask{
		ID:        parsedID,
		StartDate: parsedStartDate,
		EndDate:   parsedEndDate,
	})
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskOnly] error updating employee task: ", err)
		return nil, err
	}

	findById, err := uc.Repository.FindByID(employeeTask.ID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskOnly] error finding employee task by id: ", err)
		return nil, err
	}
	if findById == nil {
		return nil, errors.New("employee task not found")
	}

	return uc.DTO.ConvertEntityToResponse(findById), nil
}

func (uc *EmployeeTaskUseCase) CreateEmployeeTasksForRecruitment(req *request.CreateEmployeeTasksForRecruitment) error {
	templateTasks, err := uc.TemplateTaskRepository.FindAllByKeys(map[string]interface{}{
		"organization_type": req.OrganizationType,
		"status":            entity.TEMPLATE_TASK_STATUS_ENUM_ACTIVE,
	})
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTasksForRecruitment] error finding all template tasks: ", err)
		return err
	}

	parsedEmployeeID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTasksForRecruitment] error parsing employee id: ", err)
		return err
	}

	parsedJoinedDate, err := time.Parse("2006-01-02 15:04:05 -0700 MST", req.JoinedDate)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTasksForRecruitment] error parsing joined date: ", err)
		return err
	}

	_, err = uc.EmployeeHiringRepository.CreateEmployeeHiring(&entity.EmployeeHiring{
		EmployeeID: parsedEmployeeID,
		HiringDate: parsedJoinedDate,
	})
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTasksForRecruitment] error creating employee hiring: ", err)
	}

	// create employee tasks
	if len(*templateTasks) > 0 {
		for _, templateTask := range *templateTasks {
			empTaskExist, err := uc.Repository.FindByKeys(map[string]interface{}{
				"employee_id":      parsedEmployeeID,
				"template_task_id": templateTask.ID,
			})
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTasksForRecruitment] error finding employee task by keys: ", err)
				continue
			}
			if empTaskExist == nil {
				_, err = uc.Repository.CreateEmployeeTask(&entity.EmployeeTask{
					EmployeeID:     &parsedEmployeeID,
					TemplateTaskID: &templateTask.ID,
					StartDate:      parsedJoinedDate,
					EndDate:        parsedJoinedDate.AddDate(0, 0, *templateTask.DueDuration),
					CoverPath:      templateTask.CoverPath,
					Name:           templateTask.Name,
					Description:    templateTask.Description,
					Status:         entity.EMPLOYEE_TASK_STATUS_ENUM_ACTIVE,
					Kanban:         entity.EMPLOYEE_TASK_KANBAN_ENUM_TODO,
					Priority:       entity.EmployeeTaskPriorityEnum(templateTask.Priority),
					IsDone:         "NO",
					Source:         "ONBOARDING",
				})
				if err != nil {
					uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTasksForRecruitment] error creating employee task: ", err)
					continue
				}
			}
		}
	}

	return nil
}

func (uc *EmployeeTaskUseCase) CountKanbanProgressByEmployeeID(employeeID uuid.UUID) (*response.EmployeeTaskProgressResponse, error) {
	totalTask, err := uc.Repository.CountByKeys(map[string]interface{}{
		"employee_id": employeeID,
	})
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CountKanbanProgressByEmployeeID] error counting total task: ", err)
		return nil, err
	}

	toDo, err := uc.Repository.CountByKeys(map[string]interface{}{
		"employee_id": employeeID,
		"kanban":      entity.EMPLOYEE_TASK_KANBAN_ENUM_TODO,
	})
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CountKanbanProgressByEmployeeID] error counting to do: ", err)
		return nil, err
	}

	inProgress, err := uc.Repository.CountByKeys(map[string]interface{}{
		"employee_id": employeeID,
		"kanban":      entity.EPMLOYEE_TASK_KANBAN_ENUM_IN_PROGRESS,
	})
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CountKanbanProgressByEmployeeID] error counting in progress: ", err)
		return nil, err
	}

	needReview, err := uc.Repository.CountByKeys(map[string]interface{}{
		"employee_id": employeeID,
		"kanban":      entity.EMPLOYEE_TASK_KANBAN_ENUM_NEED_REVIEW,
	})
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CountKanbanProgressByEmployeeID] error counting need review: ", err)
		return nil, err
	}

	completed, err := uc.Repository.CountByKeys(map[string]interface{}{
		"employee_id": employeeID,
		"kanban":      entity.EMPLOYEE_TASK_KANBAN_ENUM_COMPLETED,
	})
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CountKanbanProgressByEmployeeID] error counting completed: ", err)
		return nil, err
	}

	return &response.EmployeeTaskProgressResponse{
		EmployeeID: employeeID,
		TotalTask:  int(totalTask),
		ToDo:       int(toDo),
		InProgress: int(inProgress),
		NeedReview: int(needReview),
		Completed:  int(completed),
	}, nil
}
