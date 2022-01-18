package doctors

import "time"

type DoctorCore struct {
	ID         int
	CreatedBy  int
	UpdatedBy  int
	Speciality SpecialityCore
	Room       RoomCore

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

type SpecialityCore struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RoomCore struct {
	ID        int
	Floor     string
	Code      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IBusiness interface {
	FindDoctors() ([]DoctorCore, error)
	FindDoctosrByIds(ids []int) ([]DoctorCore, error)
	FindDoctorById(id int) (DoctorCore, error)
	FindDoctorByEmail(email string) (DoctorCore, error)
	CreateDoctor(doctor DoctorCore) error
	EditDoctor(doctor DoctorCore) error
	EditDoctorPassword(id int, updatedBy int, oldPassword string, newPassword string) error
	RemoveDoctorById(id int, updatedBy int) error

	FindSpecialities() ([]SpecialityCore, error)
	FindSpecialityById(id int) (SpecialityCore, error)
	CreateSpeciality(speciality SpecialityCore) error
	EditSpeciality(speciality SpecialityCore) error
	RemoveSpeciality(id int) error

	FindRooms() ([]RoomCore, error)
	CreateRoom(room RoomCore) error
	EditRoom(room RoomCore) error
	RemoveRoomById(id int) error
}

type IData interface {
	// needs join with specialities & rooms
	SelectDoctors() ([]DoctorCore, error)
	SelectDoctorsByIds(ids []int) ([]DoctorCore, error) // used by shedules, include speicality & room
	SelectDoctorsBySpecialityId(id int) ([]DoctorCore, error)
	SelectDoctorsByRoomId(id int) ([]DoctorCore, error)
	SelectDoctorById(id int) (DoctorCore, error)

	SelectDoctorByEmail(email string) (DoctorCore, error)
	InsertDoctor(doctor DoctorCore) error
	UpdateDoctor(doctor DoctorCore) error
	DeleteDoctorById(id int, updatedBy int) error

	SelectSpecialities() ([]SpecialityCore, error)
	SelectSpecialityById(id int) (SpecialityCore, error)
	InsertSpeciality(speciality SpecialityCore) error
	UpdateSpeciality(speciality SpecialityCore) error
	DeleteSpecialityId(id int) error

	SelectRooms() ([]RoomCore, error)
	SelectRoomById(id int) (RoomCore, error)
	SelectRoomByCode(code string) (RoomCore, error)
	InsertRoom(room RoomCore) error
	UpdateRoom(room RoomCore) error
	DeleteRoomById(id int) error
}
