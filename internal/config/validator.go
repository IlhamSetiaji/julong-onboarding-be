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
	return validate
}
