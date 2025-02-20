package handler

import (
	"net/http"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-onboarding-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IEmployeeTaskAttachmentHandler interface {
	FindByID(ctx *gin.Context)
	DeleteEmployeeTaskAttachment(ctx *gin.Context)
}

type EmployeeTaskAttachmentHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IEmployeeTaskAttachmentUseCase
}

func NewEmployeeTaskAttachmentHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IEmployeeTaskAttachmentUseCase,
) IEmployeeTaskAttachmentHandler {
	return &EmployeeTaskAttachmentHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func EmployeeTaskAttachmentHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IEmployeeTaskAttachmentHandler {
	useCase := usecase.EmployeeTaskAttachmentUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	return NewEmployeeTaskAttachmentHandler(log, viper, validate, useCase)
}

// FindByID Find employee task attachment by ID
//
// @Summary Find employee task attachment by ID
// @Description Find employee task attachment by ID
// @Tags Employee Task Attachments
// @Accept json
// @Produce json
// @Param id path string true "Employee Task Attachment ID"
// @Success 200 {object} response.EmployeeTaskAttachmentResponse
// @Router /api/employee-task-attachments/{id} [get]
func (h *EmployeeTaskAttachmentHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[EmployeeTaskAttachmentHandler.FindByID] " + err.Error())
		utils.BadRequestResponse(ctx, "invalid id", err.Error())
		return
	}

	attachment, err := h.UseCase.FindByID(parsedId)
	if err != nil {
		h.Log.Error("[EmployeeTaskAttachmentHandler.FindByID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	if attachment == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "template task attachment not found", "template task attachment not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success get template task attachment", attachment)
}

// DeleteEmployeeTaskAttachment Delete employee task attachment by ID
//
// @Summary Delete employee task attachment by ID
// @Description Delete employee task attachment by ID
// @Tags Employee Task Attachments
// @Accept json
// @Produce json
// @Param id path string true "Employee Task Attachment ID"
// @Success 200 {string} string "success delete employee task attachment"
// @Router /api/employee-task-attachments/{id} [delete]
func (h *EmployeeTaskAttachmentHandler) DeleteEmployeeTaskAttachment(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[EmployeeTaskAttachmentHandler.DeleteEmployeeTaskAttachment] " + err.Error())
		utils.BadRequestResponse(ctx, "invalid id", err.Error())
		return
	}

	err = h.UseCase.DeleteEmployeeTaskAttachment(parsedId)
	if err != nil {
		h.Log.Error("[EmployeeTaskAttachmentHandler.DeleteEmployeeTaskAttachment] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success delete employee task attachment", nil)
}
