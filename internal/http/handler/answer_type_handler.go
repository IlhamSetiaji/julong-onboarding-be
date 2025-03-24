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

type IAnswerTypeHandler interface {
	FindAll(ctx *gin.Context)
}

type AnswerTypeHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IAnswerTypeUseCase
}

func NewAnswerTypeHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IAnswerTypeUseCase,
) IAnswerTypeHandler {
	return &AnswerTypeHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func AnswerTypeHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IAnswerTypeHandler {
	useCase := usecase.AnswerTypeUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewAnswerTypeHandler(log, viper, validate, useCase)
}

func (h *AnswerTypeHandler) FindAll(ctx *gin.Context) {
	answerTypes, err := h.UseCase.FindAll()
	if err != nil {
		h.Log.Errorf("[AnswerTypeHandler.FindAll] error when getting answer types: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error when getting answer types", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success find all answer types", answerTypes)
}
