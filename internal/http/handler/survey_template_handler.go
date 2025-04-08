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
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ISurveyTemplateHandler interface {
	CreateSurveyTemplate(ctx *gin.Context)
	UpdateSurveyTemplate(ctx *gin.Context)
	FindAllSurveyTemplatesPaginated(ctx *gin.Context)
	FindSurveyTemplateByID(ctx *gin.Context)
	DeleteSurveyTemplate(ctx *gin.Context)
}

type SurveyTemplateHandler struct {
	Log               *logrus.Logger
	Viper             *viper.Viper
	Validate          *validator.Validate
	UseCase           usecase.ISurveyTemplateUseCase
	QuestionUseCase   usecase.IQuestionUseCase
	AnswerTypeUseCase usecase.IAnswerTypeUseCase
	DB                *gorm.DB
}

func NewSurveyTemplateHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.ISurveyTemplateUseCase,
	questionUseCase usecase.IQuestionUseCase,
	answerTypeUseCase usecase.IAnswerTypeUseCase,
	db *gorm.DB) ISurveyTemplateHandler {
	return &SurveyTemplateHandler{
		Log:               log,
		Viper:             viper,
		Validate:          validate,
		UseCase:           useCase,
		QuestionUseCase:   questionUseCase,
		AnswerTypeUseCase: answerTypeUseCase,
		DB:                db,
	}
}

func SurveyTemplateHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ISurveyTemplateHandler {
	useCase := usecase.SurveyTemplateUseCaseFactory(log, viper)
	questionUseCase := usecase.QuestionUseCaseFactory(log, viper)
	answerTypeUseCase := usecase.AnswerTypeUseCaseFactory(log)
	validate := config.NewValidator(viper)
	db := config.NewDatabase()
	return NewSurveyTemplateHandler(
		log,
		viper,
		validate,
		useCase,
		questionUseCase,
		answerTypeUseCase,
		db,
	)
}

