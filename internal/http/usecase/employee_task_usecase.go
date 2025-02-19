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
	FindAllByEmployeeID(employeeID uuid.UUID) (*[]response.EmployeeTaskResponse, error)
}

type EmployeeTaskUseCase struct {
	Log                              *logrus.Logger
	DTO                              dto.IEmployeeTaskDTO
	Repository                       repository.IEmployeeTaskRepository
	Viper                            *viper.Viper
	TemplateTaskRepository           repository.ITemplateTaskRepository
	EmployeeTaskAttachmentRepository repository.IEmployeeTaskAttachmentRepository
	EmployeeTaskChecklistRepository  repository.IEmployeeTaskChecklistRepository
}

func NewEmployeeTaskUseCase(
	log *logrus.Logger,
	dto dto.IEmployeeTaskDTO,
	repository repository.IEmployeeTaskRepository,
	viper *viper.Viper,
	templateTaskRepository repository.ITemplateTaskRepository,
	etaRepo repository.IEmployeeTaskAttachmentRepository,
	etcRepo repository.IEmployeeTaskChecklistRepository,
) IEmployeeTaskUseCase {
	return &EmployeeTaskUseCase{
		Log:                              log,
		DTO:                              dto,
		Repository:                       repository,
		Viper:                            viper,
		TemplateTaskRepository:           templateTaskRepository,
		EmployeeTaskAttachmentRepository: etaRepo,
		EmployeeTaskChecklistRepository:  etcRepo,
	}
}

func EmployeeTaskUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IEmployeeTaskUseCase {
	etDTO := dto.EmployeeTaskDTOFactory(log, viper)
	repo := repository.EmployeeTaskRepositoryFactory(log)
	ttRepository := repository.TemplateTaskRepositoryFactory(log)
	etaRepo := repository.EmployeeTaskAttachmentRepositoryFactory(log)
	etcRepo := repository.EmployeeTaskChecklistRepositoryFactory(log)
	return NewEmployeeTaskUseCase(log, etDTO, repo, viper, ttRepository, etaRepo, etcRepo)
}

