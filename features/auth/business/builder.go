package business

import (
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/final-project-alterra/hospital-management-system-api/features/auth"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
)

type authBusinessBuilder struct {
	adminBusiness  admins.IBusiness
	doctorBusiness doctors.IBusiness
	nurseBusiness  nurses.IBusiness
}

func NewAuthBusinessBuilder() *authBusinessBuilder {
	return &authBusinessBuilder{}
}

func (a *authBusinessBuilder) Build() auth.IBusiness {
	authBusiness := &authBusiness{
		adminBusiness:  a.adminBusiness,
		doctorBusiness: a.doctorBusiness,
		nurseBusiness:  a.nurseBusiness,
	}

	a.adminBusiness = nil
	a.doctorBusiness = nil
	a.nurseBusiness = nil

	return authBusiness
}

func (a *authBusinessBuilder) SetAdminBusiness(ab admins.IBusiness) *authBusinessBuilder {
	a.adminBusiness = ab
	return a
}

func (a *authBusinessBuilder) SetDoctorBusiness(db doctors.IBusiness) *authBusinessBuilder {
	a.doctorBusiness = db
	return a
}

func (a *authBusinessBuilder) SetNurseBusiness(nb nurses.IBusiness) *authBusinessBuilder {
	a.nurseBusiness = nb
	return a
}
