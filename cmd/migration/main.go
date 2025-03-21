package main

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
)

func main() {
	viper := config.NewViper()
	log := config.NewLogrus(viper)
	db := config.NewDatabase()

	// migrate the schema
	err := db.AutoMigrate(
		&entity.Cover{},
		&entity.TemplateTask{},
		&entity.TemplateTaskAttachment{},
		&entity.TemplateTaskChecklist{},
		&entity.EmployeeTask{},
		&entity.EmployeeTaskAttachment{},
		&entity.EmployeeTaskFiles{},
		&entity.EmployeeHiring{},
		&entity.EmployeeTaskChecklist{},
		&entity.Event{},
		&entity.EventEmployee{},
		&entity.SurveyTemplate{},
		&entity.AnswerType{},
		&entity.Question{},
		&entity.QuestionOption{},
		&entity.SurveyResponse{},
	)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("Migration success")
	}

	answerTypes := []entity.AnswerType{
		{
			Name: "Multiple Choice",
		},
		{
			Name: "Short Answer",
		},
		{
			Name: "Long Answer",
		},
		{
			Name: "Checkbox",
		},
		{
			Name: "Dropdown",
		},
		{
			Name: "Rating",
		},
		{
			Name: "Link",
		}, 
		{
			Name: "Attachment",
		},
	}

	for _, answerType := range answerTypes {
		err = db.Create(&answerType).Error
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Info("Seed AnswerType success")
}
