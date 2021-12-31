package migration

import (
	"github.com/final-project-alterra/hospital-management-system-api/config"
	adminsData "github.com/final-project-alterra/hospital-management-system-api/features/admins/data"
	doctorsData "github.com/final-project-alterra/hospital-management-system-api/features/doctors/data"
	nursesData "github.com/final-project-alterra/hospital-management-system-api/features/nurses/data"
)

func AutoMigrate() {
	db := config.DB

	err := db.AutoMigrate(
		&adminsData.Admin{},
		&doctorsData.Room{},
		&doctorsData.Speciality{},
		&doctorsData.Doctor{},
		&nursesData.Nurse{},
	)

	if err != nil {
		panic(err)
	}
}
