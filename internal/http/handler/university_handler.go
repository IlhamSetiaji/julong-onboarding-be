package handler

import (
	"net/http"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-onboarding-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IUniversityHandler interface {
	FindAll(ctx *gin.Context)
}

type UniversityHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IUniversityUseCase
}

func NewUniversityHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IUniversityUseCase,
) IUniversityHandler {
	return &UniversityHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func UniversityHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IUniversityHandler {
	useCase := usecase.UniversityUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewUniversityHandler(log, viper, validate, useCase)
}

// FindAll find all universities
//
//		@Summary		Find all universities
//		@Description	Find all universities
//		@Tags			Universities
//		@Accept			json
//	 @Produce		json
//		@Success		200	{object} response.UniversityResponse
//		@Security BearerAuth
//		@Router			/universities [get]
func (h *UniversityHandler) FindAll(ctx *gin.Context) {
	responses, err := h.UseCase.FindAll()
	if err != nil {
		h.Log.Errorf("[UniversityHandler.FindAll] error when finding all university: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success get all universities", responses)
}
