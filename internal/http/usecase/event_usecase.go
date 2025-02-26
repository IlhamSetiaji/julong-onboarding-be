package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/dto"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type IEventUseCase interface {
	CreateEvent(ctx context.Context, req *request.CreateEventRequest) (*response.EventResponse, error)
	UpdateEvent(ctx context.Context, req *request.UpdateEventRequest) (*response.EventResponse, error)
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.EventResponse, int64, error)
	FindByID(id uuid.UUID) (*response.EventResponse, error)
}

type EventUseCase struct {
	Log                     *logrus.Logger
	DTO                     dto.IEventDTO
	Repository              repository.IEventRepository
	EventEmployeeRepository repository.IEventEmployeeRepository
	TemplateTaskRepository  repository.ITemplateTaskRepository
	Viper                   *viper.Viper
	EmployeeTaskRepository  repository.IEmployeeTaskRepository
	DB                      *gorm.DB
}

func NewEventUseCase(
	log *logrus.Logger,
	dto dto.IEventDTO,
	repository repository.IEventRepository,
	eventEmployeeRepository repository.IEventEmployeeRepository,
	templateTaskRepository repository.ITemplateTaskRepository,
	viper *viper.Viper,
	employeeTaskRepository repository.IEmployeeTaskRepository,
	db *gorm.DB,
) IEventUseCase {
	return &EventUseCase{
		Log:                     log,
		DTO:                     dto,
		Repository:              repository,
		EventEmployeeRepository: eventEmployeeRepository,
		TemplateTaskRepository:  templateTaskRepository,
		Viper:                   viper,
		EmployeeTaskRepository:  employeeTaskRepository,
		DB:                      db,
	}
}

func EventUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IEventUseCase {
	eventDTO := dto.EventDTOFactory(log, viper)
	eventRepository := repository.EventRepositoryFactory(log)
	eventEmployeeRepository := repository.EventEmployeeRepositoryFactory(log)
	templateTaskRepository := repository.TemplateTaskRepositoryFactory(log)
	employeeTaskRepository := repository.EmployeeTaskRepositoryFactory(log)
	db := config.NewDatabase()
	return NewEventUseCase(
		log,
		eventDTO,
		eventRepository,
		eventEmployeeRepository,
		templateTaskRepository,
		viper,
		employeeTaskRepository,
		db,
	)
}

