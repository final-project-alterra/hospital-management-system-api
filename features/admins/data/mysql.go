package data

import (
	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"gorm.io/gorm"
)

type MySQLRepo struct {
	db *gorm.DB
}

func NewMySQLRepo(db *gorm.DB) *MySQLRepo {
	return &MySQLRepo{
		db: db,
	}
}

func (r *MySQLRepo) SelectAdmins() ([]admins.AdminCore, error) {
	const op errors.Op = "admins.data.SelectAdmins"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	data := []Admin{}
	err := r.db.Find(&data).Error
	if err != nil {
		return ToSliceAdminCore(data), errors.E(err, op, errMessage, errors.KindServerError)
	}
	return ToSliceAdminCore(data), nil
}

func (r *MySQLRepo) SelectAdminById(id int) (admins.AdminCore, error) {
	const op errors.Op = "admins.data.SelectAdminById"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	data := Admin{}
	err := r.db.First(&data, id).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			errMessage = "Admin not found"
			return data.ToAdminCore(), errors.E(err, op, errMessage, errors.KindNotFound)
		default:
			return data.ToAdminCore(), errors.E(err, op, errMessage, errors.KindServerError)
		}
	}
	return data.ToAdminCore(), nil
}

func (r *MySQLRepo) SelectAdminByEmail(email string) (admins.AdminCore, error) {
	const op errors.Op = "admins.data.SelectAdminByEmail"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	data := Admin{}
	err := r.db.Where("email = ?", email).First(&data).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			errMessage = "Admin not found"
			return data.ToAdminCore(), errors.E(err, op, errMessage, errors.KindNotFound)
		default:
			return data.ToAdminCore(), errors.E(err, op, errMessage, errors.KindServerError)
		}
	}
	return data.ToAdminCore(), nil
}

func (r *MySQLRepo) InsertAdmin(admin admins.AdminCore) error {
	const op errors.Op = "admins.data.InsertAdmin"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	createdBy := uint(admin.CreatedBy)

	data := Admin{
		CreatedBy: &createdBy,
		Email:     admin.Email,
		Password:  admin.Password,
		Name:      admin.Name,
		Age:       admin.Age,
		ImageUrl:  admin.ImageUrl,
		Phone:     admin.Phone,
		Address:   admin.Address,
		Gender:    admin.Gender,
	}

	err := r.db.Create(&data).Error
	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}
	return nil
}

func (r *MySQLRepo) UpdateAdmin(admin admins.AdminCore) error {
	const op errors.Op = "admins.data.UpdateAdmin"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	createdBy := uint(admin.CreatedBy)
	updatedBy := uint(admin.UpdatedBy)

	data := Admin{
		Model: gorm.Model{
			ID:        uint(admin.ID),
			CreatedAt: admin.CreatedAt,
		},
		CreatedBy: &createdBy,
		UpdatedBy: &updatedBy,
		Email:     admin.Email,
		Password:  admin.Password,
		Name:      admin.Name,
		Age:       admin.Age,
		ImageUrl:  admin.ImageUrl,
		Phone:     admin.Phone,
		Address:   admin.Address,
		Gender:    admin.Gender,
	}

	err := r.db.Save(&data).Error
	if err != nil {
		return errors.E(err, op, errMessage, errors.KindServerError)
	}
	return nil
}

func (r *MySQLRepo) DeleteAdminById(id int, updatedBy int) error {
	const op errors.Op = "admins.data.DeleteAdminById"
	var errMessage errors.ErrClientMessage = "Something went wrong"

	deleteTransaction := func(tx *gorm.DB) error {
		data := Admin{}
		err := tx.First(&data, id).Error

		// data not found, so no need to delete
		if err != nil && err == gorm.ErrRecordNotFound {
			return nil
		}

		if err != nil {
			return errors.E(err, op, errMessage, errors.KindServerError)
		}

		updatedByConverted := uint(updatedBy)
		data.UpdatedBy = &updatedByConverted
		err = tx.Save(&data).Error
		if err != nil {
			return errors.E(err, op, errMessage, errors.KindServerError)
		}

		err = tx.Delete(&Admin{}, id).Error
		if err != nil {
			return errors.E(err, op, errMessage, errors.KindServerError)
		}
		return nil
	}

	err := r.db.Transaction(deleteTransaction)
	if err != nil {
		return err
	}
	return nil
}
