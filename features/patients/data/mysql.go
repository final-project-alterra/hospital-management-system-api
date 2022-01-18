package data

import (
	"time"

	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/patients"
	"gorm.io/gorm"
)

type mySQLRepo struct {
	db *gorm.DB
}

func NewMySQLRepo(db *gorm.DB) *mySQLRepo {
	return &mySQLRepo{db: db}
}

func (r *mySQLRepo) SelectPatients() ([]patients.PatientCore, error) {
	const op errors.Op = "patients.data.SelectPatients"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	patientRecords := []Patient{}
	err := r.db.Find(&patientRecords).Error
	if err != nil {
		return nil, errors.E(err, op, errMessage, errors.KindServerError)
	}
	return toSlicePatientCore(patientRecords), nil
}

func (r *mySQLRepo) SelectPatientsByIds(ids []int) ([]patients.PatientCore, error) {
	const op errors.Op = "patients.data.SelectPatientsByIds"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	patientRecords := []Patient{}
	err := r.db.Where("id IN (?)", ids).Find(&patientRecords).Error
	if err != nil {
		return nil, errors.E(err, op, errMessage, errors.KindServerError)
	}
	return toSlicePatientCore(patientRecords), nil
}

func (r *mySQLRepo) SelectPatientById(id int) (patients.PatientCore, error) {
	const op errors.Op = "patients.data.SelectPatientById"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	patientRecord := Patient{}
	err := r.db.First(&patientRecord, id).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			errMessage = "Patient not found"
			return patients.PatientCore{}, errors.E(err, op, errMessage, errors.KindNotFound)
		default:
			return patients.PatientCore{}, errors.E(err, op, errMessage, errors.KindServerError)
		}
	}
	return patientRecord.toPatientCore(), nil
}

func (r *mySQLRepo) SelectPatientByNIK(nik string) (patients.PatientCore, error) {
	const op errors.Op = "patients.data.SelectPatientByNIK"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	patientRecord := Patient{}
	err := r.db.Where("nik = ?", nik).First(&patientRecord).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			errMessage = "Patient not found"
			return patients.PatientCore{}, errors.E(err, op, errMessage, errors.KindNotFound)
		default:
			return patients.PatientCore{}, errors.E(err, op, errMessage, errors.KindServerError)
		}
	}
	return patientRecord.toPatientCore(), nil
}

func (r *mySQLRepo) InsertPatient(patient patients.PatientCore) error {
	const op errors.Op = "patients.data.SelectPatientByNIK"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	newPatientRecord := Patient{
		CreatedBy: patient.CreatedBy,
		NIK:       patient.NIK,
		Name:      patient.Name,
		BirthDate: patient.BirthDate,
		Phone:     patient.Phone,
		Address:   patient.Address,
		Gender:    patient.Gender,
	}

	err := r.db.Create(&newPatientRecord).Error
	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}
	return nil
}

func (r *mySQLRepo) UpdatePatient(patient patients.PatientCore) error {
	const op errors.Op = "patients.data.UpdatePatient"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	updatedPatientRecord := Patient{
		Model: gorm.Model{
			ID:        uint(patient.ID),
			CreatedAt: patient.CreatedAt,
		},
		UpdatedBy: patient.UpdatedBy,
		CreatedBy: patient.CreatedBy,
		NIK:       patient.NIK,
		Name:      patient.Name,
		BirthDate: patient.BirthDate,
		Phone:     patient.Phone,
		Address:   patient.Address,
		Gender:    patient.Gender,
	}

	err := r.db.Save(&updatedPatientRecord).Error
	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}
	return nil
}

func (r *mySQLRepo) DeletePatientById(id int, updatedBy int) error {
	const op errors.Op = "patients.data.DeletePatientById"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	err := r.db.
		Exec("UPDATE patients SET deleted_at = ?, updated_by = ? WHERE id = ?", time.Now(), updatedBy, id).
		Error

	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}
	return nil
}
