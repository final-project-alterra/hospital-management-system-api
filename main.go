package main

import (
	"path"

	"github.com/final-project-alterra/hospital-management-system-api/config"
	"github.com/final-project-alterra/hospital-management-system-api/migration"
	"github.com/final-project-alterra/hospital-management-system-api/routes"
	"github.com/final-project-alterra/hospital-management-system-api/utils/project"
)

func main() {
	mainDirectory, err := project.MainDirectory()
	if err != nil {
		panic(err)
	}

	config.LoadENV(path.Join(mainDirectory, ".env"))
	config.InitTimeLoc(config.ENV.TIMEZONE)
	config.ConnectDB()
	migration.AutoMigrate()

	e := routes.SetupRoutes()
	e.Logger.Fatal(e.Start(":" + config.ENV.PORT))
}
