package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/entity"
)

func LoadJSON(filename string) ([]entity.University, error) {
	var universities []entity.University
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &universities)
	if err != nil {
		return nil, err
	}

	return universities, nil
}

func main() {
	viper := config.NewViper()
	log := config.NewLogrus(viper)
	db := config.NewDatabase()

	absPath, err := filepath.Abs("./cmd/university/world_universities_and_domains.json")
	if err != nil {
		log.Fatal("failed to get absolute path: ", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Fatal("file does not exist: ", absPath)
	}

	universities, err := LoadJSON(absPath)
	if err != nil {
		log.Fatal("failed to load JSON data: ", err)
	}

	for _, university := range universities {
		db.Create(&university)
	}

	log.Info("Seed University success")
}
