package business

import (
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/final-project-alterra/hospital-management-system-api/features/patients"
	"github.com/final-project-alterra/hospital-management-system-api/features/schedules"
)

type patientBusinessBuilder struct {
	repo              patients.IData
	adminBusiness     admins.IBusiness
	schedulesBusiness schedules.IBusiness
}

func NewPatientBusinessBuilder() *patientBusinessBuilder {
	return &patientBusinessBuilder{}
}

func (p *patientBusinessBuilder) Build() *patientBusiness {
	business := &patientBusiness{
		data:              p.repo,
		adminBusiness:     p.adminBusiness,
		schedulesBusiness: p.schedulesBusiness,
	}

	p.repo = nil
	p.adminBusiness = nil
	p.schedulesBusiness = nil

	return business
}

func (p *patientBusinessBuilder) SetData(data patients.IData) *patientBusinessBuilder {
	p.repo = data
	return p
}

func (p *patientBusinessBuilder) SetAdminBusiness(b admins.IBusiness) *patientBusinessBuilder {
	p.adminBusiness = b
	return p
}

func (p *patientBusinessBuilder) SetScheduleBusiness(s schedules.IBusiness) *patientBusinessBuilder {
	p.schedulesBusiness = s
	return p
}
