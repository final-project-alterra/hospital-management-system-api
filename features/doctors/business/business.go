package business

import (
	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
)

type doctorBusiness struct {
	data          doctors.IData
	adminBusiness admins.IBusiness
}

func (d *doctorBusiness) FindDoctors() ([]doctors.DoctorCore, error) {
	const op errors.Op = "doctors.business.FindDoctors"

	doctorsData, err := d.data.SelectDoctors()
	if err != nil {
		return []doctors.DoctorCore{}, errors.E(err, op)
	}

	return doctorsData, nil
}

func (d *doctorBusiness) FindDoctosrByIds(ids []int) ([]doctors.DoctorCore, error) {
	const op errors.Op = "doctors.business.FindDoctosrByIds"

	doctorsData, err := d.data.SelectDoctorsByIds(ids)
	if err != nil {
		return []doctors.DoctorCore{}, errors.E(err, op)
	}

	return doctorsData, nil
}

func (d *doctorBusiness) FindDoctorById(id int) (doctors.DoctorCore, error) {
	const op errors.Op = "doctors.business.FindDoctorById"

	doctorData, err := d.data.SelectDoctorById(id)
	if err != nil {
		return doctors.DoctorCore{}, errors.E(err, op)
	}

	return doctorData, nil
}

func (d *doctorBusiness) FindDoctorByEmail(email string) (doctors.DoctorCore, error) {
	const op errors.Op = "doctors.business.FindDoctorByEmail"

	doctorData, err := d.data.SelectDoctorByEmail(email)
	if err != nil {
		return doctors.DoctorCore{}, errors.E(err, op)
	}

	return doctorData, nil
}

func (d *doctorBusiness) CreateDoctor(doctor doctors.DoctorCore) error {
	const op errors.Op = "doctors.business.CreateAdmin"
	var errMessage errors.ErrClientMessage

	// Check admin who is creating this new doctor, if not found or error, return error
	_, err := d.adminBusiness.FindAdminById(doctor.CreatedBy)
	if err != nil {
		switch errors.Kind(err) {
		case errors.KindNotFound:
			errMessage = "Admin who wants to create this doctor is not found"
			return errors.E(err, op, errMessage)
		default:
			return errors.E(err, op)
		}
	}

	// Check speciality, if not found or error, return error
	_, err = d.data.SelectSpecialityById(doctor.Speciality.ID)
	if err != nil {
		return errors.E(err, op)
	}

	// Check room, if not found or error, return error
	_, err = d.data.SelectRoomById(doctor.Room.ID)
	if err != nil {
		return errors.E(err, op)
	}

	// Check wheter email is already registered
	_, err = d.data.SelectDoctorByEmail(doctor.Email)
	if err == nil {
		err = errors.New("Duplicate email when createing doctor admin")
		errMessage = "Email already exists"
		return errors.E(err, op, errMessage, errors.KindUnprocessable)
	}
	if errors.Kind(err) != errors.KindNotFound {
		return errors.E(err, op)
	}

	doctor.Password, err = hash.Generate(doctor.Password)
	if err != nil {
		errMessage = "Something went wrong"
		return errors.E(err, op, errMessage, errors.KindServerError)
	}

	err = d.data.InsertDoctor(doctor)
	if err != nil {
		return errors.E(err, op)
	}
	return nil
}

func (d *doctorBusiness) EditDoctor(doctor doctors.DoctorCore) error {
	const op errors.Op = "doctors.business.EditDoctor"
	var errMessage errors.ErrClientMessage

	existingDoctor, err := d.data.SelectDoctorById(doctor.ID)
	if err != nil {
		return errors.E(err, op)
	}

	// Check speciality, if not found or error, return error
	_, err = d.data.SelectSpecialityById(doctor.Speciality.ID)
	if err != nil {
		return errors.E(err, op)
	}

	// Check room, if not found or error, return error
	_, err = d.data.SelectRoomById(doctor.Room.ID)
	if err != nil {
		return errors.E(err, op)
	}

	_, err = d.adminBusiness.FindAdminById(doctor.UpdatedBy)
	if err != nil {
		switch errors.Kind(err) {
		case errors.KindNotFound:
			errMessage = "Admin who wanted to update was not found"
			return errors.E(err, op, errMessage)
		default:
			return errors.E(err, op)
		}
	}

	existingDoctor.UpdatedBy = doctor.UpdatedBy
	existingDoctor.Name = doctor.Name
	existingDoctor.Age = doctor.Age
	existingDoctor.Address = doctor.Address
	existingDoctor.Phone = doctor.Phone
	existingDoctor.Gender = doctor.Gender

	err = d.data.UpdateDoctor(existingDoctor)
	if err != nil {
		return errors.E(err, op)
	}
	return nil
}

func (d *doctorBusiness) EditDoctorPassword(id int, updatedBy int, oldPassword string, newPassword string) error {
	const op errors.Op = "doctors.business.EditDoctorPassword"
	var errMessage errors.ErrClientMessage

	_, err := d.adminBusiness.FindAdminById(updatedBy)
	if err != nil {
		switch errors.Kind(err) {
		case errors.KindNotFound:
			errMessage = "Admin who wants to update is not found"
			return errors.E(err, op, errMessage)
		default:
			return errors.E(err, op)
		}
	}

	existingDoctor, err := d.data.SelectDoctorById(id)
	if err != nil {
		switch errors.Kind(err) {
		case errors.KindNotFound:
			errMessage = "Admin who wants to be updated is not found"
			return errors.E(err, op, errMessage)
		default:
			return errors.E(err, op)
		}
	}

	doesMatch := hash.Validate(existingDoctor.Password, oldPassword)
	if !doesMatch {
		err = errors.New("Wrong password")
		errMessage = "Wrong old password"
		return errors.E(err, op, errMessage, errors.KindUnprocessable)
	}

	existingDoctor.UpdatedBy = updatedBy
	existingDoctor.Password, err = hash.Generate(newPassword)
	if err != nil {
		return errors.E(err, op)
	}

	err = d.data.UpdateDoctor(existingDoctor)
	if err != nil {
		return errors.E(err, op)
	}

	return nil
}

// func (d *doctorBusiness) RemoveDoctorById(id int, updatedBy int) error {
// 	return nil
// }

// func (d *doctorBusiness) FindSpecialities() ([]SpecialityCore, error) {
// 	return nil

// }

// func (d *doctorBusiness) FindSpecialityById(id int) (SpecialityCore, error) {

// }

// func (d *doctorBusiness) CreateSpeciality(speciality SpecialityCore) error {

// }

// func (d *doctorBusiness) EditSpeciality(speciality SpecialityCore) error {

// }

// func (d *doctorBusiness) RemoveSpeciality(id int) error {

// }

// func (d *doctorBusiness) FindRooms() ([]RoomCore, error) {

// }

// func (d *doctorBusiness) CreateRoom(room RoomCore) error {

// }

// func (d *doctorBusiness) EditRoom(room RoomCore) error {

// }

// func (d *doctorBusiness) RemoveRoomById(id int) error {

// }
