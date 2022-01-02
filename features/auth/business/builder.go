package business

import (
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
)

type authBusinessBuilder struct {
	authBusiness
}

func NewAuthBusinessBuilder() *authBusinessBuilder {
	return &authBusinessBuilder{}
}

func (a *authBusinessBuilder) Build() *authBusiness {
	business := a.authBusiness
	a.authBusiness = authBusiness{}
	return &business
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
