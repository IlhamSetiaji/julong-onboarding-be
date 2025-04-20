package usecase

import (
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/dto"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/service"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IEmployeeTaskUseCase interface {
	CreateEmployeeTask(req *request.CreateEmployeeTaskRequest) (*response.EmployeeTaskResponse, error)
	CreateEmployeeTaskMidsuit(req *request.CreateEmployeeTaskMidsuitRequest) (*response.EmployeeTaskResponse, error)
	UpdateEmployeeTask(req *request.UpdateEmployeeTaskRequest) (*response.EmployeeTaskResponse, error)
	UpdateEmployeeTaskMidsuit(req *request.UpdateEmployeeTaskMidsuitRequest) (*response.EmployeeTaskResponse, error)
	DeleteEmployeeTask(id uuid.UUID) error
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.EmployeeTaskResponse, int64, error)
	FindAllPaginatedByEmployeeID(employeeID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}) (*[]response.EmployeeTaskResponse, int64, error)
	FindByID(id uuid.UUID) (*response.EmployeeTaskResponse, error)
	CountByKanbanAndEmployeeID(kanban entity.EmployeeTaskKanbanEnum, employeeID uuid.UUID) (int64, error)
	FindAllByEmployeeID(employeeID uuid.UUID) (*response.EmployeeTaskKanbanResponse, error)
	FindAllByEmployeeIDAndKanbanPaginated(employeeID uuid.UUID, kanban entity.EmployeeTaskKanbanEnum, page, pageSize int, search string, sort map[string]interface{}) (*[]response.EmployeeTaskResponse, int64, error)
	UpdateEmployeeTaskOnly(req *request.UpdateEmployeeTaskOnlyRequest) (*response.EmployeeTaskResponse, error)
	CreateEmployeeTasksForRecruitment(req *request.CreateEmployeeTasksForRecruitment) error
	CountKanbanProgressByEmployeeID(employeeID uuid.UUID) (*response.EmployeeTaskProgressResponse, error)
	FindByIDForResponse(id uuid.UUID) (*response.EmployeeTaskResponse, error)
	FindAllPaginatedSurvey(page, pageSize int, search string, sort map[string]interface{}) (*[]response.EmployeeTaskResponse, int64, error)
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
	SurveyTemplateRepository         repository.ISurveyTemplateRepository
	MidsuitService                   service.IMidsuitService
	EmployeeMessage                  messaging.IEmployeeMessage
	OrganizationMessage              messaging.IOrganizationMessage
	JobPlafonMessage                 messaging.IJobPlafonMessage
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
	stRepo repository.ISurveyTemplateRepository,
	midsuitService service.IMidsuitService,
	employeeMessage messaging.IEmployeeMessage,
	organizationMessage messaging.IOrganizationMessage,
	jobPlafonMessage messaging.IJobPlafonMessage,
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
		SurveyTemplateRepository:         stRepo,
		MidsuitService:                   midsuitService,
		EmployeeMessage:                  employeeMessage,
		OrganizationMessage:              organizationMessage,
		JobPlafonMessage:                 jobPlafonMessage,
	}
}

func EmployeeTaskUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IEmployeeTaskUseCase {
	etDTO := dto.EmployeeTaskDTOFactory(log, viper)
	repo := repository.EmployeeTaskRepositoryFactory(log)
	ttRepository := repository.TemplateTaskRepositoryFactory(log)
	etaRepo := repository.EmployeeTaskAttachmentRepositoryFactory(log)
	etcRepo := repository.EmployeeTaskChecklistRepositoryFactory(log)
	ehRepo := repository.EmployeeHiringRepositoryFactory(log)
	stRepo := repository.SurveyTemplateRepositoryFactory(log)
	midsuitService := service.MidsuitServiceFactory(viper, log)
	employeeMessage := messaging.EmployeeMessageFactory(log)
	organizationMessage := messaging.OrganizationMessageFactory(log)
	jobPlafonMessage := messaging.JobPlafonMessageFactory(log)
	return NewEmployeeTaskUseCase(log, etDTO, repo, viper, ttRepository, etaRepo, etcRepo, ehRepo, stRepo, midsuitService, employeeMessage, organizationMessage, jobPlafonMessage)
}

