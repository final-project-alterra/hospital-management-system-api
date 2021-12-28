package factory

import (
	"github.com/final-project-alterra/hospital-management-system-api/config"

	adminsBusiness "github.com/final-project-alterra/hospital-management-system-api/features/admins/business"
	adminsData "github.com/final-project-alterra/hospital-management-system-api/features/admins/data"
	adminsPresentation "github.com/final-project-alterra/hospital-management-system-api/features/admins/presentation"
)

type Presenter struct {
	AdminPresentation *adminsPresentation.AdminPresentation
}

func New() *Presenter {
	adminBuilder := adminsBusiness.NewAdminBusinessBuilder()

	adminData := adminsData.NewMySQLRepo(config.DB)
	adminPresentation := adminsPresentation.NewAdminPresentation(adminBuilder.SetData(adminData).Build())

	return &Presenter{
		AdminPresentation: adminPresentation,
	}
}
