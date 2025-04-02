package request

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
	"github.com/go-playground/validator/v10"
)

func MaritalStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.MaritalStatusEnum(status) {
	case entity.MARITAL_STATUS_ENUM_SINGLE,
		entity.MARITAL_STATUS_ENUM_MARRIED,
		entity.MARITAL_STATUS_ENUM_DIVORCED,
		entity.MARITAL_STATUS_ENUM_WIDOWED,
		entity.MARITAL_STATUS_ENUM_ANY:
		return true
	default:
		return false
	}
}

func UserStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.UserStatus(status) {
	case entity.USER_ACTIVE,
		entity.USER_INACTIVE,
		entity.USER_PENDING:
		return true
	default:
		return false
	}
}

func UserGenderValidation(fl validator.FieldLevel) bool {
	gender := fl.Field().String()
	if gender == "" {
		return true
	}
	switch entity.UserGender(gender) {
	case entity.MALE,
		entity.FEMALE:
		return true
	default:
		return false
	}
}

func EducationLevelValidation(fl validator.FieldLevel) bool {
	level := fl.Field().String()
	if level == "" {
		return true
	}
	switch entity.EducationLevelEnum(level) {
	case entity.EDUCATION_LEVEL_ENUM_DOCTORAL,
		entity.EDUCATION_LEVEL_ENUM_MASTER,
		entity.EDUCATION_LEVEL_ENUM_BACHELOR,
		entity.EDUCATION_LEVEL_ENUM_D1,
		entity.EDUCATION_LEVEL_ENUM_D2,
		entity.EDUCATION_LEVEL_ENUM_D3,
		entity.EDUCATION_LEVEL_ENUM_D4,
		entity.EDUCATION_LEVEL_ENUM_SD,
		entity.EDUCATION_LEVEL_ENUM_SMA,
		entity.EDUCATION_LEVEL_ENUM_SMP,
		entity.EDUCATION_LEVEL_ENUM_TK:
		return true
	default:
		return false
	}
}

func TemplateTaskPriorityValidation(fl validator.FieldLevel) bool {
	priority := fl.Field().String()
	if priority == "" {
		return true
	}
	switch entity.TemplateTaskPriorityEnum(priority) {
	case entity.TEMPLATE_TASK_PRIORITY_ENUM_LOW,
		entity.TEMPLATE_TASK_PRIORITY_ENUM_MEDIUM,
		entity.TEMPLATE_TASK_PRIORITY_ENUM_HIGH:
		return true
	default:
		return false
	}
}

func TemplateTaskStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.TemplateTaskStatusEnum(status) {
	case entity.TEMPLATE_TASK_STATUS_ENUM_ACTIVE,
		entity.TEMPLATE_TASK_STATUS_ENUM_INACTIVE:
		return true
	default:
		return false
	}
}

func EmployeeTaskPriorityValidation(fl validator.FieldLevel) bool {
	priority := fl.Field().String()
	if priority == "" {
		return true
	}
	switch entity.EmployeeTaskPriorityEnum(priority) {
	case entity.EMPLOYEE_TASK_PRIORITY_ENUM_LOW,
		entity.EMPLOYEE_TASK_PRIORITY_ENUM_MEDIUM,
		entity.EMPLOYEE_TASK_PRIORITY_ENUM_HIGH:
		return true
	default:
		return false
	}
}

func EmployeeTaskStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.EmployeeTaskStatusEnum(status) {
	case entity.EMPLOYEE_TASK_STATUS_ENUM_ACTIVE,
		entity.EMPLOYEE_TASK_STATUS_ENUM_INACTIVE:
		return true
	default:
		return false
	}
}

func EmployeeTaskKanbanValidation(fl validator.FieldLevel) bool {
	kanban := fl.Field().String()
	if kanban == "" {
		return true
	}
	switch entity.EmployeeTaskKanbanEnum(kanban) {
	case entity.EMPLOYEE_TASK_KANBAN_ENUM_TODO,
		entity.EPMLOYEE_TASK_KANBAN_ENUM_IN_PROGRESS,
		entity.EMPLOYEE_TASK_KANBAN_ENUM_NEED_REVIEW,
		entity.EMPLOYEE_TASK_KANBAN_ENUM_COMPLETED:
		return true
	default:
		return false
	}
}

func EventStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.EventStatusEnum(status) {
	case entity.EVENT_STATUS_ENUM_UPCOMING,
		entity.EVENT_STATUS_ENUM_ONGOING,
		entity.EVENT_STATUS_ENUM_FINISHED:
		return true
	default:
		return false
	}
}

func SurveyTemplateStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.SurveyTemplateStatusEnum(status) {
	case entity.SURVEY_TEMPLATE_STATUS_ENUM_DRAFT,
		entity.SURVEY_TEMPLATE_STATUS_ENUM_SUBMITTED:
		return true
	default:
		return false
	}
}
