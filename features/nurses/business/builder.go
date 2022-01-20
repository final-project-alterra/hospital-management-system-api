package business

import (
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
	"github.com/final-project-alterra/hospital-management-system-api/features/schedules"
)

type nurseBusinessBuilder struct {
	nurseRepo        nurses.IData
	adminBusiness    admins.IBusiness
	doctorBusiness   doctors.IBusiness
	scheduleBusiness schedules.IBusiness
}

func NewNurseBusinessBuilder() *nurseBusinessBuilder {
	return &nurseBusinessBuilder{}
}

func (n *nurseBusinessBuilder) Build() nurses.IBusiness {
	nurseBusiness := &nurseBusiness{
		data:             n.nurseRepo,
		adminBusiness:    n.adminBusiness,
		doctorBusiness:   n.doctorBusiness,
		scheduleBusiness: n.scheduleBusiness,
	}

	n.nurseRepo = nil
	n.adminBusiness = nil
	n.doctorBusiness = nil
	n.scheduleBusiness = nil

	return nurseBusiness
}

func (n *nurseBusinessBuilder) SetData(data nurses.IData) *nurseBusinessBuilder {
	n.nurseRepo = data
	return n
}

func (n *nurseBusinessBuilder) SetAdminBusiness(adminBusiness admins.IBusiness) *nurseBusinessBuilder {
	n.adminBusiness = adminBusiness
	return n
}

func (n *nurseBusinessBuilder) SetDoctorBusiness(doctorBusiness doctors.IBusiness) *nurseBusinessBuilder {
	n.doctorBusiness = doctorBusiness
	return n
}

func (n *nurseBusinessBuilder) SetScheduleBusiness(scheduleBusiness schedules.IBusiness) *nurseBusinessBuilder {
	n.scheduleBusiness = scheduleBusiness
	return n
}
