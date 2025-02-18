package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
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

type ITemplateTaskHandler interface {
	CreateTemplateTask(ctx *gin.Context)
	UpdateTemplateTask(ctx *gin.Context)
	DeleteTemplateTask(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
}

type TemplateTaskHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.ITemplateTaskUseCase
	DB       *gorm.DB
}

func NewTemplateTaskHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.ITemplateTaskUseCase,
	db *gorm.DB,
) ITemplateTaskHandler {
	return &TemplateTaskHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
		DB:       db,
	}
}

func TemplateTaskHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ITemplateTaskHandler {
	useCase := usecase.TemplateTaskUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	db := config.NewDatabase()
	return NewTemplateTaskHandler(log, viper, validate, useCase, db)
}

// CreateTemplateTask create new template task
//
// @Summary Create new template task
// @Description Create new template task
// @Tags Template Tasks
// @Accept multipart/form-data
// @Produce json
// @Param body body request.CreateTemplateTaskRequest true "Create Template Task"
// @Success 201 {object} response.TemplateTaskResponse
// @Security BearerAuth
// @Router /template-tasks [post]
func (h *TemplateTaskHandler) CreateTemplateTask(ctx *gin.Context) {
	var req request.CreateTemplateTaskRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.Log.Error("[TemplateTaskHandler.CreateTemplateTask] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[TemplateTaskHandler.CreateTemplateTask] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	// Handle cover file upload
	if req.CoverFile != nil {
		timestamp := time.Now().UnixNano()
		filePath := "storage/template_tasks/covers/" + strconv.FormatInt(timestamp, 10) + "_" + req.CoverFile.Filename
		if err := ctx.SaveUploadedFile(req.CoverFile, filePath); err != nil {
			h.Log.Error("failed to save cover file: ", err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save cover file", err.Error())
			return
		}

		req.CoverFile = nil
		req.CoverPath = filePath
	}

	// Handle attachments file upload manually
	form, _ := ctx.MultipartForm()
	files := form.File["template_task_attachments[file]"]          // Matches form-data key from Postman
	checklistNames := form.Value["template_task_checklists[name]"] // Matches form-data key from Postman
	checklistIds := form.Value["template_task_checklists[id]"]     // Matches form-data key from Postman

	if len(files) > 0 {
		for _, file := range files {
			timestamp := time.Now().UnixNano()
			filePath := "storage/template_tasks/attachments/" + strconv.FormatInt(timestamp, 10) + "_" + file.Filename
			if err := ctx.SaveUploadedFile(file, filePath); err != nil {
				h.Log.Error("failed to save attachment file: ", err)
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save attachment file", err.Error())
				return
			}

			req.TemplateTaskAttachments = append(req.TemplateTaskAttachments, request.TemplateTaskAttachmentRequest{
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

			req.TemplateTaskChecklists = append(req.TemplateTaskChecklists, request.TemplateTaskChecklistRequest{
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

	res, err := h.UseCase.CreateTemplateTask(&req)
	if err != nil {
		h.Log.Error("[TemplateTaskHandler.CreateTemplateTask] " + err.Error())
		tx.Rollback()
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	if err := tx.Commit().Error; err != nil {
		h.Log.Warnf("Failed commit transaction : %+v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to commit transaction", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "success create template task", res)
}

// UpdateTemplateTask update template task by id
//
// @Summary Update template task by id
// @Description Update template task by id
// @Tags Template Tasks
// @Accept multipart/form-data
// @Produce json
// @Param body body request.UpdateTemplateTaskRequest true "Update Template Task"
// @Success 200 {object} response.TemplateTaskResponse
// @Security BearerAuth
// @Router /template-tasks/update [put]
func (h *TemplateTaskHandler) UpdateTemplateTask(ctx *gin.Context) {
	var req request.UpdateTemplateTaskRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.Log.Error("[TemplateTaskHandler.CreateTemplateTask] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[TemplateTaskHandler.CreateTemplateTask] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	// Handle cover file upload
	if req.CoverFile != nil {
		timestamp := time.Now().UnixNano()
		filePath := "storage/template_tasks/covers/" + strconv.FormatInt(timestamp, 10) + "_" + req.CoverFile.Filename
		if err := ctx.SaveUploadedFile(req.CoverFile, filePath); err != nil {
			h.Log.Error("failed to save cover file: ", err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save cover file", err.Error())
			return
		}

		req.CoverFile = nil
		req.CoverPath = filePath
	}

	// Handle attachments file upload manually
	form, _ := ctx.MultipartForm()
	files := form.File["template_task_attachments[file]"]          // Matches form-data key from Postman
	checklistNames := form.Value["template_task_checklists[name]"] // Matches form-data key from Postman
	checklistIds := form.Value["template_task_checklists[id]"]     // Matches form-data key from Postman

	if len(files) > 0 {
		for _, file := range files {
			timestamp := time.Now().UnixNano()
			filePath := "storage/template_tasks/attachments/" + strconv.FormatInt(timestamp, 10) + "_" + file.Filename
			if err := ctx.SaveUploadedFile(file, filePath); err != nil {
				h.Log.Error("failed to save attachment file: ", err)
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save attachment file", err.Error())
				return
			}

			req.TemplateTaskAttachments = append(req.TemplateTaskAttachments, request.TemplateTaskAttachmentRequest{
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

			req.TemplateTaskChecklists = append(req.TemplateTaskChecklists, request.TemplateTaskChecklistRequest{
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

	res, err := h.UseCase.UpdateTemplateTask(&req)
	if err != nil {
		h.Log.Error("[TemplateTaskHandler.CreateTemplateTask] " + err.Error())
		tx.Rollback()
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	if err := tx.Commit().Error; err != nil {
		h.Log.Warnf("Failed commit transaction : %+v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to commit transaction", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "success create template task", res)
}

// DeleteTemplateTask delete template task by id
//
// @Summary Delete template task by id
// @Description Delete template task by id
// @Tags Template Tasks
// @Param id path string true "Template Task ID"
// @Success 204
// @Security BearerAuth
// @Router /template-tasks/{id} [delete]
func (h *TemplateTaskHandler) DeleteTemplateTask(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[TemplateTaskHandler.DeleteTemplateTask] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	tx := h.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		h.Log.Warnf("Failed begin transaction : %+v", tx.Error)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to begin transaction", tx.Error.Error())
		return
	}
	defer tx.Rollback()

	err = h.UseCase.DeleteTemplateTask(parsedId)
	if err != nil {
		h.Log.Error("[TemplateTaskHandler.DeleteTemplateTask] " + err.Error())
		tx.Rollback()
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	if err := tx.Commit().Error; err != nil {
		h.Log.Warnf("Failed commit transaction : %+v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to commit transaction", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusNoContent, "success delete template task", nil)
}

// FindByID find template task by id
//
// @Summary Find template task by id
// @Description Find template task by id
// @Tags Template Tasks
// @Param id path string true "Template Task ID"
// @Success 200 {object} response.TemplateTaskResponse
// @Security BearerAuth
// @Router /template-tasks/{id} [get]
func (h *TemplateTaskHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[TemplateTaskHandler.FindByID] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	res, err := h.UseCase.FindByID(parsedId)
	if err != nil {
		h.Log.Error("[TemplateTaskHandler.FindByID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success get template task", res)
}

// FindAllPaginated find all template tasks with pagination
//
// @Summary Find all template tasks with pagination
// @Description Find all template tasks with pagination
// @Tags Template Tasks
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param page_size query int false "Page Size"
// @Param search query string false "Search"
// @Param created_at query string false "Created At"
// @Success 200 {object} response.TemplateTaskResponse
// @Security BearerAuth
// @Router /template-tasks [get]
func (h *TemplateTaskHandler) FindAllPaginated(ctx *gin.Context) {
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
		h.Log.Error("[TemplateTaskHandler.FindAllPaginated] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success get template tasks", gin.H{
		"data": res,
		"meta": gin.H{
			"total": total,
		},
	})
}
