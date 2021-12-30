package business

import (
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/final-project-alterra/hospital-management-system-api/features/patients"
)

type patientBusinessBuilder struct {
	patientBusiness
}

func NewPatientBusinessBuilder() *patientBusinessBuilder {
	return &patientBusinessBuilder{}
}

func (p *patientBusinessBuilder) Build() *patientBusiness {
	business := p.patientBusiness
	p.patientBusiness = patientBusiness{}
	return &business
}

func (p *patientBusinessBuilder) SetData(data patients.IData) *patientBusinessBuilder {
	p.data = data
	return p
}

func (p *patientBusinessBuilder) SetAdminBusiness(b admins.IBusiness) *patientBusinessBuilder {
	p.adminBusiness = b
	return p
}