func (uc *EmployeeTaskUseCase) CreateEmployeeTask(req *request.CreateEmployeeTaskRequest) (*response.EmployeeTaskResponse, error) {
	var templateTaskUUID *uuid.UUID
	if req.TemplateTaskID != nil && *req.TemplateTaskID != "" {
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

	var surveyTemplateUUID *uuid.UUID
	if req.SurveyTemplateID != nil && *req.SurveyTemplateID != "" {
		parsedSurveyTemplateID, err := uuid.Parse(*req.SurveyTemplateID)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error parsing survey template id: ", err)
			return nil, err
		}
		surveyTemplate, err := uc.SurveyTemplateRepository.FindByKeys(map[string]interface{}{
			"id": parsedSurveyTemplateID,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error finding survey template by id: ", err)
			return nil, err
		}
		if surveyTemplate == nil {
			return nil, errors.New("survey template not found")
		}

		surveyTemplateUUID = &parsedSurveyTemplateID
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

	// post to midsuit
	var midsuitID string
	if uc.Viper.GetString("midsuit.sync") == "ACTIVE" {
		empResp, err := uc.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
			ID: *req.EmployeeID,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error sending find employee by id message: ", err)
			return nil, err
		}
		if empResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] employee not found in midsuit")
			return nil, errors.New("employee not found in midsuit")
		}

		orgResp, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
			ID: empResp.OrganizationID.String(),
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error sending find organization by id message: ", err)
			return nil, err
		}
		if orgResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] organization not found in midsuit")
			return nil, errors.New("organization not found in midsuit")
		}

		jobId := empResp.EmployeeJob["job_id"].(string)
		jobResp, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
			ID: jobId,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error sending find job by id message: ", err)
			return nil, err
		}
		if jobResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] job not found in midsuit")
			return nil, errors.New("job not found in midsuit")
		}

		jobLevelId := empResp.EmployeeJob["job_level_id"].(string)
		jobLevelResp, err := uc.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
			ID: jobLevelId,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error sending find job level by id message: ", err)
			return nil, err
		}
		if jobLevelResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] job level not found in midsuit")
			return nil, errors.New("job level not found in midsuit")
		}

		orgStructureId := empResp.EmployeeJob["organization_structure_id"].(string)
		orgStructureResp, err := uc.OrganizationMessage.SendFindOrganizationStructureByIDMessage(request.SendFindOrganizationStructureByIDMessageRequest{
			ID: orgStructureId,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error sending find organization structure by id message: ", err)
			return nil, err
		}
		if orgStructureResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] organization structure not found in midsuit")
			return nil, errors.New("organization structure not found in midsuit")
		}

		midsuitPayload := &request.SyncEmployeeTaskMidsuitRequest{
			AdOrgId: request.AdOrgId{
				ID: func() int {
					id, err := strconv.Atoi(orgResp.MidsuitID)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error converting orgResp.MidsuitID to int: ", err)
						return 0 // or handle the error appropriately
					}
					return id
				}(),
				// ID: 1000024,
			},
			Name: req.Name,
			Category: request.TaskCategory{
				Identifier: "Onboarding",
				ModelName:  "ad_ref_list",
			},
			StartDate: parsedStartDate.String(),
			EndDate:   parsedEndDate.String(),
			HCEmployeeID: request.HcEmployeeId{
				ID: func() int {
					id, err := strconv.Atoi(empResp.MidsuitID)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error converting empResp.MidsuitID to int: ", err)
						return 0 // or handle the error appropriately
					}
					return id
				}(),
				// ID: 1000812,
			},
			HCJobID: request.HcJobId{
				ID: func() int {
					id, err := strconv.Atoi(jobResp.MidsuitID)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error converting jobResp.MidsuitID to int: ", err)
						return 0 // or handle the error appropriately
					}
					return id
				}(),
				// ID: 1000472,
			},
			HCJobLevelID: request.HcJobLevelId{
				ID: func() int {
					id, err := strconv.Atoi(jobLevelResp.MidsuitID)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error converting jobLevelResp.MidsuitID to int: ", err)
						return 0 // or handle the error appropriately
					}
					return id
				}(),
				// ID: 1000095,
			},
			HCOrgID: request.HcOrgId{
				ID: func() int {
					id, err := strconv.Atoi(orgStructureResp.MidsuitID)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error converting orgStructureResp.MidsuitID to int: ", err)
						return 0 // or handle the error appropriately
					}
					return id
				}(),
				// ID: 1000622,
			},
		}
		authResp, err := uc.MidsuitService.AuthOneStep()
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}

		midsuitEmpTask, err := uc.MidsuitService.SyncEmployeeTaskMidsuit(*midsuitPayload, authResp.Token)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error syncing employee task to midsuit: ", err)
			return nil, err
		}

		midsuitID = *midsuitEmpTask
	}

	employeeTask, err := uc.Repository.CreateEmployeeTask(&entity.EmployeeTask{
		CoverPath:        req.CoverPath,
		EmployeeID:       &parsedEmployeeID,
		TemplateTaskID:   templateTaskUUID,
		SurveyTemplateID: surveyTemplateUUID,
		Name:             req.Name,
		Priority:         entity.EmployeeTaskPriorityEnum(req.Priority),
		Description:      req.Description,
		StartDate:        parsedStartDate,
		EndDate:          parsedEndDate,
		Source:           "ONBOARDING",
		MidsuitID:        &midsuitID,
	})
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error creating employee task: ", err)
		return nil, err
	}

	// create employee task attachments
	for _, attachmentReq := range req.EmployeeTaskAttachments {
		if uc.Viper.GetString("midsuit.sync") == "ACTIVE" {
			// Read the file from the given path
			fileContent, err := os.ReadFile(attachmentReq.Path)
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error reading file: ", err)
				return nil, err
			}

			// Extract the file name from the path
			fileName := filepath.Base(attachmentReq.Path)

			// Encode the file content to base64
			encodedData := base64.StdEncoding.EncodeToString(fileContent)

			// Create the payload
			midsuitAttachmentPayload := &request.SyncEmployeeTaskAttachmentMidsuitRequest{
				Name: fileName,
				Data: encodedData,
			}

			// Log the payload for debugging
			// uc.Log.Info("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] midsuit attachment payload: ", midsuitAttachmentPayload)

			// Sync to midsuit
			authResp, err := uc.MidsuitService.AuthOneStep()
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] " + err.Error())
				return nil, err
			}

			midsuitIDInt, err := strconv.Atoi(midsuitID)
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error converting midsuitID to int: ", err)
				return nil, err
			}
			_, err = uc.MidsuitService.SyncEmployeeTaskAttachmentMidsuit(midsuitIDInt, *midsuitAttachmentPayload, authResp.Token)
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error syncing employee task attachment to midsuit: ", err)
				return nil, err
			}
		}

		_, err = uc.EmployeeTaskAttachmentRepository.CreateEmployeeTaskAttachment(&entity.EmployeeTaskAttachment{
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
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error deleting employee task checklists: ", err)
		return nil, err
	}

	// create employee task checklists
	for _, checklistReq := range req.EmployeeTaskChecklists {
		// sync emp task checklist midsuit
		if uc.Viper.GetString("midsuit.sync") == "ACTIVE" {
			empResp, err := uc.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
				ID: *req.EmployeeID,
			})
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error sending find employee by id message: ", err)
				return nil, err
			}
			if empResp == nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] employee not found in midsuit")
				return nil, errors.New("employee not found in midsuit")
			}

			orgResp, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
				ID: empResp.OrganizationID.String(),
			})
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error sending find organization by id message: ", err)
				return nil, err
			}
			if orgResp == nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] organization not found in midsuit")
				return nil, errors.New("organization not found in midsuit")
			}

			midsuitChecklistPayload := &request.SyncEmployeeTaskChecklistMidsuitRequest{
				AdOrgId: request.AdOrgId{
					ID: func() int {
						id, err := strconv.Atoi(orgResp.MidsuitID)
						if err != nil {
							uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error converting orgResp.MidsuitID to int: ", err)
							return 0 // or handle the error appropriately
						}
						return id
					}(),
				},
				Name:      checklistReq.Name,
				IsChecked: false,
				HCTaskID: request.HCTaskID{
					ID: func() int {
						id, err := strconv.Atoi(midsuitID)
						if err != nil {
							uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error converting employeeTask.MidsuitID to int: ", err)
							return 0 // or handle the error appropriately
						}
						return id
					}(),
				},
				HCEmployeeID: request.HcEmployeeId{
					ID: func() int {
						id, err := strconv.Atoi(empResp.MidsuitID)
						if err != nil {
							uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error converting empResp.MidsuitID to int: ", err)
							return 0 // or handle the error appropriately
						}
						return id
					}(),
				},
				ModelName: "hc_taskchecklist",
			}

			authResp, err := uc.MidsuitService.AuthOneStep()
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return nil, err
			}

			_, err = uc.MidsuitService.SyncEmployeeTaskChecklistMidsuit(*midsuitChecklistPayload, authResp.Token)
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error syncing employee task checklist to midsuit: ", err)
				return nil, err
			}
		}

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

