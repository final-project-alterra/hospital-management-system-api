package business

import (
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
	"github.com/final-project-alterra/hospital-management-system-api/features/schedules"
)

type doctorBusinessBuilder struct {
	doctorRepo       doctors.IData
	adminBusiness    admins.IBusiness
	nurseBusiness    nurses.IBusiness
	scheduleBusiness schedules.IBusiness
}

func NewDoctorBusinessBuilder() *doctorBusinessBuilder {
	return &doctorBusinessBuilder{}
}

func (b *doctorBusinessBuilder) SetData(data doctors.IData) *doctorBusinessBuilder {
	b.doctorRepo = data
	return b
}

func (b *doctorBusinessBuilder) SetAdminBusiness(ab admins.IBusiness) *doctorBusinessBuilder {
	b.adminBusiness = ab
	return b
}

func (b *doctorBusinessBuilder) SetNurseBusiness(nb nurses.IBusiness) *doctorBusinessBuilder {
	b.nurseBusiness = nb
	return b
}

func (b *doctorBusinessBuilder) SetScheduleBusiness(sb schedules.IBusiness) *doctorBusinessBuilder {
	b.scheduleBusiness = sb
	return b
}

func (b *doctorBusinessBuilder) Build() doctors.IBusiness {
	doctorBusiness := &doctorBusiness{
		data:             b.doctorRepo,
		adminBusiness:    b.adminBusiness,
		nurseBusiness:    b.nurseBusiness,
		scheduleBusiness: b.scheduleBusiness,
	}

	b.doctorRepo = nil
	b.adminBusiness = nil
	b.nurseBusiness = nil
	b.scheduleBusiness = nil

	return doctorBusiness
}
