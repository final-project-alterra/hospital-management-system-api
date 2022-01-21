package nurses

import "time"

type NurseCore struct {
	ID        int
	CreatedBy int
	UpdatedBy int
	Email     string
	Password  string
	Name      string
	BirthDate string
	ImageUrl  string
	Phone     string
	Address   string
	Gender    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IBusiness interface {
	FindNurses() ([]NurseCore, error)
	FindNursesByIds(ids []int) ([]NurseCore, error)
	FindNurseById(id int) (NurseCore, error)
	FindNurseByEmail(email string) (NurseCore, error)
	CreateNurse(nurse NurseCore) error
	EditNurse(nurse NurseCore) error
	EditNurseImageProfile(nurse NurseCore) error
	EditNursePassword(id int, updatedBy int, oldPassword string, newPassword string) error
	RemoveNurseById(id int, updatedBy int) error
}

type IData interface {
	SelectNurses() ([]NurseCore, error)
	SelectNursesByIds(ids []int) ([]NurseCore, error)
	SelectNurseById(id int) (NurseCore, error)
	SelectNurseByEmail(email string) (NurseCore, error)
	InsertNurse(nurse NurseCore) error
	UpdateNurse(nurse NurseCore) error
	DeleteNurseById(id int, updatedBy int) error
}
