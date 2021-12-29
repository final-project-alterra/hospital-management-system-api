package factory

import (
	"github.com/final-project-alterra/hospital-management-system-api/config"

	adminsBusiness "github.com/final-project-alterra/hospital-management-system-api/features/admins/business"
	adminsData "github.com/final-project-alterra/hospital-management-system-api/features/admins/data"
	adminsPresentation "github.com/final-project-alterra/hospital-management-system-api/features/admins/presentation"

	doctorsBusiness "github.com/final-project-alterra/hospital-management-system-api/features/doctors/business"
	doctorsData "github.com/final-project-alterra/hospital-management-system-api/features/doctors/data"
	doctorsPresentation "github.com/final-project-alterra/hospital-management-system-api/features/doctors/presentation"
)

type Presenter struct {
	AdminPresentation  *adminsPresentation.AdminPresentation
	DoctorPresentation *doctorsPresentation.DoctorPresentation
}

func New() *Presenter {
	adminBuilder := adminsBusiness.NewAdminBusinessBuilder()
	doctorBuilder := doctorsBusiness.NewDoctorBusinessBuilder()

	adminData := adminsData.NewMySQLRepo(config.DB)
	doctorData := doctorsData.NewMySQLRepo(config.DB)

	adminPresentation := adminsPresentation.NewAdminPresentation(adminBuilder.SetData(adminData).Build())
	doctorPresentation := doctorsPresentation.NewDoctorPresentation(
		doctorBuilder.SetData(doctorData).SetAdminBusiness(adminBuilder.SetData(adminData).Build()),
	)

	return &Presenter{
		AdminPresentation:  adminPresentation,
		DoctorPresentation: doctorPresentation,
	}
}
