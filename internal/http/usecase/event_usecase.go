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

type IEventUseCase interface {
	CreateEvent(req *request.CreateEventRequest) (*response.EventResponse, error)
	UpdateEvent(req *request.UpdateEventRequest) (*response.EventResponse, error)
	DeleteEvent(id uuid.UUID) error
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
}

func NewEventUseCase(
	log *logrus.Logger,
	dto dto.IEventDTO,
	repository repository.IEventRepository,
	eventEmployeeRepository repository.IEventEmployeeRepository,
	templateTaskRepository repository.ITemplateTaskRepository,
	viper *viper.Viper,
	employeeTaskRepository repository.IEmployeeTaskRepository,
) IEventUseCase {
	return &EventUseCase{
		Log:                     log,
		DTO:                     dto,
		Repository:              repository,
		EventEmployeeRepository: eventEmployeeRepository,
		TemplateTaskRepository:  templateTaskRepository,
		Viper:                   viper,
		EmployeeTaskRepository:  employeeTaskRepository,
	}
}

func EventUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IEventUseCase {
	eventDTO := dto.EventDTOFactory(log, viper)
	eventRepository := repository.EventRepositoryFactory(log)
	eventEmployeeRepository := repository.EventEmployeeRepositoryFactory(log)
	templateTaskRepository := repository.TemplateTaskRepositoryFactory(log)
	employeeTaskRepository := repository.EmployeeTaskRepositoryFactory(log)
	return NewEventUseCase(
		log,
		eventDTO,
		eventRepository,
		eventEmployeeRepository,
		templateTaskRepository,
		viper,
		employeeTaskRepository,
	)
}

func (uc *EventUseCase) CreateEvent(req *request.CreateEventRequest) (*response.EventResponse, error) {
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

	event, err := uc.Repository.CreateEvent(&entity.Event{
		Name:           req.Name,
		TemplateTaskID: parsedTemplateTaskID,
		StartDate:      parsedStartDate,
		EndDate:        parsedEndDate,
		Description:    req.Description,
		Status:         entity.EventStatusEnum(req.Status),
	})
	if err != nil {
		return nil, err
	}

	// delete event employees
	err = uc.EventEmployeeRepository.DeleteByEventID(event.ID)
	if err != nil {
		return nil, err
	}

	// create event employee
	for _, eventEmployee := range req.EventEmployees {
		parsedEmployeeID, err := uuid.Parse(eventEmployee.EmployeeID)
		if err != nil {
			return nil, err
		}

		_, err = uc.EventEmployeeRepository.CreateEventEmployee(&entity.EventEmployee{
			EventID:    event.ID,
			EmployeeID: &parsedEmployeeID,
		})
		if err != nil {
			return nil, err
		}

		employeeTasks, err := uc.EmployeeTaskRepository.FindAllByEmployeeID(parsedEmployeeID)
		if err != nil {
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
					return nil, err
				}
			}
		}
	}

	findById, err := uc.Repository.FindByID(event.ID)
	if err != nil {
		return nil, err
	}
	if findById == nil {
		return nil, errors.New("event not found")
	}

	resp := uc.DTO.ConvertEntityToResponse(findById)
	return resp, nil
}

func (uc *EventUseCase) UpdateEvent(req *request.UpdateEventRequest) (*response.EventResponse, error) {
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
		return nil, err
	}

	// delete event employees
	err = uc.EventEmployeeRepository.DeleteByEventID(event.ID)
	if err != nil {
		return nil, err
	}

	// create event employee
	for _, eventEmployee := range req.EventEmployees {
		parsedEmployeeID, err := uuid.Parse(eventEmployee.EmployeeID)
		if err != nil {
			return nil, err
		}

		_, err = uc.EventEmployeeRepository.CreateEventEmployee(&entity.EventEmployee{
			EventID:    event.ID,
			EmployeeID: &parsedEmployeeID,
		})
		if err != nil {
			return nil, err
		}

		employeeTasks, err := uc.EmployeeTaskRepository.FindAllByEmployeeID(parsedEmployeeID)
		if err != nil {
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
					return nil, err
				}
			}
		}
	}

	findById, err := uc.Repository.FindByID(event.ID)
	if err != nil {
		return nil, err
	}
	if findById == nil {
		return nil, errors.New("event not found")
	}

	resp := uc.DTO.ConvertEntityToResponse(findById)
	return resp, nil
}

func (uc *EventUseCase) DeleteEvent(id uuid.UUID) error {
	exist, err := uc.Repository.FindByID(id)
	if err != nil {
		return err
	}
	if exist == nil {
		return errors.New("event not found")
	}

	err = uc.Repository.DeleteEvent(exist)
	if err != nil {
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
