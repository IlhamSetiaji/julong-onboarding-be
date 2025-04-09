package handler

import (
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/helper"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-onboarding-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ISurveyResponseHandler interface {
	CreateOrUpdateSurveyResponses(ctx *gin.Context)
}

type SurveyResponseHandler struct {
	Log        *logrus.Logger
	Viper      *viper.Viper
	Validate   *validator.Validate
	UseCase    usecase.ISurveyResponseUseCase
	UserHelper helper.IUserHelper
}

func NewSurveyResponseHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.ISurveyResponseUseCase,
	userHelper helper.IUserHelper,
) ISurveyResponseHandler {
	return &SurveyResponseHandler{
		Log:        log,
		Viper:      viper,
		Validate:   validate,
		UseCase:    useCase,
		UserHelper: userHelper,
	}
}

func SurveyResponseHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ISurveyResponseHandler {
	useCase := usecase.SurveyResponseUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	userHelper := helper.UserHelperFactory(log)
	return NewSurveyResponseHandler(log, viper, validate, useCase, userHelper)
}

func (h *SurveyResponseHandler) CreateOrUpdateSurveyResponses(ctx *gin.Context) {
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
		h.Log.Error("Failed to parse form-data: ", err)
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	questionID := ctx.Request.FormValue("question_id")
	answerIDs := ctx.PostFormArray("answers[id]")
	jobPostingIDs := ctx.PostFormArray("answers[survey_template_id]")
	userProfileIDs := ctx.PostFormArray("answers[employee_task_id]")
	answers := ctx.PostFormArray("answers[answer]")
	answerFiles := ctx.Request.MultipartForm.File["answers[][answer_file]"]
	// Process each answer
	var payload request.SurveyResponseRequest
	payload.QuestionID = questionID
	for i := range userProfileIDs {
		jobPostingID := jobPostingIDs[i]
		userProfileID := userProfileIDs[i]
		var answer string
		if len(answers) > i {
			answer = answers[i]
		} else {
			answer = ""
		}

		var answerID *string
		if len(answerIDs) > i {
			answerID = &answerIDs[i]
		} else {
			answerID = nil
		}

		h.Log.Infof("answer: %v", answer)

		var answerFilePath string

		if len(answerFiles) > i {
			file := answerFiles[i]
			timestamp := time.Now().UnixNano()
			filePath := "storage/answers/files/" + strconv.FormatInt(timestamp, 10) + "_" + file.Filename
			if err := ctx.SaveUploadedFile(file, filePath); err != nil {
				h.Log.Error("Failed to save answer file: ", err)
				utils.ErrorResponse(ctx, 500, "error", "Failed to save answer file")
				return
			}
			answerFilePath = filePath
		}

		payload.Answers = append(payload.Answers, request.AnswerRequest{
			ID:               answerID,
			SurveyTemplateID: jobPostingID,
			EmployeeTaskID:   userProfileID,
			Answer:           answer,
			AnswerPath:       answerFilePath,
		})
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("Error when validating payload: %v", err)
		utils.ErrorResponse(ctx, 400, "error", err.Error())
		return
	}

	h.Log.Infof("payload: %v", payload)

	questionResponse, err := h.UseCase.CreateOrUpdateSurveyResponses(&payload)
	if err != nil {
		h.Log.Errorf("Error when creating or updating question responses: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	// embed url to answer file
	// for i, qr := range questionResponse.SurveyResponses {
	// 	if qr.AnswerFile != "" {
	// 		(questionResponse.SurveyResponses)[i].AnswerFile = h.Viper.GetString("app.url") + qr.AnswerFile
	// 	}
	// }

	utils.SuccessResponse(ctx, 201, "success answer question", questionResponse)
}
