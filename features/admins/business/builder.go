package business

import "github.com/final-project-alterra/hospital-management-system-api/features/admins"

type adminBusinessBuilder struct {
	adminRepo admins.IData
}

func NewAdminBusinessBuilder() *adminBusinessBuilder {
	return &adminBusinessBuilder{}
}

func (b *adminBusinessBuilder) SetData(data admins.IData) *adminBusinessBuilder {
	b.adminRepo = data
	return b
}

func (b *adminBusinessBuilder) Build() admins.IBusiness {
	adminBusiness := &adminBusiness{
		data: b.adminRepo,
	}

	b.adminRepo = nil

	return adminBusiness
}
