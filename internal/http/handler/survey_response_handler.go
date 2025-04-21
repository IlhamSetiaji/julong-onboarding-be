package handler

import (
	"fmt"
	"net/http"
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
	"github.com/xuri/excelize/v2"
)

// contains checks if a string exists in a slice of strings.
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

type ISurveyResponseHandler interface {
	CreateOrUpdateSurveyResponses(ctx *gin.Context)
	CreateOrUpdateSurveyResponsesBulk(ctx *gin.Context)
	ExportSurveyResponses(ctx *gin.Context)
}

type SurveyResponseHandler struct {
	Log                 *logrus.Logger
	Viper               *viper.Viper
	Validate            *validator.Validate
	UseCase             usecase.ISurveyResponseUseCase
	UserHelper          helper.IUserHelper
	EmployeeTaskUseCase usecase.IEmployeeTaskUseCase
}

func NewSurveyResponseHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.ISurveyResponseUseCase,
	userHelper helper.IUserHelper,
	employeeTaskUseCase usecase.IEmployeeTaskUseCase,
) ISurveyResponseHandler {
	return &SurveyResponseHandler{
		Log:                 log,
		Viper:               viper,
		Validate:            validate,
		UseCase:             useCase,
		UserHelper:          userHelper,
		EmployeeTaskUseCase: employeeTaskUseCase,
	}
}

func SurveyResponseHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ISurveyResponseHandler {
	useCase := usecase.SurveyResponseUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	userHelper := helper.UserHelperFactory(log)
	employeeTaskUseCase := usecase.EmployeeTaskUseCaseFactory(log, viper)
	return NewSurveyResponseHandler(log, viper, validate, useCase, userHelper, employeeTaskUseCase)
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

func (h *SurveyResponseHandler) CreateOrUpdateSurveyResponsesBulk(ctx *gin.Context) {
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
		h.Log.Error("Failed to parse form-data: ", err)
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	surveyTemplateID := ctx.Request.FormValue("survey_template_id")
	employeeTaskID := ctx.Request.FormValue("employee_task_id")
	kanban := ctx.Request.FormValue("kanban")
	answerIDs := ctx.PostFormArray("answers[id]")
	questionIDs := ctx.PostFormArray("answers[question_id]")
	answers := ctx.PostFormArray("answers[answer]")

	// Process each answer
	var payload request.SurveyResponseBulkRequest
	payload.SurveyTemplateID = surveyTemplateID
	payload.EmployeeTaskID = employeeTaskID
	payload.Kanban = kanban

	for i := range questionIDs {
		questionID := questionIDs[i]

		// Get the answer text if available
		var answer string
		if len(answers) > i {
			answer = answers[i]
		} else {
			answer = ""
		}

		// Get the answer ID if available
		var answerID *string
		if len(answerIDs) > i {
			answerID = &answerIDs[i]
		} else {
			answerID = nil
		}

		// Get the answer file if available
		var answerFilePath string
		fileKey := fmt.Sprintf("answers[%d][answer_file]", i) // Dynamically construct the file key
		fileHeaders := ctx.Request.MultipartForm.File[fileKey]
		if len(fileHeaders) > 0 && fileHeaders[0] != nil {
			file := fileHeaders[0]
			timestamp := time.Now().UnixNano()
			filePath := "storage/answers/files/" + strconv.FormatInt(timestamp, 10) + "_" + file.Filename
			if err := ctx.SaveUploadedFile(file, filePath); err != nil {
				h.Log.Error("Failed to save answer file: ", err)
				utils.ErrorResponse(ctx, 500, "error", "Failed to save answer file")
				return
			}
			answerFilePath = filePath

			h.Log.Infof("answer file path: %v", answerFilePath)
		}

		// Append the answer to the payload
		payload.Answers = append(payload.Answers, request.AnswerBulkRequest{
			ID:         answerID,
			QuestionID: questionID,
			Answer:     answer,
			AnswerPath: answerFilePath,
		})
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("Error when validating payload: %v", err)
		utils.ErrorResponse(ctx, 400, "error", err.Error())
		return
	}

	h.Log.Infof("payload: %v", payload)

	questionResponse, err := h.UseCase.CreateOrUpdateSurveyResponsesBulk(&payload)
	if err != nil {
		h.Log.Errorf("Error when creating or updating question responses: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, "success answer question", questionResponse)
}

func (h *SurveyResponseHandler) ExportSurveyResponses(ctx *gin.Context) {
	employeeTasks, err := h.EmployeeTaskUseCase.FindAllSurvey()
	if err != nil {
		h.Log.Errorf("Error when getting employee tasks: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	if len(*employeeTasks) == 0 {
		h.Log.Error("No employee tasks found")
		utils.ErrorResponse(ctx, 404, "error", "No employee tasks found")
		return
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export survey responses", err.Error())
			return
		}
	}()

	f.SetSheetName("Sheet1", "Survey Responses")

	headers := []string{"Employee Task Name", "Employee Name", "Survey Name"}

	var maxQuestions int
	for _, employeeTask := range *employeeTasks {
		if len(employeeTask.SurveyTemplate.Questions) > maxQuestions {
			maxQuestions = len(employeeTask.SurveyTemplate.Questions)
		}

		for _, question := range employeeTask.SurveyTemplate.Questions {
			if !contains(headers, question.Question) {
				headers = append(headers, question.Question)
			}
		}
	}
	// for i := 1; i <= maxQuestions; i++ {
	// 	headers = append(headers, fmt.Sprintf("Question %d", i))
	// }

	// Define header style
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "#000000",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#90EE90"}, // Light green background
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		h.Log.Errorf("Error when creating header style: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	// Write headers to the Excel file
	for i, header := range headers {
		col := string(rune('A' + i))
		cell := fmt.Sprintf("%s1", col)
		f.SetCellValue("Survey Responses", cell, header)
		f.SetCellStyle("Survey Responses", cell, cell, headerStyle)
	}

	// Write data rows
	for rowIndex, employeeTask := range *employeeTasks {
		row := rowIndex + 2 // Start from the second row
		f.SetCellValue("Survey Responses", fmt.Sprintf("A%d", row), employeeTask.EmployeeName)
		f.SetCellValue("Survey Responses", fmt.Sprintf("B%d", row), employeeTask.SurveyTemplate.Title)

		// Fetch survey responses
		surveyResponses, err := h.EmployeeTaskUseCase.FindByIDForResponse(employeeTask.ID.String())
		if err != nil {
			h.Log.Errorf("Error when getting survey responses: %v", err)
			utils.ErrorResponse(ctx, 500, "error", err.Error())
			return
		}

		// Write answers to the corresponding question columns
		for questionIndex, question := range surveyResponses.SurveyTemplate.Questions {
			col := string(rune('C' + questionIndex))
			var concatenatedValue string

			if len(question.SurveyResponses) > 0 {
				for _, answer := range question.SurveyResponses {
					if concatenatedValue != "" {
						concatenatedValue += ", "
					}
					if answer.AnswerFile == "" {
						concatenatedValue += answer.Answer
					} else {
						concatenatedValue += h.Viper.GetString("app.url") + answer.AnswerFile
					}
				}
			}

			f.SetCellValue("Survey Responses", fmt.Sprintf("%s%d", col, row), concatenatedValue)
		}
	}

	// Set column widths for better readability
	f.SetColWidth("Survey Responses", "A", "C", 20)
	f.SetColWidth("Survey Responses", "D", string(rune('C'+maxQuestions-1)), 30)

	// Set response headers for file download
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=survey_responses.xlsx")
	ctx.Header("Content-Transfer-Encoding", "binary")

	// Write the Excel file to the response
	if err := f.Write(ctx.Writer); err != nil {
		fmt.Println(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export survey responses", err.Error())
		return
	}
}
