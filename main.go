package main

import (
	"github.com/final-project-alterra/hospital-management-system-api/config"
	"github.com/final-project-alterra/hospital-management-system-api/migration"
	"github.com/final-project-alterra/hospital-management-system-api/routes"
)

func main() {
	path := ".env"
	config.LoadENV(path)
	config.ConnectDB()
	migration.AutoMigrate()

	e := routes.SetupRoutes()
	e.Logger.Fatal(e.Start(":" + config.ENV.PORT))
}
