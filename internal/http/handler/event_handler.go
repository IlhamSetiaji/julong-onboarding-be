package handler

import (
	"net/http"
	"strconv"

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

type IEventHandler interface {
	CreateEvent(ctx *gin.Context)
	UpdateEvent(ctx *gin.Context)
	DeleteEvent(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	FindByID(ctx *gin.Context)
}

type EventHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IEventUseCase
	DB       *gorm.DB
}

func NewEventHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IEventUseCase,
	db *gorm.DB,
) IEventHandler {
	return &EventHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
		DB:       db,
	}
}

func EventHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IEventHandler {
	db := config.NewDatabase()
	validate := config.NewValidator(viper)
	useCase := usecase.EventUseCaseFactory(log, viper)
	return NewEventHandler(log, viper, validate, useCase, db)
}

// CreateEvent creates a new event
//
// @Summary Create a new event
// @Description Create a new event
// @Tags Events
// @Accept json
// @Produce json
// @Param event body request.CreateEventRequest true "Event data"
// @Security BearerAuth
// @Success 201 {object} response.EventResponse
// @Router /events [post]
func (h *EventHandler) CreateEvent(ctx *gin.Context) {
	var req request.CreateEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Error("[EventHandler.CreateEvent] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[EventHandler.CreateEvent] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	event, err := h.UseCase.CreateEvent(ctx, &req)
	if err != nil {
		h.Log.Error("[EventHandler.CreateEvent] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create event", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Event created", event)
}

// UpdateEvent updates an event
//
// @Summary Update an event
// @Description Update an event
// @Tags Events
// @Accept json
// @Produce json
// @Param event body request.UpdateEventRequest true "Event data"
// @Security BearerAuth
// @Success 200 {object} response.EventResponse
// @Router /events [put]
func (h *EventHandler) UpdateEvent(ctx *gin.Context) {
	var req request.UpdateEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Error("[EventHandler.UpdateEvent] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[EventHandler.UpdateEvent] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	event, err := h.UseCase.UpdateEvent(ctx, &req)
	if err != nil {
		h.Log.Error("[EventHandler.UpdateEvent] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update event", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Event updated", event)
}

// DeleteEvent deletes an event
//
// @Summary Delete an event
// @Description Delete an event
// @Tags Events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Security BearerAuth
// @Success 204 "Event deleted"
// @Router /events/{id} [delete]
func (h *EventHandler) DeleteEvent(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.Log.Error("[EventHandler.DeleteEvent] ID is required")
		utils.BadRequestResponse(ctx, "ID is required", "ID is required")
		return
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[EventHandler.DeleteEvent] " + err.Error())
		utils.BadRequestResponse(ctx, "Invalid ID", "Invalid ID")
		return
	}

	err = h.UseCase.DeleteEvent(ctx, parsedID)
	if err != nil {
		h.Log.Error("[EventHandler.DeleteEvent] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete event", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusNoContent, "Event deleted", nil)
}

// FindAllPaginated finds all events with pagination
//
// @Summary Find all events with pagination
// @Description Find all events with pagination
// @Tags Events
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param page_size query int false "Page Size"
// @Param search query string false "Search"
// @Param created_at query string false "Created At"
// @Success 200 {object} response.EventResponse
// @Router /events [get]
func (h *EventHandler) FindAllPaginated(ctx *gin.Context) {
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
	events, total, err := h.UseCase.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		h.Log.Error("[EventHandler.FindAllPaginated] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find events", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Events found", gin.H{
		"events": events,
		"total":  total,
	})
}

// FindByID finds an event by ID
//
// @Summary Find an event by ID
// @Description Find an event by ID
// @Tags Events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} response.EventResponse
// @Router /events/{id} [get]
func (h *EventHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.Log.Error("[EventHandler.FindByID] ID is required")
		utils.BadRequestResponse(ctx, "ID is required", "ID is required")
		return
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[EventHandler.FindByID] " + err.Error())
		utils.BadRequestResponse(ctx, "Invalid ID", "Invalid ID")
		return
	}

	event, err := h.UseCase.FindByID(parsedID)
	if err != nil {
		h.Log.Error("[EventHandler.FindByID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find event", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Event found", event)
}
