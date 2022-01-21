package business

import (
	"os"
	"path"

	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
	"github.com/final-project-alterra/hospital-management-system-api/features/schedules"
	"github.com/final-project-alterra/hospital-management-system-api/utils/hash"
	"github.com/final-project-alterra/hospital-management-system-api/utils/project"
)

type nurseBusiness struct {
	data             nurses.IData
	adminBusiness    admins.IBusiness
	doctorBusiness   doctors.IBusiness
	scheduleBusiness schedules.IBusiness
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

	if err = n.checkEmail(nurse.Email); err != nil {
		return errors.E(err, op)
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

	_, err := n.adminBusiness.FindAdminById(nurse.UpdatedBy)
	if err != nil {
		return errors.E(err, op)
	}

	existingNurse, err := n.data.SelectNurseById(nurse.ID)
	if err != nil {
		return errors.E(op, err)
	}

	existingNurse.Name = nurse.Name
	existingNurse.BirthDate = nurse.BirthDate
	existingNurse.Phone = nurse.Phone
	existingNurse.Address = nurse.Address
	existingNurse.Gender = nurse.Gender

	err = n.data.UpdateNurse(existingNurse)
	if err != nil {
		return errors.E(op, err)
	}
	return nil
}

func (nb *nurseBusiness) EditNurseImageProfile(nurse nurses.NurseCore) error {
	const op errors.Op = "nurses.business.EditNurseImageProfile"

	newImage := path.Join(project.GetMainDir(), "files", nurse.ImageUrl)

	_, err := nb.adminBusiness.FindAdminById(nurse.UpdatedBy)
	if err != nil {
		go os.Remove(newImage)
		return errors.E(err, op)
	}

	existingNurse, err := nb.data.SelectNurseById(nurse.ID)
	if err != nil {
		go os.Remove(newImage)
		return errors.E(err, op)
	}
	oldImage := path.Join(project.GetMainDir(), "files", existingNurse.ImageUrl)

	existingNurse.ImageUrl = nurse.ImageUrl
	existingNurse.UpdatedBy = nurse.UpdatedBy

	err = nb.data.UpdateNurse(existingNurse)
	if err != nil {
		go os.Remove(newImage)
		return errors.E(err, op)
	}

	go os.Remove(oldImage)

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

	existingNurse, err := n.data.SelectNurseById(id)
	if err != nil {
		return errors.E(err, op)
	}
	existingImage := path.Join(project.GetMainDir(), "files", existingNurse.ImageUrl)

	err = n.scheduleBusiness.RemoveNurseFromNextWorkSchedules(id)
	if err != nil {
		return errors.E(err, op)
	}

	err = n.data.DeleteNurseById(id, updatedBy)
	if err != nil {
		return errors.E(err, op)
	}

	go os.Remove(existingImage)

	return nil
}

// Private methods
func (nb *nurseBusiness) checkEmail(email string) error {
	const op errors.Op = "admins.business.checkEmail"
	var errMsg errors.ErrClientMessage = "Email already exist"

	adminCh := make(chan error)
	doctorCh := make(chan error)
	nurseCh := make(chan error)

	go func() {
		_, err := nb.data.SelectNurseByEmail(email)
		if err != nil {
			if errors.Kind(err) == errors.KindNotFound {
				err = nil
			}
			nurseCh <- err
			return
		}
		nurseCh <- errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindUnprocessable)
	}()

	go func() {
		_, err := nb.adminBusiness.FindAdminByEmail(email)
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
		_, err := nb.doctorBusiness.FindDoctorByEmail(email)
		if err != nil {
			if errors.Kind(err) == errors.KindNotFound {
				err = nil
			}
			doctorCh <- err
			return
		}
		doctorCh <- errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindUnprocessable)
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
