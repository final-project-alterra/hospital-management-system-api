package factory

import (
	"github.com/final-project-alterra/hospital-management-system-api/config"

	adminsBusiness "github.com/final-project-alterra/hospital-management-system-api/features/admins/business"
	adminsData "github.com/final-project-alterra/hospital-management-system-api/features/admins/data"
	adminsPresentation "github.com/final-project-alterra/hospital-management-system-api/features/admins/presentation"

	doctorsBusiness "github.com/final-project-alterra/hospital-management-system-api/features/doctors/business"
	doctorsData "github.com/final-project-alterra/hospital-management-system-api/features/doctors/data"
	doctorsPresentation "github.com/final-project-alterra/hospital-management-system-api/features/doctors/presentation"

	nursesBusiness "github.com/final-project-alterra/hospital-management-system-api/features/nurses/business"
	nursesData "github.com/final-project-alterra/hospital-management-system-api/features/nurses/data"
	nursesPresentation "github.com/final-project-alterra/hospital-management-system-api/features/nurses/presentation"

	authsBusiness "github.com/final-project-alterra/hospital-management-system-api/features/auth/business"
	authsPresentation "github.com/final-project-alterra/hospital-management-system-api/features/auth/presentation"

	patientsBusiness "github.com/final-project-alterra/hospital-management-system-api/features/patients/business"
	patientsData "github.com/final-project-alterra/hospital-management-system-api/features/patients/data"
	patientsPresentation "github.com/final-project-alterra/hospital-management-system-api/features/patients/presentation"
)

type Presenter struct {
	AuthPresentation    *authsPresentation.AuthPresetation
	AdminPresentation   *adminsPresentation.AdminPresentation
	DoctorPresentation  *doctorsPresentation.DoctorPresentation
	NursePresentation   *nursesPresentation.NursePresentation
	PatientPresentation *patientsPresentation.PatientPresentation
}

func New() *Presenter {
	adminBuilder := adminsBusiness.NewAdminBusinessBuilder()
	doctorBuilder := doctorsBusiness.NewDoctorBusinessBuilder()
	nurseBuilder := nursesBusiness.NewNurseBusinessBuilder()
	authBuilder := authsBusiness.NewAuthBusinessBuilder()
	patientBuilder := patientsBusiness.NewPatientBusinessBuilder()

	adminData := adminsData.NewMySQLRepo(config.DB)
	doctorData := doctorsData.NewMySQLRepo(config.DB)
	nurseData := nursesData.NewMySQLRepo(config.DB)
	patientData := patientsData.NewMySQLRepo(config.DB)

	adminBusiness := adminBuilder.SetData(adminData).Build()
	doctorBusiness := doctorBuilder.SetData(doctorData).SetAdminBusiness(adminBusiness).Build()
	nurseBusiness := nurseBuilder.SetData(nurseData).SetAdminBusiness(adminBusiness).Build()
	patientBusiness := patientBuilder.SetData(patientData).SetAdminBusiness(adminBusiness).Build()
	authBusiness := authBuilder.
		SetAdminBusiness(adminBusiness).
		SetDoctorBusiness(doctorBusiness).
		SetNurseBusiness(nurseBusiness).
		Build()

	adminPresentation := adminsPresentation.NewAdminPresentation(adminBusiness)
	doctorPresentation := doctorsPresentation.NewDoctorPresentation(doctorBusiness)
	nursePresentation := nursesPresentation.NewNursePresentation(nurseBusiness)
	patientPresentation := patientsPresentation.NewPatientPresentation(patientBusiness)
	authPresentation := authsPresentation.NewAuthPresentation(authBusiness)

	return &Presenter{
		AuthPresentation:    authPresentation,
		AdminPresentation:   adminPresentation,
		DoctorPresentation:  doctorPresentation,
		NursePresentation:   nursePresentation,
		PatientPresentation: patientPresentation,
	}
}
