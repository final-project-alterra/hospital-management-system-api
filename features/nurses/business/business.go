package business

import (
	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
	"github.com/final-project-alterra/hospital-management-system-api/utils/hash"
)

type nurseBusiness struct {
	data          nurses.IData
	adminBusiness admins.IBusiness
}

func (n *nurseBusiness) FindNurses() ([]nurses.NurseCore, error) {
	const op errors.Op = "nurses.business.FindNurses"

	nursesData, err := n.data.SelectNurses()
	if err != nil {
		return []nurses.NurseCore{}, errors.E(op, err)
	}
	return nursesData, nil
}

func (n *nurseBusiness) FindNursesByIds(ids []int) ([]nurses.NurseCore, error) {
	const op errors.Op = "nurses.business.FindNursesByIds"

	nursesData, err := n.data.SelectNursesByIds(ids)
	if err != nil {
		return []nurses.NurseCore{}, errors.E(op, err)
	}
	return nursesData, nil
}

func (n *nurseBusiness) FindNurseById(id int) (nurses.NurseCore, error) {
	const op errors.Op = "nurses.business.FindNurseById"

	nurseData, err := n.data.SelectNurseById(id)
	if err != nil {
		return nurses.NurseCore{}, errors.E(op, err)
	}
	return nurseData, nil
}

func (n *nurseBusiness) FindNurseByEmail(email string) (nurses.NurseCore, error) {
	const op errors.Op = "nurses.business.FindNurseByEmail"

	nurseData, err := n.data.SelectNurseByEmail(email)
	if err != nil {
		return nurses.NurseCore{}, errors.E(op, err)
	}
	return nurseData, nil
}

func (n *nurseBusiness) CreateNurse(nurse nurses.NurseCore) error {
	const op errors.Op = "nurses.business.CreateNurse"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	_, err := n.adminBusiness.FindAdminById(nurse.CreatedBy)
	if err != nil {
		switch errors.Kind(err) {
		case errors.KindNotFound:
			errMessage = "Admin who wants to create this nurse is not found"
			return errors.E(err, op, errMessage, errors.KindNotFound)

		default:
			return errors.E(err, op, errMessage, errors.KindServerError)
		}
	}

	_, err = n.data.SelectNurseByEmail(nurse.Email)
	if err == nil {
		err = errors.New("Email is already used")
		errMessage = "Email is already used"
		return errors.E(err, op, errMessage, errors.KindUnprocessable)
	}
	if errors.Kind(err) != errors.KindNotFound {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}

	nurse.Password, err = hash.Generate(nurse.Password)
	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}

	err = n.data.InsertNurse(nurse)
	if err != nil {
		return errors.E(op, err)
	}
	return nil
}

func (n *nurseBusiness) EditNurse(nurse nurses.NurseCore) error {
	const op errors.Op = "nurses.business.EditNurse"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	_, err := n.adminBusiness.FindAdminById(nurse.UpdatedBy)
	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}

	existingNurse, err := n.data.SelectNurseById(nurse.ID)
	if err != nil {
		return errors.E(op, err)
	}

	existingNurse.Name = nurse.Name
	existingNurse.Age = nurse.Age
	existingNurse.Phone = nurse.Phone
	existingNurse.Address = nurse.Address
	existingNurse.Gender = nurse.Gender

	err = n.data.UpdateNurse(existingNurse)
	if err != nil {
		return errors.E(op, err)
	}
	return nil
}

func (n *nurseBusiness) EditNursePassword(id int, updatedBy int, oldPassword string, newPassword string) error {
	const op errors.Op = "nurses.business.EditNursePassword"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	_, err := n.adminBusiness.FindAdminById(updatedBy)
	if err != nil {
		return errors.E(err, op)
	}

	existingNurse, err := n.data.SelectNurseById(id)
	if err != nil {
		return errors.E(op, err)
	}

	doesMatch := hash.Validate(existingNurse.Password, oldPassword)
	if !doesMatch {
		err = errors.New("Wrong old password")
		errMessage = "Wrong old password"

		return errors.E(err, op, errMessage, errors.KindUnprocessable)
	}

	existingNurse.Password, err = hash.Generate(newPassword)
	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}

	err = n.data.UpdateNurse(existingNurse)
	if err != nil {
		return errors.E(err, op)
	}
	return nil
}

func (n *nurseBusiness) RemoveNurseById(id int, updatedBy int) error {
	const op errors.Op = "nurses.business.RemoveNurseById"

	_, err := n.adminBusiness.FindAdminById(updatedBy)
	if err != nil {
		return errors.E(err, op)
	}

	err = n.data.DeleteNurseById(id, updatedBy)
	if err != nil {
		return errors.E(err, op)
	}
	return nil
}
