package business

import (
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
)

type doctorBusinessBuilder struct {
	doctorBusiness
}

func NewDoctorBusinessBuilder() *doctorBusinessBuilder {
	return &doctorBusinessBuilder{}
}

func (b *doctorBusinessBuilder) SetData(data doctors.IData) *doctorBusinessBuilder {
	b.data = data
	return b
}

func (b *doctorBusinessBuilder) SetAdminBusiness(ab admins.IBusiness) *doctorBusinessBuilder {
	b.adminBusiness = ab
	return b
}

func (b *doctorBusinessBuilder) Build() doctors.IBusiness {
	business := b.doctorBusiness
	b.doctorBusiness = doctorBusiness{}

	return &business
}
