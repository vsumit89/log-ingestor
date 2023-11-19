package main

import (
	"logswift/internal/app"
	"logswift/pkg/logger"
	"os"

	"gopkg.in/yaml.v2"
)

func main() {
	log := logger.GetInstance()

	log.Info("starting the application")

	log.Info("reading config file")

	file, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Error("error reading config file", "error", err.Error())
		return
	}

	err = yaml.Unmarshal(file, &app.AppCfg)
	if err != nil {
		log.Error("error unmarshalling config file", "error", err.Error())
		return
	}
	log.Info("config file read successfully")

	appInstance := app.NewApp()

	err = appInstance.Start()
	if err != nil {
		log.Error("Error starting the application", "error", err.Error())
		return
	}
}
