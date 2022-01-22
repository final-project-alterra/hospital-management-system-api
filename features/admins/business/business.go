package business

import (
	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
	"github.com/final-project-alterra/hospital-management-system-api/utils/hash"
)

type adminBusiness struct {
	data           admins.IData
	doctorBusiness doctors.IBusiness
	nurseBusiness  nurses.IBusiness
}

func (ab *adminBusiness) FindAdmins() ([]admins.AdminCore, error) {
	const op errors.Op = "admins.business.FindAdmins"

	dataAdmins, err := ab.data.SelectAdmins()
	if err != nil {
		return []admins.AdminCore{}, errors.E(err, op)
	}
	return dataAdmins, nil
}

func (ab *adminBusiness) FindAdminById(id int) (admins.AdminCore, error) {
	const op errors.Op = "admins.business.FindAdminById"

	dataAdmin, err := ab.data.SelectAdminById(id)
	if err != nil {
		return admins.AdminCore{}, errors.E(err, op)
	}
	return dataAdmin, nil
}

func (ab *adminBusiness) FindAdminByEmail(email string) (admins.AdminCore, error) {
	const op errors.Op = "admins.business.FindAdminByEmail"

	dataAdmin, err := ab.data.SelectAdminByEmail(email)
	if err != nil {
		return admins.AdminCore{}, errors.E(err, op)
	}
	return dataAdmin, nil
}

func (ab *adminBusiness) CreateAdmin(admin admins.AdminCore) error {
	const op errors.Op = "admins.business.CreateAdmin"
	var errMessage errors.ErrClientMessage

	// Check admin who is creating this new admin, if not found or error, return error
	_, err := ab.data.SelectAdminById(admin.CreatedBy)
	if err != nil {
		switch errors.Kind(err) {
		case errors.KindNotFound:
			errMessage = "Admin who wants to update is not found"
			return errors.E(err, op, errMessage)
		default:
			return errors.E(err, op)
		}
	}

	// Check wheter email is already registered
	if err = ab.checkEmail(admin.Email); err != nil {
		return errors.E(err, op)
	}

	admin.Password, err = hash.Generate(admin.Password)
	if err != nil {
		errMessage = "Something went wrong"
		return errors.E(err, op, errMessage, errors.KindServerError)
	}

	err = ab.data.InsertAdmin(admin)
	if err != nil {
		return errors.E(err, op)
	}
	return nil
}

func (ab *adminBusiness) EditAdmin(admin admins.AdminCore) error {
	const op errors.Op = "admins.business.EditAdmin"
	var errMessage errors.ErrClientMessage

	existingAdmin, err := ab.data.SelectAdminById(admin.ID)
	if err != nil {
		switch errors.Kind(err) {
		case errors.KindNotFound:
			errMessage = "Admin who wants to be updated is not found"
			return errors.E(err, op, errMessage)
		default:
			return errors.E(err, op)
		}
	}

	_, err = ab.data.SelectAdminById(admin.UpdatedBy)
	if err != nil {
		switch errors.Kind(err) {
		case errors.KindNotFound:
			errMessage = "Admin who wants to update is not found"
			return errors.E(err, op, errMessage)
		default:
			return errors.E(err, op)
		}
	}

	existingAdmin.UpdatedBy = admin.UpdatedBy
	existingAdmin.Name = admin.Name
	existingAdmin.BirthDate = admin.BirthDate
	existingAdmin.Address = admin.Address
	existingAdmin.Phone = admin.Phone
	existingAdmin.Gender = admin.Gender

	err = ab.data.UpdateAdmin(existingAdmin)
	if err != nil {
		return errors.E(err, op)
	}
	return nil
}

func (ab *adminBusiness) EditAdminPassword(id int, updatedBy int, oldPassword string, newPassword string) error {
	const op errors.Op = "admins.business.EditAdminPassword"
	var errMessage errors.ErrClientMessage

	_, err := ab.data.SelectAdminById(updatedBy)
	if err != nil {
		switch errors.Kind(err) {
		case errors.KindNotFound:
			errMessage = "Admin who wants to update is not found"
			return errors.E(err, op, errMessage)
		default:
			return errors.E(err, op)
		}
	}

	existingAdmin, err := ab.data.SelectAdminById(id)
	if err != nil {
		switch errors.Kind(err) {
		case errors.KindNotFound:
			errMessage = "Admin who wants to be updated is not found"
			return errors.E(err, op, errMessage)
		default:
			return errors.E(err, op)
		}
	}

	doesMatch := hash.Validate(existingAdmin.Password, oldPassword)
	if !doesMatch {
		err = errors.New("Wrong password")
		errMessage = "Wrong old password"
		return errors.E(err, op, errMessage, errors.KindUnprocessable)
	}

	existingAdmin.UpdatedBy = updatedBy
	existingAdmin.Password, err = hash.Generate(newPassword)
	if err != nil {
		return errors.E(err, op)
	}

	err = ab.data.UpdateAdmin(existingAdmin)
	if err != nil {
		return errors.E(err, op)
	}

	return nil
}

func (ab *adminBusiness) RemoveAdminById(id int, updatedBy int) error {
	const op errors.Op = "admins.business.RemoveAdminById"
	var errMessage errors.ErrClientMessage

	_, err := ab.data.SelectAdminById(updatedBy)
	if err != nil {
		switch errors.Kind(err) {
		case errors.KindNotFound:
			errMessage = "Admin who wants to update is not found"
			return errors.E(err, op, errMessage)
		default:
			return errors.E(err, op)
		}
	}

	err = ab.data.DeleteAdminById(id, updatedBy)
	if err != nil {
		return errors.E(err, op)
	}
	return nil
}

// Private methods
func (ab *adminBusiness) checkEmail(email string) error {
	const op errors.Op = "admins.business.checkEmail"
	var errMsg errors.ErrClientMessage = "Email already exist"

	adminCh := make(chan error)
	doctorCh := make(chan error)
	nurseCh := make(chan error)

	go func() {
		_, err := ab.data.SelectAdminByEmail(email)
		if err != nil {
			if errors.Kind(err) == errors.KindNotFound {
				err = nil
			}
			adminCh <- err
			return
		}
		adminCh <- errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindUnprocessable)
	}()

	go func() {
		_, err := ab.doctorBusiness.FindDoctorByEmail(email)
		if err != nil {
			if errors.Kind(err) == errors.KindNotFound {
				err = nil
			}
			doctorCh <- err
			return
		}
		doctorCh <- errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindUnprocessable)
	}()

	go func() {
		_, err := ab.nurseBusiness.FindNurseByEmail(email)
		if err != nil {
			if errors.Kind(err) == errors.KindNotFound {
				err = nil
			}
			nurseCh <- err
			return
		}
		nurseCh <- errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindUnprocessable)
	}()

	adminErr := <-adminCh
	doctorErr := <-doctorCh
	nurseErr := <-nurseCh

	if adminErr == nil && doctorErr == nil && nurseErr == nil {
		return nil
	}

	if adminErr != nil {
		return adminErr
	}
	if doctorErr != nil {
		return doctorErr
	}
	return nurseErr
}