func (uc *EmployeeTaskUseCase) CreateEmployeeTask(req *request.CreateEmployeeTaskRequest) (*response.EmployeeTaskResponse, error) {
	var templateTaskUUID *uuid.UUID
	if req.TemplateTaskID != nil {
		parsedTemplateTaskID, err := uuid.Parse(*req.TemplateTaskID)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error parsing template task id: ", err)
			return nil, err
		}
		templateTask, err := uc.TemplateTaskRepository.FindByID(parsedTemplateTaskID)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error finding template task by id: ", err)
			return nil, err
		}
		if templateTask == nil {
			return nil, errors.New("template task not found")
		}

		templateTaskUUID = &parsedTemplateTaskID
	}

	parsedEmployeeID, err := uuid.Parse(*req.EmployeeID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error parsing employee id: ", err)
		return nil, err
	}

	parsedStartDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error parsing start date: ", err)
		return nil, err
	}

	parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error parsing end date: ", err)
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
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error creating employee task: ", err)
		return nil, err
	}

	// create employee task attachments
	for _, attachmentReq := range req.EmployeeTaskAttachments {
		_, err := uc.EmployeeTaskAttachmentRepository.CreateEmployeeTaskAttachment(&entity.EmployeeTaskAttachment{
			EmployeeTaskID: employeeTask.ID,
			Path:           attachmentReq.Path,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error creating employee task attachment: ", err)
			return nil, err
		}
	}

	// delete employee task checklists
	err = uc.EmployeeTaskChecklistRepository.DeleteByEmployeeTaskID(employeeTask.ID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error deleting employee task checklists: ", err)
		return nil, err
	}

	// create employee task checklists
	for _, checklistReq := range req.EmployeeTaskChecklists {
		if checklistReq.ID != nil {
			parsedChecklistID, err := uuid.Parse(*checklistReq.ID)
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error parsing checklist id: ", err)
				return nil, err
			}
			exist, err := uc.EmployeeTaskChecklistRepository.FindByKeys(map[string]interface{}{
				"id": parsedChecklistID,
			})
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error finding checklist by id: ", err)
				return nil, err
			}
			if exist == nil {
				_, err := uc.EmployeeTaskChecklistRepository.CreateEmployeeTaskChecklist(&entity.EmployeeTaskChecklist{
					EmployeeTaskID: employeeTask.ID,
					Name:           checklistReq.Name,
				})
				if err != nil {
					uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error creating employee task checklist: ", err)
					return nil, err
				}
			} else {
				var verifiedBy *uuid.UUID
				if checklistReq.VerifiedBy != nil {
					parsedVerifiedBy, err := uuid.Parse(*checklistReq.VerifiedBy)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error parsing verified by: ", err)
						return nil, err
					}
					verifiedBy = &parsedVerifiedBy
				}
				_, err := uc.EmployeeTaskChecklistRepository.UpdateEmployeeTaskChecklist(&entity.EmployeeTaskChecklist{
					ID:             parsedChecklistID,
					EmployeeTaskID: employeeTask.ID,
					Name:           checklistReq.Name,
					IsChecked:      *checklistReq.IsChecked,
					VerifiedBy:     verifiedBy,
				})
				if err != nil {
					uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error updating employee task checklist: ", err)
					return nil, err
				}
			}
		} else {
			_, err := uc.EmployeeTaskChecklistRepository.CreateEmployeeTaskChecklist(&entity.EmployeeTaskChecklist{
				EmployeeTaskID: employeeTask.ID,
				Name:           checklistReq.Name,
			})
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error creating employee task checklist: ", err)
				return nil, err
			}
		}
	}

	findById, err := uc.Repository.FindByID(employeeTask.ID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error finding employee task by id: ", err)
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
	if req.TemplateTaskID != nil {
		parsedTemplateTaskID, err := uuid.Parse(*req.TemplateTaskID)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error parsing template task id: ", err)
			return nil, err
		}
		templateTask, err := uc.TemplateTaskRepository.FindByID(parsedTemplateTaskID)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error finding template task by id: ", err)
			return nil, err
		}
		if templateTask == nil {
			return nil, errors.New("template task not found")
		}

		templateTaskUUID = &parsedTemplateTaskID
	}

	var verifiedBy *uuid.UUID
	if req.VerifiedBy != nil {
		parsedVerifiedBy, err := uuid.Parse(*req.VerifiedBy)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error parsing verified by: ", err)
			return nil, err
		}
		verifiedBy = &parsedVerifiedBy
	}

	parsedEmployeeID, err := uuid.Parse(*req.EmployeeID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error parsing employee id: ", err)
		return nil, err
	}

	parsedStartDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error parsing start date: ", err)
		return nil, err
	}

	parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error parsing end date: ", err)
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
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error creating employee task: ", err)
		return nil, err
	}

	// create employee task attachments
	for _, attachmentReq := range req.EmployeeTaskAttachments {
		_, err := uc.EmployeeTaskAttachmentRepository.CreateEmployeeTaskAttachment(&entity.EmployeeTaskAttachment{
			EmployeeTaskID: employeeTask.ID,
			Path:           attachmentReq.Path,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error creating employee task attachment: ", err)
			return nil, err
		}
	}

	// create employee task checklists
	var checklistIds []uuid.UUID
	for _, checklistReq := range req.EmployeeTaskChecklists {
		if checklistReq.ID != nil {
			parsedChecklistID, err := uuid.Parse(*checklistReq.ID)
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error parsing checklist id: ", err)
				return nil, err
			}
			exist, err := uc.EmployeeTaskChecklistRepository.FindByKeys(map[string]interface{}{
				"id": parsedChecklistID,
			})
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error finding checklist by id: ", err)
				return nil, err
			}
			if exist == nil {
				_, err := uc.EmployeeTaskChecklistRepository.CreateEmployeeTaskChecklist(&entity.EmployeeTaskChecklist{
					EmployeeTaskID: employeeTask.ID,
					Name:           checklistReq.Name,
				})
				if err != nil {
					uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error creating employee task checklist: ", err)
					return nil, err
				}
			} else {
				checklistIds = append(checklistIds, parsedChecklistID)
				var verifiedBy *uuid.UUID
				if checklistReq.VerifiedBy != nil {
					parsedVerifiedBy, err := uuid.Parse(*checklistReq.VerifiedBy)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error parsing verified by: ", err)
						return nil, err
					}
					verifiedBy = &parsedVerifiedBy
				}
				_, err := uc.EmployeeTaskChecklistRepository.UpdateEmployeeTaskChecklist(&entity.EmployeeTaskChecklist{
					ID:             parsedChecklistID,
					EmployeeTaskID: employeeTask.ID,
					Name:           checklistReq.Name,
					IsChecked:      *checklistReq.IsChecked,
					VerifiedBy:     verifiedBy,
				})
				if err != nil {
					uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error updating employee task checklist: ", err)
					return nil, err
				}
			}
		} else {
			_, err := uc.EmployeeTaskChecklistRepository.CreateEmployeeTaskChecklist(&entity.EmployeeTaskChecklist{
				EmployeeTaskID: employeeTask.ID,
				Name:           checklistReq.Name,
			})
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error creating employee task checklist: ", err)
				return nil, err
			}
		}
	}

	// delete employee task checklists
	err = uc.EmployeeTaskChecklistRepository.DeleteByEmployeeTaskIDAndNotInChecklistIDs(employeeTask.ID, checklistIds)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error deleting employee task checklists: ", err)
		return nil, err
	}

	findById, err := uc.Repository.FindByID(employeeTask.ID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error finding employee task by id: ", err)
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

func (uc *EmployeeTaskUseCase) FindAllByEmployeeID(employeeID uuid.UUID) (*[]response.EmployeeTaskResponse, error) {
	employeeTasks, err := uc.Repository.FindAllByEmployeeID(employeeID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.FindAllByEmployeeID] error finding all by employee id: ", err)
		return nil, err
	}

	var responses []response.EmployeeTaskResponse
	for _, employeeTask := range *employeeTasks {
		responses = append(responses, *uc.DTO.ConvertEntityToResponse(&employeeTask))
	}

	return &responses, nil
}
