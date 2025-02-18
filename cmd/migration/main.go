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
	)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("Migration success")
	}
}
