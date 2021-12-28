package business

import "github.com/final-project-alterra/hospital-management-system-api/features/admins"

type adminBusinessBuilder struct {
	adminBusiness
}

func NewAdminBusinessBuilder() *adminBusinessBuilder {
	return &adminBusinessBuilder{}
}

func (b *adminBusinessBuilder) SetData(data admins.IData) *adminBusinessBuilder {
	b.data = data
	return b
}

func (b *adminBusinessBuilder) Build() admins.IBusiness {
	business := b.adminBusiness
	b.adminBusiness = adminBusiness{}

	return &business
}
