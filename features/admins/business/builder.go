package business

import (
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
)

type adminBusinessBuilder struct {
	adminRepo      admins.IData
	doctorBusiness doctors.IBusiness
	nurseBusiness  nurses.IBusiness
}

func NewAdminBusinessBuilder() *adminBusinessBuilder {
	return &adminBusinessBuilder{}
}

func (b *adminBusinessBuilder) SetData(data admins.IData) *adminBusinessBuilder {
	b.adminRepo = data
	return b
}
func (b *adminBusinessBuilder) SetDoctorBusiness(db doctors.IBusiness) *adminBusinessBuilder {
	b.doctorBusiness = db
	return b
}
func (b *adminBusinessBuilder) SetNurseBusiness(nb nurses.IBusiness) *adminBusinessBuilder {
	b.nurseBusiness = nb
	return b
}

func (b *adminBusinessBuilder) Build() admins.IBusiness {
	adminBusiness := &adminBusiness{
		data:           b.adminRepo,
		doctorBusiness: b.doctorBusiness,
		nurseBusiness:  b.nurseBusiness,
	}

	b.adminRepo = nil
	b.doctorBusiness = nil
	b.nurseBusiness = nil

	return adminBusiness
}