func (h *SurveyTemplateHandler) CreateSurveyTemplate(ctx *gin.Context) {
	var req request.CreateOrUpdateQuestions
	if err := ctx.ShouldBind(&req); err != nil {
		h.Log.Error("[SurveyTemplateHandler.CreateSurveyTemplate] Error when binding request: ", err)
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[SurveyTemplateHandler.CreateSurveyTemplate] Error when validating request: ", err)
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	// Handle attachments file upload manually
	form, _ := ctx.MultipartForm()
	files := form.File["questions[attachment]"]
	answerTypes := form.Value["questions[answer_type_id]"]
	questions := form.Value["questions[question]"]
	optionText := form.Value["questions[question_options][option_text]"]
	// questionOptions := make([][]string, len(optionText))

	h.Log.Info("Answer Types: ", answerTypes)

	for i, answerType := range answerTypes {
		ans, err := h.AnswerTypeUseCase.FindByID(answerType)
		if err != nil {
			h.Log.Error("[SurveyTemplateHandler.CreateSurveyTemplate] Error when finding answer type by id: ", err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to find answer type by id", err.Error())
			return
		}

		if ans == nil {
			h.Log.Error("[SurveyTemplateHandler.CreateSurveyTemplate] Error when finding answer type by id: answer type not found")
			utils.ErrorResponse(ctx, http.StatusNotFound, "failed to find answer type by id", "answer type not found")
			return
		}

		h.Log.Info("Answer Type: ", ans.Name)

		if ans.Name == "Attachment" {
			if len(files) > i { // Ensure the file index exists
				file := files[i]
				timestamp := time.Now().UnixNano()
				filePath := "storage/questions/attachments/" + strconv.FormatInt(timestamp, 10) + "_" + file.Filename
				if err := ctx.SaveUploadedFile(file, filePath); err != nil {
					h.Log.Error("failed to save attachment file: ", err)
					utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save attachment file", err.Error())
					return
				}

				req.Questions = append(req.Questions, request.QuestionRequest{
					Attachment:     nil,
					AttachmentPath: filePath,
					AnswerTypeID:   answerType,
					Question:       questions[i],
				})
			}
		} else {
			if len(questions) > i { // Ensure the question index exists
				var questionOptions []request.QuestionOptionRequest
				if len(optionText) > 0 {
					for _, option := range optionText {
						questionOptions = append(questionOptions, request.QuestionOptionRequest{
							OptionText: option,
						})
					}
				}
				req.Questions = append(req.Questions, request.QuestionRequest{
					AnswerTypeID:    answerType,
					Question:        questions[i],
					QuestionOptions: questionOptions,
				})
			}
		}
	}

	h.Log.Info("Questions: ", req.Questions)
	h.Log.Info("QuestionOptions: ", optionText)

	tx := h.DB.WithContext(ctx.Request.Context()).Begin()
	if tx.Error != nil {
		h.Log.Error("[SurveyTemplateHandler.CreateSurveyTemplate] Error when starting transaction: ", tx.Error)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to start transaction", tx.Error.Error())
		return
	}
	defer tx.Rollback()

	res, err := h.QuestionUseCase.CreateOrUpdateQuestions(&req)
	if err != nil {
		h.Log.Error("[SurveyTemplateHandler.CreateSurveyTemplate] Error when creating or updating questions: ", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to create or update questions", err.Error())
		return
	}

	if err := tx.Commit().Error; err != nil {
		h.Log.Warnf("[SurveyTemplateHandler.CreateSurveyTemplate] Error when committing transaction: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to commit transaction", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res)
}

func (h *SurveyTemplateHandler) UpdateSurveyTemplate(ctx *gin.Context) {
	var req request.CreateOrUpdateQuestions
	if err := ctx.ShouldBind(&req); err != nil {
		h.Log.Error("[SurveyTemplateHandler.UpdateSurveyTemplate] Error when binding request: ", err)
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[SurveyTemplateHandler.UpdateSurveyTemplate] Error when validating request: ", err)
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	// Handle attachments file upload manually
	form, _ := ctx.MultipartForm()
	files := form.File["questions[attachment]"]
	answerTypes := form.Value["questions[answer_type_id]"]
	questions := form.Value["questions[question]"]
	optionText := form.Value["questions[question_options][option_text]"]
	// questionOptions := make([][]string, len(optionText))

	for i, answerType := range answerTypes {
		ans, err := h.AnswerTypeUseCase.FindByID(answerType)
		if err != nil {
			h.Log.Error("[SurveyTemplateHandler.CreateSurveyTemplate] Error when finding answer type by id: ", err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to find answer type by id", err.Error())
			return
		}

		if ans == nil {
			h.Log.Error("[SurveyTemplateHandler.CreateSurveyTemplate] Error when finding answer type by id: answer type not found")
			utils.ErrorResponse(ctx, http.StatusNotFound, "failed to find answer type by id", "answer type not found")
			return
		}

		h.Log.Info("Answer Type: ", ans.Name)

		if ans.Name == "Attachment" {
			if len(files) > i { // Ensure the file index exists
				file := files[i]
				timestamp := time.Now().UnixNano()
				filePath := "storage/questions/attachments/" + strconv.FormatInt(timestamp, 10) + "_" + file.Filename
				if err := ctx.SaveUploadedFile(file, filePath); err != nil {
					h.Log.Error("failed to save attachment file: ", err)
					utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save attachment file", err.Error())
					return
				}

				req.Questions = append(req.Questions, request.QuestionRequest{
					Attachment:     nil,
					AttachmentPath: filePath,
					AnswerTypeID:   answerType,
					Question:       questions[i],
				})
			}
		} else {
			if len(questions) > i { // Ensure the question index exists
				var questionOptions []request.QuestionOptionRequest
				if len(optionText) > 0 {
					for _, option := range optionText {
						questionOptions = append(questionOptions, request.QuestionOptionRequest{
							OptionText: option,
						})
					}
				}
				req.Questions = append(req.Questions, request.QuestionRequest{
					AnswerTypeID:    answerType,
					Question:        questions[i],
					QuestionOptions: questionOptions,
				})
			}
		}
	}

	tx := h.DB.WithContext(ctx.Request.Context()).Begin()
	if tx.Error != nil {
		h.Log.Error("[SurveyTemplateHandler.UpdateSurveyTemplate] Error when starting transaction: ", tx.Error)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to start transaction", tx.Error.Error())
		return
	}
	defer tx.Rollback()

	res, err := h.QuestionUseCase.CreateOrUpdateQuestions(&req)
	if err != nil {
		h.Log.Error("[SurveyTemplateHandler.UpdateSurveyTemplate] Error when creating or updating questions: ", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to create or update questions", err.Error())
		return
	}

	if err := tx.Commit().Error; err != nil {
		h.Log.Warnf("[SurveyTemplateHandler.UpdateSurveyTemplate] Error when committing transaction: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to commit transaction", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res)
}

func (h *SurveyTemplateHandler) FindAllSurveyTemplatesPaginated(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		h.Log.Error("[SurveyTemplateHandler.FindAllSurveyTemplate] Error when parsing page query param: ", err)
		utils.BadRequestResponse(ctx, "invalid page query param", err.Error())
		return
	}

	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil {
		h.Log.Error("[SurveyTemplateHandler.FindAllSurveyTemplate] Error when parsing page_size query param: ", err)
		utils.BadRequestResponse(ctx, "invalid page_size query param", err.Error())
		return
	}

	search := ctx.Query("search")
	sort := map[string]interface{}{
		"created_at": "desc",
	}

	res, total, err := h.UseCase.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		h.Log.Error("[SurveyTemplateHandler.FindAllSurveyTemplate] Error when finding all survey template: ", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to find all survey template", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", gin.H{
		"total":            total,
		"survey_templates": res,
	})
}

func (h *SurveyTemplateHandler) FindSurveyTemplateByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.Log.Error("[SurveyTemplateHandler.FindSurveyTemplateByID] Error when getting id from url param: id is empty")
		utils.BadRequestResponse(ctx, "id is required", "id is required")
		return
	}

	res, err := h.UseCase.FindSurveyTemplateByID(id)
	if err != nil {
		h.Log.Error("[SurveyTemplateHandler.FindSurveyTemplateByID] Error when finding survey template by id: ", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to find survey template by id", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res)
}

func (h *SurveyTemplateHandler) DeleteSurveyTemplate(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.Log.Error("[SurveyTemplateHandler.DeleteSurveyTemplate] Error when getting id from url param: id is empty")
		utils.BadRequestResponse(ctx, "id is required", "id is required")
		return
	}

	err := h.UseCase.DeleteSurveyTemplate(id)
	if err != nil {
		h.Log.Error("[SurveyTemplateHandler.DeleteSurveyTemplate] Error when deleting survey template: ", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to delete survey template", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", nil)
}
