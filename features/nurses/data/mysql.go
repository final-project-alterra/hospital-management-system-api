package data

import (
	"time"

	"github.com/final-project-alterra/hospital-management-system-api/config"
	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mySQLRepo struct {
	db *gorm.DB
}

func NewMySQLRepo(db *gorm.DB) *mySQLRepo {
	return &mySQLRepo{
		db: db,
	}
}

func (r *mySQLRepo) SelectNurses() ([]nurses.NurseCore, error) {
	const op errors.Op = "nurses.data.SelectNurses"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	var nurseRecords []Nurse
	err := r.db.Find(&nurseRecords).Error
	if err != nil {
		return []nurses.NurseCore{}, errors.E(err, op, errMessage, errors.KindServerError)
	}

	return ToSliceNurseCore(nurseRecords), err
}

func (r *mySQLRepo) SelectNursesByIds(ids []int) ([]nurses.NurseCore, error) {
	const op errors.Op = "nurses.data.SelectNursesByIds"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	var nurseRecords []Nurse
	err := r.db.Where("id in (?)", ids).Find(&nurseRecords).Error
	if err != nil {
		return []nurses.NurseCore{}, errors.E(err, op, errMessage, errors.KindServerError)
	}

	return ToSliceNurseCore(nurseRecords), err
}

func (r *mySQLRepo) SelectNurseById(id int) (nurses.NurseCore, error) {
	const op errors.Op = "nurses.data.SelectNurseById"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	var nurseRecord Nurse
	err := r.db.First(&nurseRecord, id).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			errMessage = "Nurse not found"
			return nurses.NurseCore{}, errors.E(err, op, errMessage, errors.KindNotFound)
		default:
			return nurses.NurseCore{}, errors.E(err, op, errMessage, errors.KindServerError)
		}
	}

	return nurseRecord.ToNurseCore(), err
}

func (r *mySQLRepo) SelectNurseByEmail(email string) (nurses.NurseCore, error) {
	const op errors.Op = "nurses.data.SelectNurseByEmail"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	var nurseRecord Nurse
	err := r.db.Where("email = ?", email).First(&nurseRecord).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			errMessage = "Nurse not found"
			return nurses.NurseCore{}, errors.E(err, op, errMessage, errors.KindNotFound)
		default:
			return nurses.NurseCore{}, errors.E(err, op, errMessage, errors.KindServerError)
		}
	}

	return nurseRecord.ToNurseCore(), err
}

func (r *mySQLRepo) InsertNurse(nurse nurses.NurseCore) error {
	const op errors.Op = "nurses.data.InsertNurse"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	newNurse := Nurse{
		CreatedBy: nurse.CreatedBy,
		Email:     nurse.Email,
		Password:  nurse.Password,
		Name:      nurse.Name,
		BirthDate: nurse.BirthDate,
		ImageUrl:  nurse.ImageUrl,
		Phone:     nurse.Phone,
		Address:   nurse.Address,
		Gender:    nurse.Gender,
	}

	err := r.db.Create(&newNurse).Error
	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}
	return nil
}

func (r *mySQLRepo) UpdateNurse(nurse nurses.NurseCore) error {
	const op errors.Op = "nurses.data.UpdateNurse"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	updatedNurse := Nurse{
		Model: gorm.Model{
			ID:        uint(nurse.ID),
			CreatedAt: nurse.CreatedAt,
		},
		CreatedBy: nurse.CreatedBy,
		UpdatedBy: nurse.UpdatedBy,

		Email:     nurse.Email,
		Name:      nurse.Name,
		BirthDate: nurse.BirthDate,
		ImageUrl:  nurse.ImageUrl,
		Phone:     nurse.Phone,
		Address:   nurse.Address,
		Password:  nurse.Password,
		Gender:    nurse.Gender,
	}

	err := r.db.Save(&updatedNurse).Error
	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}
	return nil
}

func (r *mySQLRepo) DeleteNurseById(id int, updatedBy int) error {
	const op errors.Op = "nurses.data.DeleteNurseById"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	email := uuid.New().String()
	now := time.Now().In(config.GetTimeLoc())
	err := r.db.
		Exec("UPDATE nurses SET updated_by = ?, deleted_at = ?, email = ? WHERE id = ?", updatedBy, now, email, id).
		Error

	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}
	return nil
}
