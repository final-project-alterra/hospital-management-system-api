package business

import (
	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/final-project-alterra/hospital-management-system-api/features/patients"
)

type patientBusiness struct {
	data          patients.IData
	adminBusiness admins.IBusiness
}

func (p *patientBusiness) FindPatients() ([]patients.PatientCore, error) {
	const op errors.Op = "patients.business.FindPatients"

	patientsData, err := p.data.SelectPatients()
	if err != nil {
		return []patients.PatientCore{}, errors.E(err, op)
	}
	return patientsData, nil
}

func (p *patientBusiness) FindPatientsByIds(ids []int) ([]patients.PatientCore, error) {
	const op errors.Op = "patients.business.FindPatientsByIds"

	patientsData, err := p.data.SelectPatientsByIds(ids)
	if err != nil {
		return []patients.PatientCore{}, errors.E(err, op)
	}
	return patientsData, nil
}

func (p *patientBusiness) FindPatientById(id int) (patients.PatientCore, error) {
	const op errors.Op = "patients.business.FindPatientById"

	patientData, err := p.data.SelectPatientById(id)
	if err != nil {
		return patients.PatientCore{}, errors.E(err, op)
	}
	return patientData, nil
}

func (p *patientBusiness) CreatePatient(patient patients.PatientCore) error {
	const op errors.Op = "patients.business.CreatePatient"
	var errMessage errors.ErrClientMessage

	_, err := p.adminBusiness.FindAdminById(patient.CreatedBy)
	if err != nil {
		return errors.E(err, op)
	}

	_, err = p.data.SelectPatientByNIK(patient.NIK)
	if err == nil {
		err = errors.New("NIK already exists")
		errMessage = "NIK already exists"
		return errors.E(err, op, errMessage, errors.KindUnprocessable)
	}
	if errors.Kind(err) != errors.KindNotFound {
		return errors.E(err, op)
	}

	err = p.data.InsertPatient(patient)
	if err != nil {
		return errors.E(err, op)
	}

	return nil
}

func (p *patientBusiness) EditPatient(patient patients.PatientCore) error {
	const op errors.Op = "patients.business.EditPatient"

	_, err := p.adminBusiness.FindAdminById(patient.CreatedBy)
	if err != nil {
		return errors.E(err, op)
	}

	existingPatient, err := p.data.SelectPatientById(patient.ID)
	if err != nil {
		return errors.E(err, op)
	}

	existingPatient.UpdatedBy = patient.UpdatedBy
	existingPatient.Name = patient.Name
	existingPatient.Age = patient.Age
	existingPatient.Phone = patient.Phone
	existingPatient.Address = patient.Address
	existingPatient.Gender = patient.Gender

	err = p.data.UpdatePatient(existingPatient)
	if err != nil {
		return errors.E(err, op)
	}

	return nil
}

func (p *patientBusiness) RemovePatientById(id int, updatedBy int) error {
	const op errors.Op = "patients.business.EditPatient"

	_, err := p.adminBusiness.FindAdminById(updatedBy)
	if err != nil {
		return errors.E(err, op)
	}

	err = p.data.DeletePatientById(id, updatedBy)
	if err != nil {
		return errors.E(err, op)
	}

	return nil
}
