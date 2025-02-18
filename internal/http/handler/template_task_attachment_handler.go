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

type ITemplateTaskAttachmentHandler interface {
	FindByID(ctx *gin.Context)
	DeleteTemplateTaskAttachment(ctx *gin.Context)
}

type TemplateTaskAttachmentHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.ITemplateTaskAttachmentUseCase
}

func NewTemplateTaskAttachmentHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.ITemplateTaskAttachmentUseCase,
) ITemplateTaskAttachmentHandler {
	return &TemplateTaskAttachmentHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func TemplateTaskAttachmentHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ITemplateTaskAttachmentHandler {
	useCase := usecase.TemplateTaskAttachmentUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	return NewTemplateTaskAttachmentHandler(log, viper, validate, useCase)
}

// FindByID Find template task attachment by ID
//
// @Summary Find template task attachment by ID
// @Description Find template task attachment by ID
// @Tags Template Task Attachments
// @Accept json
// @Produce json
// @Param id path string true "Template Task Attachment ID"
// @Success 200 {object} response.TemplateTaskAttachmentResponse
// @Router /api/template-tasks/{id} [get]
func (h *TemplateTaskAttachmentHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[CoverHandler.FindByID] " + err.Error())
		utils.BadRequestResponse(ctx, "invalid id", err.Error())
		return
	}

	attachment, err := h.UseCase.FindByID(parsedId)
	if err != nil {
		h.Log.Error("[TemplateTaskAttachmentHandler.FindByID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	if attachment == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "template task attachment not found", "template task attachment not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success get template task attachment", attachment)
}

// DeleteTemplateTaskAttachment delete template task attachment by ID
//
// @Summary Delete template task attachment by ID
// @Description Delete template task attachment by ID
// @Tags Template Task Attachments
// @Accept json
// @Produce json
// @Param id path string true "Template Task Attachment ID"
// @Success 204 {string} string "success delete template task attachment"
// @Router /api/template-tasks/{id} [delete]
func (h *TemplateTaskAttachmentHandler) DeleteTemplateTaskAttachment(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[TemplateTaskAttachmentHandler.DeleteTemplateTaskAttachment] " + err.Error())
		utils.BadRequestResponse(ctx, "invalid id", err.Error())
		return
	}

	err = h.UseCase.DeleteTemplateTaskAttachment(parsedId)
	if err != nil {
		h.Log.Error("[TemplateTaskAttachmentHandler.DeleteTemplateTaskAttachment] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusNoContent, "success delete template task attachment", nil)
}
