package admins

import "time"

type AdminCore struct {
	ID        int
	Email     string
	Password  string
	Name      string
	Age       int
	ImageUrl  string
	Phone     string
	Address   string
	Gender    rune
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IBusiness interface {
	FindAdmins() ([]AdminCore, error)
	FindAdminById(id int) (AdminCore, error)
	FindAdminByEmail(email string) (AdminCore, error)
	CreateAdmin(admin AdminCore) error
	EditAdmin(admin AdminCore) error
	EditAdminPassword(id int, oldPassword string, newPassword string) error
	RemoveAdminById(id int) error
}

type IData interface {
	SelectAdmins() ([]AdminCore, error)
	SelectAdminById(id int) (AdminCore, error)
	SelectAdminByEmail(email string) (AdminCore, error)
	InsertAdmin(admin AdminCore) error
	UpdateAdmin(admin AdminCore) error
	DeleteAdminById(id int) error
}
