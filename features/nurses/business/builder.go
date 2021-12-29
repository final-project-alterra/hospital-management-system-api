package business

import (
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
)

type nurseBusinessBuilder struct {
	nurseBusiness
}

func NewNurseBusinessBuilder() *nurseBusinessBuilder {
	return &nurseBusinessBuilder{}
}

func (n *nurseBusinessBuilder) Build() nurses.IBusiness {
	business := n.nurseBusiness
	n.nurseBusiness = nurseBusiness{}
	return &business
}

func (n *nurseBusinessBuilder) SetData(data nurses.IData) *nurseBusinessBuilder {
	n.data = data
	return n
}

func (n *nurseBusinessBuilder) SetAdminBusiness(adminBusiness admins.IBusiness) *nurseBusinessBuilder {
	n.adminBusiness = adminBusiness
	return n
}
