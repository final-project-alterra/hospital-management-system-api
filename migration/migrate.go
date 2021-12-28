package migration

import (
	"github.com/final-project-alterra/hospital-management-system-api/config"
	adminsData "github.com/final-project-alterra/hospital-management-system-api/features/admins/data"
)

func AutoMigrate() {
	db := config.DB

	err := db.AutoMigrate(&adminsData.Admin{})
	if err != nil {
		panic(err)
	}
}