func (uc *EmployeeTaskUseCase) CreateEmployeeTaskMidsuit(req *request.CreateEmployeeTaskMidsuitRequest) (*response.EmployeeTaskResponse, error) {
	// Implementation for CreateEmployeeTaskMidsuit
	var templateTaskUUID *uuid.UUID
	if req.TemplateTaskID != nil && *req.TemplateTaskID != "" {
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

	var surveyTemplateUUID *uuid.UUID
	if req.SurveyTemplateID != nil && *req.SurveyTemplateID != "" {
		parsedSurveyTemplateID, err := uuid.Parse(*req.SurveyTemplateID)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error parsing survey template id: ", err)
			return nil, err
		}
		surveyTemplate, err := uc.SurveyTemplateRepository.FindByKeys(map[string]interface{}{
			"id": parsedSurveyTemplateID,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error finding survey template by id: ", err)
			return nil, err
		}
		if surveyTemplate == nil {
			return nil, errors.New("survey template not found")
		}

		surveyTemplateUUID = &parsedSurveyTemplateID
	}

	empRespMessage, err := uc.EmployeeMessage.SendFindEmployeeByMidsuitIDMessage(*req.EmployeeMidsuitID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error sending find employee by midsuit id message: ", err)
		return nil, err
	}
	if empRespMessage == nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] employee not found in midsuit")
		return nil, errors.New("employee not found in midsuit")
	}

	parsedEmployeeID := empRespMessage.ID

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

	// post to midsuit
	var midsuitID string
	if uc.Viper.GetString("midsuit.sync") == "ACTIVE" {
		empResp, err := uc.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
			ID: parsedEmployeeID.String(),
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error sending find employee by id message: ", err)
			return nil, err
		}
		if empResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] employee not found in midsuit")
			return nil, errors.New("employee not found in midsuit")
		}

		orgResp, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
			ID: empResp.OrganizationID.String(),
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error sending find organization by id message: ", err)
			return nil, err
		}
		if orgResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] organization not found in midsuit")
			return nil, errors.New("organization not found in midsuit")
		}

		jobId := empResp.EmployeeJob["job_id"].(string)
		jobResp, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
			ID: jobId,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error sending find job by id message: ", err)
			return nil, err
		}
		if jobResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] job not found in midsuit")
			return nil, errors.New("job not found in midsuit")
		}

		jobLevelId := empResp.EmployeeJob["job_level_id"].(string)
		jobLevelResp, err := uc.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
			ID: jobLevelId,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error sending find job level by id message: ", err)
			return nil, err
		}
		if jobLevelResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] job level not found in midsuit")
			return nil, errors.New("job level not found in midsuit")
		}

		orgStructureId := empResp.EmployeeJob["organization_structure_id"].(string)
		orgStructureResp, err := uc.OrganizationMessage.SendFindOrganizationStructureByIDMessage(request.SendFindOrganizationStructureByIDMessageRequest{
			ID: orgStructureId,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error sending find organization structure by id message: ", err)
			return nil, err
		}
		if orgStructureResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] organization structure not found in midsuit")
			return nil, errors.New("organization structure not found in midsuit")
		}

		midsuitPayload := &request.SyncEmployeeTaskMidsuitRequest{
			AdOrgId: request.AdOrgId{
				ID: func() int {
					id, err := strconv.Atoi(orgResp.MidsuitID)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error converting orgResp.MidsuitID to int: ", err)
						return 0 // or handle the error appropriately
					}
					return id
				}(),
				// ID: 1000024,
			},
			Name: req.Name,
			Category: request.TaskCategory{
				Identifier: "Onboarding",
				ModelName:  "ad_ref_list",
			},
			StartDate: parsedStartDate.String(),
			EndDate:   parsedEndDate.String(),
			HCEmployeeID: request.HcEmployeeId{
				ID: func() int {
					id, err := strconv.Atoi(empResp.MidsuitID)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error converting empResp.MidsuitID to int: ", err)
						return 0 // or handle the error appropriately
					}
					return id
				}(),
				// ID: 1000812,
			},
			HCJobID: request.HcJobId{
				ID: func() int {
					id, err := strconv.Atoi(jobResp.MidsuitID)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error converting jobResp.MidsuitID to int: ", err)
						return 0 // or handle the error appropriately
					}
					return id
				}(),
				// ID: 1000472,
			},
			HCJobLevelID: request.HcJobLevelId{
				ID: func() int {
					id, err := strconv.Atoi(jobLevelResp.MidsuitID)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error converting jobLevelResp.MidsuitID to int: ", err)
						return 0 // or handle the error appropriately
					}
					return id
				}(),
				// ID: 1000095,
			},
			HCOrgID: request.HcOrgId{
				ID: func() int {
					id, err := strconv.Atoi(orgStructureResp.MidsuitID)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error converting orgStructureResp.MidsuitID to int: ", err)
						return 0 // or handle the error appropriately
					}
					return id
				}(),
				// ID: 1000622,
			},
		}
		authResp, err := uc.MidsuitService.AuthOneStep()
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}

		midsuitEmpTask, err := uc.MidsuitService.SyncEmployeeTaskMidsuit(*midsuitPayload, authResp.Token)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error syncing employee task to midsuit: ", err)
			return nil, err
		}

		midsuitID = *midsuitEmpTask
	}

	employeeTask, err := uc.Repository.CreateEmployeeTask(&entity.EmployeeTask{
		CoverPath:        req.CoverPath,
		EmployeeID:       &parsedEmployeeID,
		TemplateTaskID:   templateTaskUUID,
		SurveyTemplateID: surveyTemplateUUID,
		Name:             req.Name,
		Priority:         entity.EmployeeTaskPriorityEnum(req.Priority),
		Description:      req.Description,
		StartDate:        parsedStartDate,
		EndDate:          parsedEndDate,
		Source:           "ONBOARDING",
		MidsuitID:        &midsuitID,
	})
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error creating employee task: ", err)
		return nil, err
	}

	// create employee task attachments
	for _, attachmentReq := range req.EmployeeTaskAttachments {
		if uc.Viper.GetString("midsuit.sync") == "ACTIVE" {
			// Read the file from the given path
			fileContent, err := os.ReadFile(attachmentReq.Path)
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error reading file: ", err)
				return nil, err
			}

			// Extract the file name from the path
			fileName := filepath.Base(attachmentReq.Path)

			// Encode the file content to base64
			encodedData := base64.StdEncoding.EncodeToString(fileContent)

			// Create the payload
			midsuitAttachmentPayload := &request.SyncEmployeeTaskAttachmentMidsuitRequest{
				Name: fileName,
				Data: encodedData,
			}

			// Log the payload for debugging
			// uc.Log.Info("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] midsuit attachment payload: ", midsuitAttachmentPayload)

			// Sync to midsuit
			authResp, err := uc.MidsuitService.AuthOneStep()
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] " + err.Error())
				return nil, err
			}

			midsuitIDInt, err := strconv.Atoi(midsuitID)
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error converting midsuitID to int: ", err)
				return nil, err
			}
			_, err = uc.MidsuitService.SyncEmployeeTaskAttachmentMidsuit(midsuitIDInt, *midsuitAttachmentPayload, authResp.Token)
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error syncing employee task attachment to midsuit: ", err)
				return nil, err
			}
		}

		_, err = uc.EmployeeTaskAttachmentRepository.CreateEmployeeTaskAttachment(&entity.EmployeeTaskAttachment{
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
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error deleting employee task checklists: ", err)
		return nil, err
	}

	// create employee task checklists
	for _, checklistReq := range req.EmployeeTaskChecklists {
		// sync emp task checklist midsuit
		if uc.Viper.GetString("midsuit.sync") == "ACTIVE" {
			empResp, err := uc.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
				ID: parsedEmployeeID.String(),
			})
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error sending find employee by id message: ", err)
				return nil, err
			}
			if empResp == nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] employee not found in midsuit")
				return nil, errors.New("employee not found in midsuit")
			}

			orgResp, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
				ID: empResp.OrganizationID.String(),
			})
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error sending find organization by id message: ", err)
				return nil, err
			}
			if orgResp == nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] organization not found in midsuit")
				return nil, errors.New("organization not found in midsuit")
			}

			midsuitChecklistPayload := &request.SyncEmployeeTaskChecklistMidsuitRequest{
				AdOrgId: request.AdOrgId{
					ID: func() int {
						id, err := strconv.Atoi(orgResp.MidsuitID)
						if err != nil {
							uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error converting orgResp.MidsuitID to int: ", err)
							return 0 // or handle the error appropriately
						}
						return id
					}(),
				},
				Name:      checklistReq.Name,
				IsChecked: false,
				HCTaskID: request.HCTaskID{
					ID: func() int {
						id, err := strconv.Atoi(midsuitID)
						if err != nil {
							uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error converting employeeTask.MidsuitID to int: ", err)
							return 0 // or handle the error appropriately
						}
						return id
					}(),
				},
				HCEmployeeID: request.HcEmployeeId{
					ID: func() int {
						id, err := strconv.Atoi(empResp.MidsuitID)
						if err != nil {
							uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error converting empResp.MidsuitID to int: ", err)
							return 0 // or handle the error appropriately
						}
						return id
					}(),
				},
				ModelName: "hc_taskchecklist",
			}

			authResp, err := uc.MidsuitService.AuthOneStep()
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return nil, err
			}

			_, err = uc.MidsuitService.SyncEmployeeTaskChecklistMidsuit(*midsuitChecklistPayload, authResp.Token)
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error syncing employee task checklist to midsuit: ", err)
				return nil, err
			}
		}

		if checklistReq.MidsuitID != nil {
			exist, err := uc.EmployeeTaskChecklistRepository.FindByKeys(map[string]interface{}{
				"midsuit_id": *checklistReq.MidsuitID,
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
				var isChecked string
				if checklistReq.IsChecked != nil {
					isChecked = *checklistReq.IsChecked
				} else {
					isChecked = "NO"
				}
				_, err := uc.EmployeeTaskChecklistRepository.UpdateEmployeeTaskChecklist(&entity.EmployeeTaskChecklist{
					ID:             exist.ID,
					EmployeeTaskID: employeeTask.ID,
					Name:           checklistReq.Name,
					IsChecked:      isChecked,
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

	if uc.Viper.GetString("midsuit.sync") == "ACTIVE" {
		empResp, err := uc.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
			ID: *req.EmployeeID,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find employee by id message: ", err)
			return nil, err
		}
		if empResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] employee not found in midsuit")
			return nil, errors.New("employee not found in midsuit")
		}

		orgResp, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
			ID: empResp.OrganizationID.String(),
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find organization by id message: ", err)
			return nil, err
		}
		if orgResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] organization not found in midsuit")
			return nil, errors.New("organization not found in midsuit")
		}

		jobId := empResp.EmployeeJob["job_id"].(string)
		jobResp, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
			ID: jobId,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find job by id message: ", err)
			return nil, err
		}
		if jobResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] job not found in midsuit")
			return nil, errors.New("job not found in midsuit")
		}

		jobLevelId := empResp.EmployeeJob["job_level_id"].(string)
		jobLevelResp, err := uc.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
			ID: jobLevelId,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find job level by id message: ", err)
			return nil, err
		}
		if jobLevelResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] job level not found in midsuit")
			return nil, errors.New("job level not found in midsuit")
		}

		orgStructureId := empResp.EmployeeJob["organization_structure_id"].(string)
		orgStructureResp, err := uc.OrganizationMessage.SendFindOrganizationStructureByIDMessage(request.SendFindOrganizationStructureByIDMessageRequest{
			ID: orgStructureId,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find organization structure by id message: ", err)
			return nil, err
		}
		if orgStructureResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] organization structure not found in midsuit")
			return nil, errors.New("organization structure not found in midsuit")
		}

		midsuitPayload := &request.SyncEmployeeTaskMidsuitRequest{
			AdOrgId: request.AdOrgId{
				ID: func() int {
					id, err := strconv.Atoi(orgResp.MidsuitID)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error converting orgResp.MidsuitID to int: ", err)
						return 0 // or handle the error appropriately
					}
					return id
				}(),
				// ID: 1000024,
			},
			Name: req.Name,
			Category: request.TaskCategory{
				Identifier: "Onboarding",
				ModelName:  "ad_ref_list",
			},
			StartDate: parsedStartDate.String(),
			EndDate:   parsedEndDate.String(),
			HCEmployeeID: request.HcEmployeeId{
				ID: func() int {
					id, err := strconv.Atoi(empResp.MidsuitID)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error converting empResp.MidsuitID to int: ", err)
						return 0 // or handle the error appropriately
					}
					return id
				}(),
				// ID: 1000812,
			},
			HCJobID: request.HcJobId{
				ID: func() int {
					id, err := strconv.Atoi(jobResp.MidsuitID)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error converting jobResp.MidsuitID to int: ", err)
						return 0 // or handle the error appropriately
					}
					return id
				}(),
				// ID: 1000472,
			},
			HCJobLevelID: request.HcJobLevelId{
				ID: func() int {
					id, err := strconv.Atoi(jobLevelResp.MidsuitID)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error converting jobLevelResp.MidsuitID to int: ", err)
						return 0 // or handle the error appropriately
					}
					return id
				}(),
				// ID: 1000095,
			},
			HCOrgID: request.HcOrgId{
				ID: func() int {
					id, err := strconv.Atoi(orgStructureResp.MidsuitID)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error converting orgStructureResp.MidsuitID to int: ", err)
						return 0 // or handle the error appropriately
					}
					return id
				}(),
				// ID: 1000622,
			},
		}
		authResp, err := uc.MidsuitService.AuthOneStep()
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}

		midsuitIDInt, err := strconv.Atoi(*empTask.MidsuitID)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error converting empTask.MidsuitID to int: ", err)
			return nil, err
		}
		_, err = uc.MidsuitService.SyncUpdateEmployeeTaskMidsuit(midsuitIDInt, *midsuitPayload, authResp.Token)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error syncing employee task to midsuit: ", err)
			return nil, err
		}
	}

	employeeTask, err := uc.Repository.UpdateEmployeeTask(&entity.EmployeeTask{
		ID:               parsedID,
		CoverPath:        req.CoverPath,
		EmployeeID:       &parsedEmployeeID,
		TemplateTaskID:   templateTaskUUID,
		SurveyTemplateID: surveyTemplateUUID,
		Name:             req.Name,
		Priority:         entity.EmployeeTaskPriorityEnum(req.Priority),
		Description:      req.Description,
		StartDate:        parsedStartDate,
		EndDate:          parsedEndDate,
		VerifiedBy:       verifiedBy,
		IsDone:           req.IsDone,
		Proof:            req.ProofPath,
		Status:           entity.EmployeeTaskStatusEnum(req.Status),
		Kanban:           entity.EmployeeTaskKanbanEnum(req.Kanban),
		Notes:            req.Notes,
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

func (uc *EmployeeTaskUseCase) UpdateEmployeeTaskMidsuit(req *request.UpdateEmployeeTaskMidsuitRequest) (*response.EmployeeTaskResponse, error) {
	empTask, err := uc.Repository.FindByKeys(map[string]interface{}{
		"midsuit_id": req.ID,
	})
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

	var verifiedBy *uuid.UUID
	if req.VerifiedBy != nil && *req.VerifiedBy != "" {
		parsedVerifiedBy, err := uuid.Parse(*req.VerifiedBy)
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error parsing verified by: ", err)
			return nil, err
		}
		verifiedBy = &parsedVerifiedBy
	}

	empRespMessage, err := uc.EmployeeMessage.SendFindEmployeeByMidsuitIDMessage(*req.EmployeeMidsuitID)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error sending find employee by midsuit id message: ", err)
		return nil, err
	}
	if empRespMessage == nil {
		uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] employee not found in midsuit")
		return nil, errors.New("employee not found in midsuit")
	}

	parsedEmployeeID := empRespMessage.ID

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

	if uc.Viper.GetString("midsuit.sync") == "ACTIVE" {
		empResp, err := uc.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
			ID: parsedEmployeeID.String(),
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find employee by id message: ", err)
			return nil, err
		}
		if empResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] employee not found in midsuit")
			return nil, errors.New("employee not found in midsuit")
		}

		orgResp, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
			ID: empResp.OrganizationID.String(),
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find organization by id message: ", err)
			return nil, err
		}
		if orgResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] organization not found in midsuit")
			return nil, errors.New("organization not found in midsuit")
		}

		jobId := empResp.EmployeeJob["job_id"].(string)
		jobResp, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
			ID: jobId,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find job by id message: ", err)
			return nil, err
		}
		if jobResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] job not found in midsuit")
			return nil, errors.New("job not found in midsuit")
		}

		jobLevelId := empResp.EmployeeJob["job_level_id"].(string)
		jobLevelResp, err := uc.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
			ID: jobLevelId,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find job level by id message: ", err)
			return nil, err
		}
		if jobLevelResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] job level not found in midsuit")
			return nil, errors.New("job level not found in midsuit")
		}

		orgStructureId := empResp.EmployeeJob["organization_structure_id"].(string)
		orgStructureResp, err := uc.OrganizationMessage.SendFindOrganizationStructureByIDMessage(request.SendFindOrganizationStructureByIDMessageRequest{
			ID: orgStructureId,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find organization structure by id message: ", err)
			return nil, err
		}
		if orgStructureResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] organization structure not found in midsuit")
			return nil, errors.New("organization structure not found in midsuit")
		}

		var verifiedByMidsuitID *int
		var verifiedByJobLevelID *int
		var verifiedByJobID *int
		if req.VerifiedBy != nil && *req.VerifiedBy != "" {
			parsedVerifiedBy, err := uuid.Parse(*req.VerifiedBy)
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error parsing verified by: ", err)
				return nil, err
			}

			empRespVerifiedBy, err := uc.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
				ID: parsedVerifiedBy.String(),
			})
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find employee by id message: ", err)
				return nil, err
			}
			if empRespVerifiedBy == nil {
				uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] employee not found in midsuit")
				return nil, errors.New("employee not found in midsuit")
			}

			verifiedByMidsuitIDInt, err := strconv.Atoi(empRespVerifiedBy.MidsuitID)
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error converting empRespVerifiedBy.MidsuitID to int: ", err)
				return nil, err
			}

			verifiedByMidsuitID = &verifiedByMidsuitIDInt
			verifiedByJobLevelIDInt, err := strconv.Atoi(empRespVerifiedBy.EmployeeJob["job_level_id"].(string))
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error converting empRespVerifiedBy.EmployeeJob.job_level_id to int: ", err)
				return nil, err
			}
			verifiedByJobLevelID = &verifiedByJobLevelIDInt
			verifiedByJobIDInt, err := strconv.Atoi(empRespVerifiedBy.EmployeeJob["job_id"].(string))
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error converting empRespVerifiedBy.EmployeeJob.job_id to int: ", err)
				return nil, err
			}

			verifiedByJobID = &verifiedByJobIDInt
		} else {
			verifiedByMidsuitID = nil
			verifiedByJobLevelID = nil
			verifiedByJobID = nil
		}
		if req.VerifiedBy != nil && *req.VerifiedBy != "" {
			midsuitPayload := &request.SyncEmployeeTaskMidsuitRequest{
				AdOrgId: request.AdOrgId{
					ID: func() int {
						id, err := strconv.Atoi(orgResp.MidsuitID)
						if err != nil {
							uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error converting orgResp.MidsuitID to int: ", err)
							return 0 // or handle the error appropriately
						}
						return id
					}(),
					// ID: 1000024,
				},
				Name: req.Name,
				Category: request.TaskCategory{
					Identifier: "Onboarding",
					ModelName:  "ad_ref_list",
				},
				StartDate: parsedStartDate.String(),
				EndDate:   parsedEndDate.String(),
				HCEmployeeID: request.HcEmployeeId{
					ID: func() int {
						id, err := strconv.Atoi(empResp.MidsuitID)
						if err != nil {
							uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error converting empResp.MidsuitID to int: ", err)
							return 0 // or handle the error appropriately
						}
						return id
					}(),
					// ID: 1000812,
				},
				HCApproverID: request.HcApproverId{
					ID: *verifiedByMidsuitID,
					// ID: 1000812,
				},
				HCJobID: request.HcJobId{
					ID: func() int {
						id, err := strconv.Atoi(jobResp.MidsuitID)
						if err != nil {
							uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error converting jobResp.MidsuitID to int: ", err)
							return 0 // or handle the error appropriately
						}
						return id
					}(),
					// ID: 1000472,
				},
				HCJob2ID: request.HcJobId{
					ID: *verifiedByJobID,
					// ID: 1000472,
				},
				HCJobLevel2ID: request.HcJobLevelId{
					ID: *verifiedByJobLevelID,
					// ID: 1000095,
				},
				HCJobLevelID: request.HcJobLevelId{
					ID: func() int {
						id, err := strconv.Atoi(jobLevelResp.MidsuitID)
						if err != nil {
							uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error converting jobLevelResp.MidsuitID to int: ", err)
							return 0 // or handle the error appropriately
						}
						return id
					}(),
					// ID: 1000095,
				},
				HCOrgID: request.HcOrgId{
					ID: func() int {
						id, err := strconv.Atoi(orgStructureResp.MidsuitID)
						if err != nil {
							uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error converting orgStructureResp.MidsuitID to int: ", err)
							return 0 // or handle the error appropriately
						}
						return id
					}(),
					// ID: 1000622,
				},
			}
			authResp, err := uc.MidsuitService.AuthOneStep()
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return nil, err
			}

			midsuitIDInt, err := strconv.Atoi(*empTask.MidsuitID)
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error converting empTask.MidsuitID to int: ", err)
				return nil, err
			}
			_, err = uc.MidsuitService.SyncUpdateEmployeeTaskMidsuit(midsuitIDInt, *midsuitPayload, authResp.Token)
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error syncing employee task to midsuit: ", err)
				return nil, err
			}
		}
	}

	employeeTask, err := uc.Repository.UpdateEmployeeTask(&entity.EmployeeTask{
		ID:               empTask.ID,
		CoverPath:        req.CoverPath,
		EmployeeID:       &parsedEmployeeID,
		TemplateTaskID:   templateTaskUUID,
		SurveyTemplateID: surveyTemplateUUID,
		Name:             req.Name,
		Priority:         entity.EmployeeTaskPriorityEnum(req.Priority),
		Description:      req.Description,
		StartDate:        parsedStartDate,
		EndDate:          parsedEndDate,
		VerifiedBy:       verifiedBy,
		IsDone:           req.IsDone,
		Proof:            req.ProofPath,
		Status:           entity.EmployeeTaskStatusEnum(req.Status),
		Kanban:           entity.EmployeeTaskKanbanEnum(req.Kanban),
		Notes:            req.Notes,
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
		if checklistReq.MidsuitID != nil {
			exist, err := uc.EmployeeTaskChecklistRepository.FindByKeys(map[string]interface{}{
				"midsuit_id": *checklistReq.MidsuitID,
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
				checklistIds = append(checklistIds, exist.ID)
				var verifiedBy *uuid.UUID
				if checklistReq.VerifiedBy != nil {
					// parsedVerifiedBy, err := uuid.Parse(*checklistReq.VerifiedBy)
					// if err != nil {
					// 	uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error parsing verified by: ", err)
					// 	return nil, err
					// }
					verifiedByResp, err := uc.EmployeeMessage.SendFindEmployeeByMidsuitIDMessage(*checklistReq.VerifiedBy)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find employee by midsuit id message: ", err)
						return nil, err
					}
					verifiedBy = &verifiedByResp.ID
				}
				var isChecked string
				if checklistReq.IsChecked != nil {
					isChecked = *checklistReq.IsChecked
				} else {
					isChecked = "NO"
				}
				_, err := uc.EmployeeTaskChecklistRepository.UpdateEmployeeTaskChecklist(&entity.EmployeeTaskChecklist{
					ID:             exist.ID,
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
				// post to midsuit
				var midsuitID string
				if uc.Viper.GetString("midsuit.sync") == "ACTIVE" {
					empResp, err := uc.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
						ID: parsedEmployeeID.String(),
					})
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find employee by id message: ", err)
						return err
					}
					if empResp == nil {
						uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] employee not found in midsuit")
						return errors.New("employee not found in midsuit")
					}

					orgResp, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
						ID: empResp.OrganizationID.String(),
					})
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find organization by id message: ", err)
						return err
					}
					if orgResp == nil {
						uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] organization not found in midsuit")
						return errors.New("organization not found in midsuit")
					}

					jobId := empResp.EmployeeJob["job_id"].(string)
					jobResp, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
						ID: jobId,
					})
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find job by id message: ", err)
						return err
					}
					if jobResp == nil {
						uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] job not found in midsuit")
						return errors.New("job not found in midsuit")
					}

					jobLevelId := empResp.EmployeeJob["job_level_id"].(string)
					jobLevelResp, err := uc.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
						ID: jobLevelId,
					})
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find job level by id message: ", err)
						return err
					}
					if jobLevelResp == nil {
						uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] job level not found in midsuit")
						return errors.New("job level not found in midsuit")
					}

					orgStructureId := empResp.EmployeeJob["organization_structure_id"].(string)
					orgStructureResp, err := uc.OrganizationMessage.SendFindOrganizationStructureByIDMessage(request.SendFindOrganizationStructureByIDMessageRequest{
						ID: orgStructureId,
					})
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find organization structure by id message: ", err)
						return err
					}
					if orgStructureResp == nil {
						uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] organization structure not found in midsuit")
						return errors.New("organization structure not found in midsuit")
					}

					midsuitPayload := &request.SyncEmployeeTaskMidsuitRequest{
						AdOrgId: request.AdOrgId{
							// ID: orgResp.MidsuitID,
							ID: 1000024,
						},
						Name: templateTask.Name,
						Category: request.TaskCategory{
							Identifier: "Onboarding",
							ModelName:  "ad_ref_list",
						},
						StartDate: parsedJoinedDate.String(),
						EndDate:   parsedJoinedDate.AddDate(0, 0, *templateTask.DueDuration).String(),
						HCEmployeeID: request.HcEmployeeId{
							ID: func() int {
								id, err := strconv.Atoi(empResp.MidsuitID)
								if err != nil {
									uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error converting empResp.MidsuitID to int: ", err)
									return 0 // or handle the error appropriately
								}
								return id
							}(),
							// ID: 1000108,
						},
						HCJobID: request.HcJobId{
							ID: func() int {
								id, err := strconv.Atoi(jobResp.MidsuitID)
								if err != nil {
									uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error converting jobResp.MidsuitID to int: ", err)
									return 0 // or handle the error appropriately
								}
								return id
							}(),
							// ID: 1000472,
						},
						HCJobLevelID: request.HcJobLevelId{
							ID: func() int {
								id, err := strconv.Atoi(jobLevelResp.MidsuitID)
								if err != nil {
									uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error converting jobLevelResp.MidsuitID to int: ", err)
									return 0 // or handle the error appropriately
								}
								return id
							}(),
							// ID: 1000095,
						},
						HCOrgID: request.HcOrgId{
							ID: func() int {
								id, err := strconv.Atoi(orgStructureResp.MidsuitID)
								if err != nil {
									uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error converting orgStructureResp.MidsuitID to int: ", err)
									return 0 // or handle the error appropriately
								}
								return id
							}(),
							// ID: 1000622,
						},
					}
					authResp, err := uc.MidsuitService.AuthOneStep()
					if err != nil {
						uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
						return err
					}

					midsuitEmpTask, err := uc.MidsuitService.SyncEmployeeTaskMidsuit(*midsuitPayload, authResp.Token)
					if err != nil {
						uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error syncing employee task to midsuit: ", err)
						return err
					}

					midsuitID = *midsuitEmpTask
				}
				createdEmpTask, err := uc.Repository.CreateEmployeeTask(&entity.EmployeeTask{
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
					MidsuitID:      &midsuitID,
				})
				if err != nil {
					uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTasksForRecruitment] error creating employee task: ", err)
					// continue
					return err
				}
				if len(templateTask.TemplateTaskChecklists) > 0 {
					for _, checklist := range templateTask.TemplateTaskChecklists {
						_, err := uc.EmployeeTaskChecklistRepository.CreateEmployeeTaskChecklist(&entity.EmployeeTaskChecklist{
							EmployeeTaskID: createdEmpTask.ID,
							Name:           checklist.Name,
							IsChecked:      "NO",
						})
						if err != nil {
							uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTasksForRecruitment] error creating employee task checklist: ", err)
							// continue
							return err
						}
					}
				}
				if len(templateTask.TemplateTaskAttachments) > 0 {
					for _, attachment := range templateTask.TemplateTaskAttachments {
						_, err := uc.EmployeeTaskAttachmentRepository.CreateEmployeeTaskAttachment(&entity.EmployeeTaskAttachment{
							EmployeeTaskID: createdEmpTask.ID,
							Path:           attachment.Path,
						})
						if err != nil {
							uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTasksForRecruitment] error creating employee task attachment: ", err)
							// continue
							return err
						}
					}
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

