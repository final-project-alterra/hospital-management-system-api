package data

import (
	"time"

	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
	"gorm.io/gorm"
)

type mySQLRepo struct {
	db *gorm.DB
}

func NewMySQLRepo(db *gorm.DB) *mySQLRepo {
	return &mySQLRepo{db}
}

// needs join with specialities & rooms
func (r *mySQLRepo) SelectDoctors() ([]doctors.DoctorCore, error) {
	const op errors.Op = "doctors.data.SelectDoctors"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	var doctorRecords []Doctor
	err := r.db.Preload("Speciality").Preload("Room").Find(&doctorRecords).Error
	if err != nil {
		return []doctors.DoctorCore{}, errors.E(err, op, errMessage, errors.KindServerError)
	}
	return ToSliceDoctorCore(doctorRecords), nil
}
func (r *mySQLRepo) SelectDoctorsByIds(ids []int) ([]doctors.DoctorCore, error) {
	const op errors.Op = "doctors.data.SelectDoctorsByIds"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	var doctorRecords []Doctor
	err := r.db.Preload("Speciality").Preload("Room").Find(&doctorRecords, ids).Error
	if err != nil {
		return []doctors.DoctorCore{}, errors.E(err, op, errMessage, errors.KindServerError)
	}
	return ToSliceDoctorCore(doctorRecords), nil
}
func (r *mySQLRepo) SelectDoctorsBySpecialityId(id int) ([]doctors.DoctorCore, error) {
	const op errors.Op = "doctors.data.SelectDoctorsBySpecialityId"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	var doctorRecords []Doctor
	err := r.db.Where("speciality_id = ?", id).Preload("Speciality").Preload("Room").Find(&doctorRecords).Error
	if err != nil {
		return []doctors.DoctorCore{}, errors.E(err, op, errMessage, errors.KindServerError)
	}
	return ToSliceDoctorCore(doctorRecords), nil
}
func (r *mySQLRepo) SelectDoctorsByRoomId(id int) ([]doctors.DoctorCore, error) {
	const op errors.Op = "doctors.data.SelectDoctorsByRoomId"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	var doctorRecords []Doctor
	err := r.db.Where("room_id = ?", id).Preload("Speciality").Preload("Room").Find(&doctorRecords).Error
	if err != nil {
		return []doctors.DoctorCore{}, errors.E(err, op, errMessage, errors.KindServerError)
	}
	return ToSliceDoctorCore(doctorRecords), nil
}
func (r *mySQLRepo) SelectDoctorById(id int) (doctors.DoctorCore, error) {
	const op errors.Op = "doctors.data.SelectDoctorById"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	var doctorRecord Doctor
	err := r.db.Preload("Speciality").Preload("Room").First(&doctorRecord, id).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			errMessage = "Doctor not found"
			return doctors.DoctorCore{}, errors.E(err, op, errMessage, errors.KindNotFound)

		default:
			return doctors.DoctorCore{}, errors.E(err, op, errMessage, errors.KindServerError)
		}
	}
	return doctorRecord.ToDoctorCore(), nil
}

func (r *mySQLRepo) SelectDoctorByEmail(email string) (doctors.DoctorCore, error) {
	const op errors.Op = "doctors.data.SelectDoctorById"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	var doctorRecord Doctor
	err := r.db.Where("email = ?", email).Preload("Speciality").Preload("Room").First(&doctorRecord).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			errMessage = "Doctor not found"
			return doctors.DoctorCore{}, errors.E(err, op, errMessage, errors.KindNotFound)

		default:
			return doctors.DoctorCore{}, errors.E(err, op, errMessage, errors.KindServerError)
		}
	}
	return doctorRecord.ToDoctorCore(), nil
}
func (r *mySQLRepo) InsertDoctor(doctor doctors.DoctorCore) error {
	const op errors.Op = "doctors.data.InsertDoctor"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	doctorRecord := Doctor{
		SpecialityID: uint(doctor.Speciality.ID),
		RoomID:       uint(doctor.Room.ID),
		CreatedBy:    doctor.CreatedBy,

		Name:     doctor.Name,
		Email:    doctor.Email,
		Password: doctor.Password,
		Phone:    doctor.Phone,
		Gender:   doctor.Gender,
		Age:      doctor.Age,
		ImageUrl: doctor.ImageUrl,
		Address:  doctor.Address,
	}
	err := r.db.Create(&doctorRecord).Error
	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}
	return nil
}
func (r *mySQLRepo) UpdateDoctor(doctor doctors.DoctorCore) error {
	const op errors.Op = "doctors.data.UpdateDoctor"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	updatedDoctor := Doctor{
		Model: gorm.Model{
			ID:        uint(doctor.ID),
			CreatedAt: doctor.CreatedAt,
		},
		SpecialityID: uint(doctor.Speciality.ID),
		RoomID:       uint(doctor.Room.ID),
		CreatedBy:    doctor.CreatedBy,
		UpdatedBy:    doctor.UpdatedBy,

		Name:     doctor.Name,
		Email:    doctor.Email,
		Password: doctor.Password,
		Phone:    doctor.Phone,
		Gender:   doctor.Gender,
		Age:      doctor.Age,
		ImageUrl: doctor.ImageUrl,
		Address:  doctor.Address,
	}

	err := r.db.Save(updatedDoctor).Error
	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}
	return nil
}
func (r *mySQLRepo) DeleteDoctorById(id int, updatedBy int) error {
	const op errors.Op = "doctors.data.DeleteDoctorById"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	err := r.db.
		Exec("UPDATE doctors SET updated_by = ?, deleted_at = ? WHERE id = ?", updatedBy, time.Now(), id).
		Error

	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}
	return nil
}

func (r *mySQLRepo) SelectSpecialities() ([]doctors.SpecialityCore, error)
func (r *mySQLRepo) SelectSpecialityById(id int) (doctors.SpecialityCore, error)
func (r *mySQLRepo) InsertSpeciality(speciality doctors.SpecialityCore) error
func (r *mySQLRepo) UpdateSpeciality(speciality doctors.SpecialityCore) error
func (r *mySQLRepo) DeleteSpecialityId(id int) error

func (r *mySQLRepo) SelectRooms() ([]doctors.RoomCore, error)
func (r *mySQLRepo) SelectRoomById(id int) (doctors.RoomCore, error)
func (r *mySQLRepo) SelectRoomByCode(code string) (doctors.RoomCore, error)
func (r *mySQLRepo) InsertRoom(room doctors.RoomCore) error
func (r *mySQLRepo) UpdateRoom(room doctors.RoomCore) error
func (r *mySQLRepo) DeleteRoomById(id int) error
