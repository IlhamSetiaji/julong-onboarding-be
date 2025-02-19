package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-onboarding-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type IEmployeeTaskHandler interface {
	CreateEmployeeTask(ctx *gin.Context)
	UpdateEmployeeTask(ctx *gin.Context)
	DeleteEmployeeTask(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	CountByKanbanAndEmployeeID(ctx *gin.Context)
	FindAllByEmployeeID(ctx *gin.Context)
}

type EmployeeTaskHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IEmployeeTaskUseCase
	DB       *gorm.DB
}

func NewEmployeeTaskHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IEmployeeTaskUseCase,
	db *gorm.DB,
) IEmployeeTaskHandler {
	return &EmployeeTaskHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
		DB:       db,
	}
}

func EmployeeTaskHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IEmployeeTaskHandler {
	validate := config.NewValidator(viper)
	db := config.NewDatabase()
	useCase := usecase.EmployeeTaskUseCaseFactory(log, viper)
	return NewEmployeeTaskHandler(log, viper, validate, useCase, db)
}

// CreateEmployeeTask create new employee task
//
// @Summary Create new employee task
// @Description Create new employee task
// @Tags Employee Task
// @Accept  multipart/form-data
// @Produce  json
// @Param body body request.CreateEmployeeTaskRequest true "Create Template Task"
// @Success 201 {object} response.EmployeeTaskResponse
// @Security BearerAuth
// @Router /employee-tasks [post]
func (h *EmployeeTaskHandler) CreateEmployeeTask(ctx *gin.Context) {
	var req request.CreateEmployeeTaskRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.Log.Error("[EmployeeTaskHandler.CreateEmployeeTask] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[EmployeeTaskHandler.CreateEmployeeTask] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	// Handle attachments file upload manually
	form, _ := ctx.MultipartForm()
	files := form.File["employee_task_attachments[file]"]
	checklistNames := form.Value["employee_task_checklists[name]"]
	checklistIds := form.Value["employee_task_checklists[id]"]

	if len(files) > 0 {
		for _, file := range files {
			timestamp := time.Now().UnixNano()
			filePath := "storage/employee_tasks/attachments/" + strconv.FormatInt(timestamp, 10) + "_" + file.Filename
			if err := ctx.SaveUploadedFile(file, filePath); err != nil {
				h.Log.Error("failed to save attachment file: ", err)
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save attachment file", err.Error())
				return
			}

			req.EmployeeTaskAttachments = append(req.EmployeeTaskAttachments, request.EmployeeTaskAttachmentRequest{
				File: nil,
				Path: filePath,
			})
		}
	}

	if len(checklistNames) > 0 {
		var checklistId *string
		for i, name := range checklistNames {
			if i < len(checklistIds) {
				checklistId = &checklistIds[i]
			} else {
				checklistId = nil
			}

			req.EmployeeTaskChecklists = append(req.EmployeeTaskChecklists, request.EmployeeTaskChecklistRequest{
				ID:   checklistId,
				Name: name,
			})
		}
	}

	tx := h.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		h.Log.Warnf("Failed begin transaction : %+v", tx.Error)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to begin transaction", tx.Error.Error())
		return
	}
	defer tx.Rollback()

	res, err := h.UseCase.CreateEmployeeTask(&req)
	if err != nil {
		h.Log.Error("[EmployeeTaskHandler.CreateEmployeeTask] " + err.Error())
		tx.Rollback()
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	if err := tx.Commit().Error; err != nil {
		h.Log.Warnf("Failed commit transaction : %+v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to commit transaction", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "success create employee task", res)
}

// UpdateEmployeeTask update employee task
//
// @Summary Update employee task
// @Description Update employee task
// @Tags Employee Task
// @Accept  multipart/form-data
// @Produce  json
// @Param id path string true "Employee Task ID"
// @Param body body request.UpdateEmployeeTaskRequest true "Update Employee Task"
// @Success 200 {object} response.EmployeeTaskResponse
// @Security BearerAuth
// @Router /employee-tasks/{id} [put]
func (h *EmployeeTaskHandler) UpdateEmployeeTask(ctx *gin.Context) {
	var req request.UpdateEmployeeTaskRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.Log.Error("[EmployeeTaskHandler.CreateEmployeeTask] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[EmployeeTaskHandler.CreateEmployeeTask] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	// Handle attachments file upload manually
	form, _ := ctx.MultipartForm()
	files := form.File["employee_task_attachments[file]"]
	checklistNames := form.Value["employee_task_checklists[name]"]
	checklistIds := form.Value["employee_task_checklists[id]"]
	checklistIsCheckeds := form.Value["employee_task_checklists[is_checked]"]
	checklistVerifiedBys := form.Value["employee_task_checklists[verified_by]"]

	if len(files) > 0 {
		for _, file := range files {
			timestamp := time.Now().UnixNano()
			filePath := "storage/employee_tasks/attachments/" + strconv.FormatInt(timestamp, 10) + "_" + file.Filename
			if err := ctx.SaveUploadedFile(file, filePath); err != nil {
				h.Log.Error("failed to save attachment file: ", err)
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save attachment file", err.Error())
				return
			}

			req.EmployeeTaskAttachments = append(req.EmployeeTaskAttachments, request.EmployeeTaskAttachmentRequest{
				File: nil,
				Path: filePath,
			})
		}
	}

	if len(checklistNames) > 0 {
		var checklistId *string
		var checklistIsChecked *string
		var checklistVerifiedBy *string
		for i, name := range checklistNames {
			if i < len(checklistIds) {
				checklistId = &checklistIds[i]
			} else {
				checklistId = nil
			}

			if i < len(checklistIsCheckeds) {
				checklistIsChecked = &checklistIsCheckeds[i]
			} else {
				checklistIsChecked = nil
			}

			if i < len(checklistVerifiedBys) {
				checklistVerifiedBy = &checklistVerifiedBys[i]
			} else {
				checklistVerifiedBy = nil
			}

			req.EmployeeTaskChecklists = append(req.EmployeeTaskChecklists, request.EmployeeTaskChecklistRequest{
				ID:         checklistId,
				Name:       name,
				IsChecked:  checklistIsChecked,
				VerifiedBy: checklistVerifiedBy,
			})
		}
	}

	if req.Proof != nil {
		timestamp := time.Now().UnixNano()
		filePath := "storage/covers/" + strconv.FormatInt(timestamp, 10) + "_" + req.Proof.Filename
		if err := ctx.SaveUploadedFile(req.Proof, filePath); err != nil {
			h.Log.Error("failed to save cover file: ", err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save cover file", err.Error())
			return
		}

		req.Proof = nil
		req.ProofPath = &filePath
	}

	tx := h.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		h.Log.Warnf("Failed begin transaction : %+v", tx.Error)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to begin transaction", tx.Error.Error())
		return
	}
	defer tx.Rollback()

	res, err := h.UseCase.UpdateEmployeeTask(&req)
	if err != nil {
		h.Log.Error("[EmployeeTaskHandler.CreateEmployeeTask] " + err.Error())
		tx.Rollback()
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	if err := tx.Commit().Error; err != nil {
		h.Log.Warnf("Failed commit transaction : %+v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to commit transaction", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "success create employee task", res)
}

// DeleteEmployeeTask delete employee task
//
// @Summary Delete employee task
// @Description Delete employee task
// @Tags Employee Task
// @Accept  json
// @Produce  json
// @Param id path string true "Employee Task ID"
// @Success 200 {string} string "success delete employee task"
// @Security BearerAuth
// @Router /employee-tasks/{id} [delete]
func (h *EmployeeTaskHandler) DeleteEmployeeTask(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.BadRequestResponse(ctx, "id is required", "id is required")
		return
	}

	employeeTaskID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequestResponse(ctx, "invalid id", "invalid id")
		return
	}

	if err := h.UseCase.DeleteEmployeeTask(employeeTaskID); err != nil {
		h.Log.Error("[EmployeeTaskHandler.DeleteEmployeeTask] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success delete employee task", nil)
}

// FindByID find employee task by id
//
// @Summary Find employee task by id
// @Description Find employee task by id
// @Tags Employee Task
// @Accept  json
// @Produce  json
// @Param id path string true "Employee Task ID"
// @Success 200 {object} response.EmployeeTaskResponse
// @Security BearerAuth
// @Router /employee-tasks/{id} [get]
func (h *EmployeeTaskHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.BadRequestResponse(ctx, "id is required", "id is required")
		return
	}

	employeeTaskID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequestResponse(ctx, "invalid id", "invalid id")
		return
	}

	res, err := h.UseCase.FindByID(employeeTaskID)
	if err != nil {
		h.Log.Error("[EmployeeTaskHandler.FindByID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success find employee task", res)
}

// FindAllPaginated find all employee task paginated
//
// @Summary Find all employee task paginated
// @Description Find all employee task paginated
// @Tags Employee Task
// @Accept  json
// @Produce  json
// @Param page query int false "Page"
// @Param page_size query int false "Page Size"
// @Param search query string false "Search"
// @Param created_at query string false "Created At"
// @Success 200 {object} response.EmployeeTaskResponse
// @Security BearerAuth
// @Router /employee-tasks [get]
func (h *EmployeeTaskHandler) FindAllPaginated(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	search := ctx.Query("search")
	if search == "" {
		search = ""
	}

	createdAt := ctx.Query("created_at")
	if createdAt == "" {
		createdAt = "DESC"
	}

	sort := map[string]interface{}{
		"created_at": createdAt,
	}

	res, total, err := h.UseCase.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		h.Log.Error("[EmployeeTaskHandler.FindAllPaginated] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success find all employee task", gin.H{
		"data":  res,
		"total": total,
	})
}

// CountByKanbanAndEmployeeID count employee task by kanban and employee id
//
// @Summary Count employee task by kanban and employee id
// @Description Count employee task by kanban and employee id
// @Tags Employee Task
// @Accept  json
// @Produce  json
// @Param kanban query string true "Kanban"
// @Param employee_id query string true "Employee ID"
// @Success 200 {object} response.EmployeeTaskResponse
// @Security BearerAuth
// @Router /employee-tasks/count [get]
func (h *EmployeeTaskHandler) CountByKanbanAndEmployeeID(ctx *gin.Context) {
	kanban := ctx.Query("kanban")
	if kanban == "" {
		utils.BadRequestResponse(ctx, "kanban is required", "kanban is required")
		return
	}

	employeeID := ctx.Query("employee_id")
	if employeeID == "" {
		utils.BadRequestResponse(ctx, "employee_id is required", "employee_id is required")
		return
	}

	parsedEmployeeID, err := uuid.Parse(employeeID)
	if err != nil {
		utils.BadRequestResponse(ctx, "invalid employee_id", "invalid employee_id")
		return
	}

	res, err := h.UseCase.CountByKanbanAndEmployeeID(entity.EmployeeTaskKanbanEnum(kanban), parsedEmployeeID)
	if err != nil {
		h.Log.Error("[EmployeeTaskHandler.CountByKanbanAndEmployeeID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success count employee task", res)
}

// FindAllByEmployeeID find all employee task by employee id
//
// @Summary Find all employee task by employee id
// @Description Find all employee task by employee id
// @Tags Employee Task
// @Accept  json
// @Produce  json
// @Param employee_id query string true "Employee ID"
// @Success 200 {object} response.EmployeeTaskResponse
// @Security BearerAuth
// @Router /employee-tasks/employee [get]
func (h *EmployeeTaskHandler) FindAllByEmployeeID(ctx *gin.Context) {
	employeeID := ctx.Query("employee_id")
	if employeeID == "" {
		utils.BadRequestResponse(ctx, "employee_id is required", "employee_id is required")
		return
	}

	parsedEmployeeID, err := uuid.Parse(employeeID)
	if err != nil {
		utils.BadRequestResponse(ctx, "invalid employee_id", "invalid employee_id")
		return
	}

	res, err := h.UseCase.FindAllByEmployeeID(parsedEmployeeID)
	if err != nil {
		h.Log.Error("[EmployeeTaskHandler.FindAllByEmployeeID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success find all employee task", res)
}
