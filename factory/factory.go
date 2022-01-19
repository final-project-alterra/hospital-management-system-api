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

	schedulesBusiness "github.com/final-project-alterra/hospital-management-system-api/features/schedules/business"
	schedulesData "github.com/final-project-alterra/hospital-management-system-api/features/schedules/data"
	schedulesPresentation "github.com/final-project-alterra/hospital-management-system-api/features/schedules/presentation"
)

type Presenter struct {
	AuthPresentation     *authsPresentation.AuthPresetation
	AdminPresentation    *adminsPresentation.AdminPresentation
	DoctorPresentation   *doctorsPresentation.DoctorPresentation
	NursePresentation    *nursesPresentation.NursePresentation
	PatientPresentation  *patientsPresentation.PatientPresentation
	SchedulePresentation *schedulesPresentation.SchedulePresentation
}

func New() *Presenter {
	adminBuilder := adminsBusiness.NewAdminBusinessBuilder()
	doctorBuilder := doctorsBusiness.NewDoctorBusinessBuilder()
	nurseBuilder := nursesBusiness.NewNurseBusinessBuilder()
	authBuilder := authsBusiness.NewAuthBusinessBuilder()
	patientBuilder := patientsBusiness.NewPatientBusinessBuilder()
	scheduleBuilder := schedulesBusiness.NewScheduleBusinessBuilder()

	adminData := adminsData.NewMySQLRepo(config.DB)
	doctorData := doctorsData.NewMySQLRepo(config.DB)
	nurseData := nursesData.NewMySQLRepo(config.DB)
	patientData := patientsData.NewMySQLRepo(config.DB)
	scheduleData := schedulesData.NewMySQLRepo(config.DB)

	pureScheduleBusiness := scheduleBuilder.SetData(scheduleData).Build()

	adminBusiness := adminBuilder.SetData(adminData).Build()
	doctorBusiness := doctorBuilder.
		SetData(doctorData).
		SetAdminBusiness(adminBusiness).
		SetScheduleBusiness(pureScheduleBusiness).
		Build()
	nurseBusiness := nurseBuilder.
		SetData(nurseData).
		SetAdminBusiness(adminBusiness).
		SetScheduleBusiness(pureScheduleBusiness).
		Build()
	patientBusiness := patientBuilder.
		SetData(patientData).
		SetAdminBusiness(adminBusiness).
		SetScheduleBusiness(pureScheduleBusiness).
		Build()
	authBusiness := authBuilder.
		SetAdminBusiness(adminBusiness).
		SetDoctorBusiness(doctorBusiness).
		SetNurseBusiness(nurseBusiness).
		Build()
	scheduleBusiness := scheduleBuilder.
		SetData(scheduleData).
		SetDoctorBusiness(doctorBusiness).
		SetNurseBusiness(nurseBusiness).
		SetPatientBusiness(patientBusiness).
		Build()

	adminPresentation := adminsPresentation.NewAdminPresentation(adminBusiness)
	doctorPresentation := doctorsPresentation.NewDoctorPresentation(doctorBusiness)
	nursePresentation := nursesPresentation.NewNursePresentation(nurseBusiness)
	patientPresentation := patientsPresentation.NewPatientPresentation(patientBusiness)
	authPresentation := authsPresentation.NewAuthPresentation(authBusiness)
	schedulePresentation := schedulesPresentation.NewSchedulePresentation(scheduleBusiness)

	return &Presenter{
		AuthPresentation:     authPresentation,
		AdminPresentation:    adminPresentation,
		DoctorPresentation:   doctorPresentation,
		NursePresentation:    nursePresentation,
		PatientPresentation:  patientPresentation,
		SchedulePresentation: schedulePresentation,
	}
}
