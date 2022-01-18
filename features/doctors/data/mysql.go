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
	err := r.db.Preload("Speciality").Preload("Room").Where("id IN (?)", ids).Find(&doctorRecords).Error
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

		Name:      doctor.Name,
		Email:     doctor.Email,
		Password:  doctor.Password,
		Phone:     doctor.Phone,
		Gender:    doctor.Gender,
		BirthDate: doctor.BirthDate,
		ImageUrl:  doctor.ImageUrl,
		Address:   doctor.Address,
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

		Name:      doctor.Name,
		Email:     doctor.Email,
		Password:  doctor.Password,
		Phone:     doctor.Phone,
		Gender:    doctor.Gender,
		BirthDate: doctor.BirthDate,
		ImageUrl:  doctor.ImageUrl,
		Address:   doctor.Address,
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

func (r *mySQLRepo) SelectSpecialities() ([]doctors.SpecialityCore, error) {
	const op errors.Op = "doctors.data.SelectSpecialities"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	var specialityRecords []Speciality
	err := r.db.Find(&specialityRecords).Error
	if err != nil {
		return []doctors.SpecialityCore{}, errors.E(err, op, errMessage, errors.KindServerError)
	}
	return ToSliceSpecialityCore(specialityRecords), nil
}
func (r *mySQLRepo) SelectSpecialityById(id int) (doctors.SpecialityCore, error) {
	const op errors.Op = "doctors.data.SelectSpecialityById"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	var specialityRecord Speciality
	err := r.db.First(&specialityRecord, id).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			errMessage = "Speciality not found"
			return doctors.SpecialityCore{}, errors.E(err, op, errMessage, errors.KindNotFound)

		default:
			return doctors.SpecialityCore{}, errors.E(err, op, errMessage, errors.KindServerError)
		}
	}
	return specialityRecord.ToSpecialityCore(), nil
}
func (r *mySQLRepo) InsertSpeciality(speciality doctors.SpecialityCore) error {
	const op errors.Op = "doctors.data.InsertSpeciality"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	specialityRecord := Speciality{
		Name: speciality.Name,
	}
	err := r.db.Create(&specialityRecord).Error
	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}
	return nil
}
func (r *mySQLRepo) UpdateSpeciality(speciality doctors.SpecialityCore) error {
	const op errors.Op = "doctors.data.UpdateSpeciality"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	updatedSpeciality := Speciality{
		Model: gorm.Model{
			ID:        uint(speciality.ID),
			CreatedAt: speciality.CreatedAt,
		},
		Name: speciality.Name,
	}

	err := r.db.Save(updatedSpeciality).Error
	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}
	return nil
}
func (r *mySQLRepo) DeleteSpecialityId(id int) error {
	const op errors.Op = "doctors.data.DeleteSpecialityId"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	err := r.db.Delete(&Speciality{}, id).Error
	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}
	return nil
}

func (r *mySQLRepo) SelectRooms() ([]doctors.RoomCore, error) {
	const op errors.Op = "doctors.data.SelectRooms"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	var roomRecords []Room
	err := r.db.Find(&roomRecords).Error
	if err != nil {
		return []doctors.RoomCore{}, errors.E(err, op, errMessage, errors.KindServerError)
	}
	return ToSliceRoomCore(roomRecords), nil
}
func (r *mySQLRepo) SelectRoomById(id int) (doctors.RoomCore, error) {
	const op errors.Op = "doctors.data.SelectRoomById"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	var roomRecord Room
	err := r.db.First(&roomRecord, id).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			errMessage = "Room not found"
			return doctors.RoomCore{}, errors.E(err, op, errMessage, errors.KindNotFound)

		default:
			return doctors.RoomCore{}, errors.E(err, op, errMessage, errors.KindServerError)
		}
	}
	return roomRecord.ToRoomCore(), nil
}
func (r *mySQLRepo) SelectRoomByCode(code string) (doctors.RoomCore, error) {
	const op errors.Op = "doctors.data.SelectRoomByCode"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	var roomRecord Room
	err := r.db.Where("code = ?", code).First(&roomRecord).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			errMessage = "Room not found"
			return doctors.RoomCore{}, errors.E(err, op, errMessage, errors.KindNotFound)

		default:
			return doctors.RoomCore{}, errors.E(err, op, errMessage, errors.KindServerError)
		}
	}
	return roomRecord.ToRoomCore(), nil
}
func (r *mySQLRepo) InsertRoom(room doctors.RoomCore) error {
	const op errors.Op = "doctors.data.InsertRoom"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	roomRecord := Room{
		Code:  room.Code,
		Floor: room.Floor,
	}
	err := r.db.Create(&roomRecord).Error
	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}
	return nil
}
func (r *mySQLRepo) UpdateRoom(room doctors.RoomCore) error {
	const op errors.Op = "doctors.data.UpdateRoom"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	updatedRoom := Room{
		Model: gorm.Model{
			ID:        uint(room.ID),
			CreatedAt: room.CreatedAt,
		},
		Code:  room.Code,
		Floor: room.Floor,
	}

	err := r.db.Save(updatedRoom).Error
	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}
	return nil
}
func (r *mySQLRepo) DeleteRoomById(id int) error {
	const op errors.Op = "doctors.data.DeleteRoomById"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	err := r.db.Delete(&Room{}, id).Error
	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}
	return nil
}
