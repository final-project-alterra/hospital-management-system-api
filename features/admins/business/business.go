package business

import (
	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/final-project-alterra/hospital-management-system-api/utils/hash"
)

type adminBusiness struct {
	data admins.IData
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
	_, err = ab.data.SelectAdminByEmail(admin.Email)
	if err == nil {
		err = errors.New("Duplicate email when createing new admin")
		errMessage = "Email already exists"
		return errors.E(err, op, errMessage, errors.KindUnprocessable)
	}
	if errors.Kind(err) != errors.KindNotFound {
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
			errMessage = "Admin who wants to update is not found"
			return errors.E(err, op, errMessage)
		default:
			return errors.E(err, op)
		}
	}

	_, err = ab.data.SelectAdminById(admin.UpdatedBy)
	if err != nil {
		switch errors.Kind(err) {
		case errors.KindNotFound:
			errMessage = "Admin who wants to be updated is not found"
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
