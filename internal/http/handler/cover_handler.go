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

type ICoverHandler interface {
	CreateCover(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	UpdateCover(ctx *gin.Context)
	DeleteCover(ctx *gin.Context)
}

type CoverHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.ICoverUseCase
	DB       *gorm.DB
}

func NewCoverHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.ICoverUseCase,
	db *gorm.DB,
) ICoverHandler {
	return &CoverHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
		DB:       db,
	}
}

func CoverHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ICoverHandler {
	useCase := usecase.CoverUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	db := config.NewDatabase()
	return NewCoverHandler(log, viper, validate, useCase, db)
}

// CreateCover create a new cover
//
// @Summary Create a new cover
// @Description Create a new cover
// @Tags Covers
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File"
// @Success 201 {object} response.CoverResponse
// @Security BearerAuth
// @Router /covers [post]
func (h *CoverHandler) CreateCover(ctx *gin.Context) {
	var req request.CreateCoverRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.Log.Error("[CoverHandler.CreateCover] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[CoverHandler.CreateCover] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	// handle file upload
	if req.File != nil {
		timestamp := time.Now().UnixNano()
		filePath := "storage/covers/" + strconv.FormatInt(timestamp, 10) + "_" + req.File.Filename
		if err := ctx.SaveUploadedFile(req.File, filePath); err != nil {
			h.Log.Error("failed to save cover file: ", err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save cover file", err.Error())
			return
		}

		req.File = nil
		req.Path = filePath
	}

	tx := h.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		h.Log.Warnf("Failed begin transaction : %+v", tx.Error)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to begin transaction", tx.Error.Error())
		return
	}
	defer tx.Rollback()

	res, err := h.UseCase.CreateCover(&req)
	if err != nil {
		h.Log.Error("[CoverHandler.CreateCover] " + err.Error())
		tx.Rollback()
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	if err := tx.Commit().Error; err != nil {
		h.Log.Warnf("Failed commit transaction : %+v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to commit transaction", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "success create cover", res)
}

// FindByID find cover by id
//
// @Summary Find cover by id
// @Description Find cover by id
// @Tags Covers
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} response.CoverResponse
// @Security BearerAuth
// @Router /covers/{id} [get]
func (h *CoverHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[CoverHandler.FindByID] " + err.Error())
		utils.BadRequestResponse(ctx, "invalid id", err.Error())
		return
	}

	res, err := h.UseCase.FindByID(parsedId)
	if err != nil {
		h.Log.Error("[CoverHandler.FindByID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	if res == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "cover not found", "cover not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success get cover", res)
}

// FindAllPaginated find all covers with pagination
//
// @Summary Find all covers with pagination
// @Description Find all covers with pagination
// @Tags Covers
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param page_size query int false "Page Size"
// @Param search query string false "Search"
// @Param created_at query string false "Created At"
// @Success 200 {object} response.CoverResponse
// @Security BearerAuth
// @Router /covers [get]
func (h *CoverHandler) FindAllPaginated(ctx *gin.Context) {
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
		h.Log.Error("[CoverHandler.FindAllPaginated] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success get all covers", gin.H{
		"covers": res,
		"total":  total,
	})
}

// UpdateCover update cover by id
//
// @Summary Update cover by id
// @Description Update cover by id
// @Tags Covers
// @Accept multipart/form-data
// @Produce json
// @Param id formData string true "ID"
// @Param file formData file true "File"
// @Success 200 {object} response.CoverResponse
// @Security BearerAuth
// @Router /covers/update [put]
func (h *CoverHandler) UpdateCover(ctx *gin.Context) {
	id := ctx.PostForm("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[CoverHandler.UpdateCover] " + err.Error())
		utils.BadRequestResponse(ctx, "invalid id", err.Error())
		return
	}

	var req request.UpdateCoverRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.Log.Error("[CoverHandler.UpdateCover] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[CoverHandler.UpdateCover] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	// handle file upload
	if req.File != nil {
		timestamp := time.Now().UnixNano()
		filePath := "storage/covers/" + strconv.FormatInt(timestamp, 10) + "_" + req.File.Filename
		if err := ctx.SaveUploadedFile(req.File, filePath); err != nil {
			h.Log.Error("failed to save cover file: ", err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save cover file", err.Error())
			return
		}

		req.File = nil
		req.Path = filePath
	}
	tx := h.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		h.Log.Warnf("Failed begin transaction : %+v", tx.Error)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to begin transaction", tx.Error.Error())
		return
	}
	defer tx.Rollback()

	res, err := h.UseCase.UpdateCover(parsedId, &req)
	if err != nil {
		h.Log.Error("[CoverHandler.UpdateCover] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	if err := tx.Commit().Error; err != nil {
		h.Log.Warnf("Failed commit transaction : %+v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to commit transaction", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success update cover", res)
}

// DeleteCover delete cover by id
//
// @Summary Delete cover by id
// @Description Delete cover by id
// @Tags Covers
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 204
// @Security BearerAuth
// @Router /covers/{id} [delete]
func (h *CoverHandler) DeleteCover(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[CoverHandler.DeleteCover] " + err.Error())
		utils.BadRequestResponse(ctx, "invalid id", err.Error())
		return
	}

	tx := h.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		h.Log.Warnf("Failed begin transaction : %+v", tx.Error)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to begin transaction", tx.Error.Error())
		return
	}
	defer tx.Rollback()

	err = h.UseCase.DeleteCover(parsedId)
	if err != nil {
		h.Log.Error("[CoverHandler.DeleteCover] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	if err := tx.Commit().Error; err != nil {
		h.Log.Warnf("Failed commit transaction : %+v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to commit transaction", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusNoContent, "success delete cover", nil)
}
