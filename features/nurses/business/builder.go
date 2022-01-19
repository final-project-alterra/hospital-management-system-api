package business

import (
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
	"github.com/final-project-alterra/hospital-management-system-api/features/schedules"
)

type nurseBusinessBuilder struct {
	nurseRepo        nurses.IData
	adminBusiness    admins.IBusiness
	scheduleBusiness schedules.IBusiness
}

func NewNurseBusinessBuilder() *nurseBusinessBuilder {
	return &nurseBusinessBuilder{}
}

func (n *nurseBusinessBuilder) Build() nurses.IBusiness {
	nurseBusiness := &nurseBusiness{
		data:             n.nurseRepo,
		adminBusiness:    n.adminBusiness,
		scheduleBusiness: n.scheduleBusiness,
	}

	n.nurseRepo = nil
	n.adminBusiness = nil
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

func (n *nurseBusinessBuilder) SetScheduleBusiness(scheduleBusiness schedules.IBusiness) *nurseBusinessBuilder {
	n.scheduleBusiness = scheduleBusiness
	return n
}