func (uc *EmployeeTaskUseCase) FindAllPaginatedByEmployeeID(employeeID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}) (*[]response.EmployeeTaskResponse, int64, error) {
	employeeTasks, total, err := uc.Repository.FindAllPaginatedByEmployeeID(employeeID, page, pageSize, search, sort)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.FindAllPaginatedByEmployeeID] error finding all employee tasks by employee id: ", err)
		return nil, 0, err
	}

	var responses []response.EmployeeTaskResponse
	for _, employeeTask := range *employeeTasks {
		responses = append(responses, *uc.DTO.ConvertEntityToResponse(&employeeTask))
	}

	return &responses, total, nil
}

func (uc *EmployeeTaskUseCase) FindByIDForResponse(id uuid.UUID) (*response.EmployeeTaskResponse, error) {
	employeeTask, err := uc.Repository.FindByIDForResponse(id)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.FindByIDForResponse] error finding employee task by id: ", err)
		return nil, err
	}

	if employeeTask == nil {
		return nil, errors.New("employee task not found")
	}

	return uc.DTO.ConvertEntityToResponse(employeeTask), nil
}

func (uc *EmployeeTaskUseCase) FindAllPaginatedSurvey(page, pageSize int, search string, sort map[string]interface{}) (*[]response.EmployeeTaskResponse, int64, error) {
	employeeTasks, total, err := uc.Repository.FindAllPaginatedSurvey(page, pageSize, search, sort)
	if err != nil {
		uc.Log.Error("[EmployeeTaskUseCase.FindAllPaginatedSurvey] error finding all employee tasks: ", err)
		return nil, 0, err
	}

	var responses []response.EmployeeTaskResponse
	for _, employeeTask := range *employeeTasks {
		responses = append(responses, *uc.DTO.ConvertEntityToResponse(&employeeTask))
	}

	return &responses, total, nil
}
