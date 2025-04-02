package config

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func NewValidator(viper *viper.Viper) *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("marital_status_validation", request.MaritalStatusValidation)
	validate.RegisterValidation("user_status_validation", request.UserStatusValidation)
	validate.RegisterValidation("user_gender_validation", request.UserGenderValidation)
	validate.RegisterValidation("education_level_validation", request.EducationLevelValidation)
	validate.RegisterValidation("template_task_status_validation", request.TemplateTaskStatusValidation)
	validate.RegisterValidation("template_task_priority_validation", request.TemplateTaskPriorityValidation)
	validate.RegisterValidation("employee_task_status_validation", request.EmployeeTaskStatusValidation)
	validate.RegisterValidation("employee_task_priority_validation", request.EmployeeTaskPriorityValidation)
	validate.RegisterValidation("employee_task_kanban_validation", request.EmployeeTaskKanbanValidation)
	validate.RegisterValidation("event_status_validation", request.EventStatusValidation)
	validate.RegisterValidation("survey_template_status_validation", request.SurveyTemplateStatusValidation)
	return validate
}
