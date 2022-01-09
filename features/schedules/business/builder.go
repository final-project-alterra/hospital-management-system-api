package business

import (
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
	"github.com/final-project-alterra/hospital-management-system-api/features/patients"
	"github.com/final-project-alterra/hospital-management-system-api/features/schedules"
)

type scheduleBusinessBuilder struct {
	repo            schedules.IData
	doctorBusiness  doctors.IBusiness
	nurseBusiness   nurses.IBusiness
	patientBusiness patients.IBusiness
}

func NewScheduleBusinessBuilder() *scheduleBusinessBuilder {
	return &scheduleBusinessBuilder{}
}

func (b *scheduleBusinessBuilder) SetData(repo schedules.IData) *scheduleBusinessBuilder {
	b.repo = repo
	return b
}

func (b *scheduleBusinessBuilder) SetPatientBusiness(p patients.IBusiness) *scheduleBusinessBuilder {
	b.patientBusiness = p
	return b
}

func (b *scheduleBusinessBuilder) SetDoctorBusiness(d doctors.IBusiness) *scheduleBusinessBuilder {
	b.doctorBusiness = d
	return b
}

func (b *scheduleBusinessBuilder) SetNurseBusiness(n nurses.IBusiness) *scheduleBusinessBuilder {
	b.nurseBusiness = n
	return b
}

func (b *scheduleBusinessBuilder) Build() *scheduleBusiness {
	business := &scheduleBusiness{
		data:            b.repo,
		patientBusiness: b.patientBusiness,
		doctorBusiness:  b.doctorBusiness,
		nurseBusiness:   b.nurseBusiness,
	}
	b.repo = nil
	b.doctorBusiness = nil
	b.nurseBusiness = nil
	b.patientBusiness = nil

	return business
}
