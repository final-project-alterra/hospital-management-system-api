package data

import (
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
	"gorm.io/gorm"
)

type Doctor struct {
	gorm.Model
	CreatedBy int
	UpdatedBy int

	SpecialityID uint
	RoomID       uint
	Speciality   Speciality
	Room         Room

	Email     string `gorm:"type:varchar(64);unique"`
	Password  string `gorm:"type:varchar(128)"`
	Name      string `gorm:"type:varchar(64)"`
	Phone     string `gorm:"type:varchar(14)"`
	Gender    string `gorm:"type:varchar(1)"`
	BirthDate string `gorm:"type:date"`
	ImageUrl  string
	Address   string
}

type Speciality struct {
	gorm.Model
	Name    string `gorm:"type:varchar(64)"`
	Doctors []Doctor
}

type Room struct {
	gorm.Model
	Floor   string `gorm:"type:varchar(16)"`
	Code    string `gorm:"type:varchar(32);unique"`
	Doctors []Doctor
}

func (d Doctor) ToDoctorCore() doctors.DoctorCore {
	return doctors.DoctorCore{
		ID:         int(d.ID),
		CreatedBy:  d.CreatedBy,
		UpdatedBy:  d.UpdatedBy,
		Speciality: d.Speciality.ToSpecialityCore(),
		Room:       d.Room.ToRoomCore(),

		Email:     d.Email,
		Password:  d.Password,
		Name:      d.Name,
		BirthDate: d.BirthDate,
		ImageUrl:  d.ImageUrl,
		Phone:     d.Phone,
		Address:   d.Address,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		Gender:    d.Gender,
	}
}

func ToSliceDoctorCore(d []Doctor) []doctors.DoctorCore {
	doctors := make([]doctors.DoctorCore, len(d))
	for i, v := range d {
		doctors[i] = v.ToDoctorCore()
	}
	return doctors
}

func (s Speciality) ToSpecialityCore() doctors.SpecialityCore {
	return doctors.SpecialityCore{
		ID:        int(s.ID),
		Name:      s.Name,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func ToSliceSpecialityCore(s []Speciality) []doctors.SpecialityCore {
	specialities := make([]doctors.SpecialityCore, len(s))
	for i, v := range s {
		specialities[i] = v.ToSpecialityCore()
	}
	return specialities
}

func (r Room) ToRoomCore() doctors.RoomCore {
	return doctors.RoomCore{
		ID:        int(r.ID),
		Floor:     r.Floor,
		Code:      r.Code,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func ToSliceRoomCore(r []Room) []doctors.RoomCore {
	rooms := make([]doctors.RoomCore, len(r))
	for i, v := range r {
		rooms[i] = v.ToRoomCore()
	}
	return rooms
}