func (uc *EventUseCase) CreateEvent(ctx context.Context, req *request.CreateEventRequest) (*response.EventResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		uc.Log.Errorf("[EventUseCase.CreateEvent] error when starting transaction: %s", tx.Error)
		return nil, errors.New("[EventUseCase.CreateEvent] error when starting transaction: " + tx.Error.Error())
	}
	defer tx.Rollback()

	parsedTemplateTaskID, err := uuid.Parse(req.TemplateTaskID)
	if err != nil {
		return nil, err
	}
	templateTask, err := uc.TemplateTaskRepository.FindByID(parsedTemplateTaskID)
	if err != nil {
		return nil, err
	}
	if templateTask == nil {
		return nil, errors.New("template task not found")
	}

	parsedStartDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, err
	}
	parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, err
	}

	var status entity.EventStatusEnum
	if parsedStartDate.After(time.Now()) {
		status = entity.EVENT_STATUS_ENUM_UPCOMING
	} else if parsedStartDate.Before(time.Now()) && parsedEndDate.After(time.Now()) {
		status = entity.EVENT_STATUS_ENUM_ONGOING
	} else {
		status = entity.EVENT_STATUS_ENUM_FINISHED
	}

	event, err := uc.Repository.CreateEvent(&entity.Event{
		Name:           req.Name,
		TemplateTaskID: parsedTemplateTaskID,
		StartDate:      parsedStartDate,
		EndDate:        parsedEndDate,
		Description:    req.Description,
		Status:         status,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// delete event employees
	err = uc.EventEmployeeRepository.DeleteByEventID(event.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// create event employee
	for _, eventEmployee := range req.EventEmployees {
		parsedEmployeeID, err := uuid.Parse(eventEmployee.EmployeeID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		_, err = uc.EventEmployeeRepository.CreateEventEmployee(&entity.EventEmployee{
			EventID:    event.ID,
			EmployeeID: &parsedEmployeeID,
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		employeeTasks, err := uc.EmployeeTaskRepository.FindAllByEmployeeID(parsedEmployeeID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if len(*employeeTasks) > 0 {
			for _, employeeTask := range *employeeTasks {
				_, err = uc.EmployeeTaskRepository.UpdateEmployeeTask(&entity.EmployeeTask{
					ID:        employeeTask.ID,
					StartDate: parsedStartDate,
					EndDate:   parsedEndDate,
				})
				if err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}

		empTaskExist, err := uc.EmployeeTaskRepository.FindByKeys(map[string]interface{}{
			"employee_id":      parsedEmployeeID,
			"template_task_id": parsedTemplateTaskID,
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if empTaskExist == nil {
			_, err = uc.EmployeeTaskRepository.CreateEmployeeTask(&entity.EmployeeTask{
				EmployeeID:     &parsedEmployeeID,
				TemplateTaskID: &parsedTemplateTaskID,
				StartDate:      parsedStartDate,
				EndDate:        parsedEndDate,
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
				tx.Rollback()
				return nil, err
			}
		}
	}

	findById, err := uc.Repository.FindByID(event.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if findById == nil {
		tx.Rollback()
		return nil, errors.New("event not found")
	}

	resp := uc.DTO.ConvertEntityToResponse(findById)
	return resp, nil
}

func (uc *EventUseCase) UpdateEvent(ctx context.Context, req *request.UpdateEventRequest) (*response.EventResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		uc.Log.Errorf("[EventUseCase.UpdateEvent] error when starting transaction: %s", tx.Error)
		return nil, errors.New("[EventUseCase.UpdateEvent] error when starting transaction: " + tx.Error.Error())
	}
	defer tx.Rollback()

	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, err
	}
	exist, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		return nil, err
	}
	if exist == nil {
		return nil, errors.New("event not found")
	}

	parsedTemplateTaskID, err := uuid.Parse(req.TemplateTaskID)
	if err != nil {
		return nil, err
	}
	templateTask, err := uc.TemplateTaskRepository.FindByID(parsedTemplateTaskID)
	if err != nil {
		return nil, err
	}
	if templateTask == nil {
		return nil, errors.New("template task not found")
	}

	parsedStartDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, err
	}
	parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, err
	}

	event, err := uc.Repository.UpdateEvent(&entity.Event{
		ID:             parsedID,
		Name:           req.Name,
		TemplateTaskID: parsedTemplateTaskID,
		StartDate:      parsedStartDate,
		EndDate:        parsedEndDate,
		Description:    req.Description,
		Status:         entity.EventStatusEnum(req.Status),
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// delete event employees
	err = uc.EventEmployeeRepository.DeleteByEventID(event.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// create event employee
	for _, eventEmployee := range req.EventEmployees {
		parsedEmployeeID, err := uuid.Parse(eventEmployee.EmployeeID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		_, err = uc.EventEmployeeRepository.CreateEventEmployee(&entity.EventEmployee{
			EventID:    event.ID,
			EmployeeID: &parsedEmployeeID,
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		employeeTasks, err := uc.EmployeeTaskRepository.FindAllByEmployeeID(parsedEmployeeID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if len(*employeeTasks) > 0 {
			for _, employeeTask := range *employeeTasks {
				_, err = uc.EmployeeTaskRepository.UpdateEmployeeTask(&entity.EmployeeTask{
					ID:        employeeTask.ID,
					StartDate: parsedStartDate,
					EndDate:   parsedEndDate,
				})
				if err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}

		empTaskExist, err := uc.EmployeeTaskRepository.FindByKeys(map[string]interface{}{
			"employee_id":      parsedEmployeeID,
			"template_task_id": parsedTemplateTaskID,
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if empTaskExist == nil {
			_, err = uc.EmployeeTaskRepository.CreateEmployeeTask(&entity.EmployeeTask{
				EmployeeID:     &parsedEmployeeID,
				TemplateTaskID: &parsedTemplateTaskID,
				StartDate:      parsedStartDate,
				EndDate:        parsedEndDate,
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
				tx.Rollback()
				return nil, err
			}
		}
	}

	findById, err := uc.Repository.FindByID(event.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if findById == nil {
		tx.Rollback()
		return nil, errors.New("event not found")
	}

	resp := uc.DTO.ConvertEntityToResponse(findById)
	return resp, nil
}

func (uc *EventUseCase) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	tx := uc.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		uc.Log.Errorf("[EventUseCase.DeleteEvent] error when starting transaction: %s", tx.Error)
		return errors.New("[EventUseCase.DeleteEvent] error when starting transaction: " + tx.Error.Error())
	}
	defer tx.Rollback()

	exist, err := uc.Repository.FindByID(id)
	if err != nil {
		return err
	}
	if exist == nil {
		return errors.New("event not found")
	}

	err = uc.Repository.DeleteEvent(exist)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (uc *EventUseCase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.EventResponse, int64, error) {
	events, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		return nil, 0, err
	}

	var responses []response.EventResponse
	for _, event := range *events {
		resp := uc.DTO.ConvertEntityToResponse(&event)
		responses = append(responses, *resp)
	}

	return &responses, total, nil
}

func (uc *EventUseCase) FindByID(id uuid.UUID) (*response.EventResponse, error) {
	event, err := uc.Repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, errors.New("event not found")
	}

	resp := uc.DTO.ConvertEntityToResponse(event)
	return resp, nil
}
