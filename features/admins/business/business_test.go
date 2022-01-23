package business_test

import (
	"os"
	"testing"
	"time"

	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
	"github.com/final-project-alterra/hospital-management-system-api/utils/files"
	"github.com/final-project-alterra/hospital-management-system-api/utils/hash"
	"github.com/final-project-alterra/hospital-management-system-api/utils/project"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	ab "github.com/final-project-alterra/hospital-management-system-api/features/admins/business"
	adminMock "github.com/final-project-alterra/hospital-management-system-api/features/admins/mocks"
	doctorMock "github.com/final-project-alterra/hospital-management-system-api/features/doctors/mocks"
	nurseMock "github.com/final-project-alterra/hospital-management-system-api/features/nurses/mocks"
)

var (
	adminsData adminMock.IData

	adminsBusiness  admins.IBusiness
	doctorsBusiness doctorMock.IBusiness
	nursesBusiness  nurseMock.IBusiness

	adminValue admins.AdminCore
	newAdmin   admins.AdminCore

	errNotFound error
	errServer   error
)

// bikin adapternya
func TestMain(m *testing.M) {
	adminsBusiness = ab.NewAdminBusinessBuilder().
		SetData(&adminsData).
		SetDoctorBusiness(&doctorsBusiness).
		SetNurseBusiness(&nursesBusiness).
		Build()

	errNotFound = errors.E(errors.New("not found"), errors.KindNotFound)
	errServer = errors.E(errors.New("error"), errors.KindServerError)

	adminValue = admins.AdminCore{
		ID:        1,
		CreatedBy: 0,
		UpdatedBy: 0,
		Email:     "riza@mail.com",
		Password:  "admin",
		Name:      "Riza",
		BirthDate: "1990-10-23",
		ImageUrl:  "google.com/img.jpg",
		Phone:     "08128288282828",
		Address:   "Jalan baru",
		Gender:    "L",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newAdmin = adminValue
	newAdmin.ID = 0
	newAdmin.CreatedBy = 1
	newAdmin.Email = "hernowo@mail.com"
	newAdmin.Name = "hernowo"

	files.Remove = func(path string) error { return nil }
	project.GetMainDir = func() string { return "" }

	os.Exit(m.Run())
}

func TestFindAdmins(t *testing.T) {
	t.Run("valid - when find admins success", func(t *testing.T) {
		adminsData.
			On("SelectAdmins").
			Return([]admins.AdminCore{adminValue}, nil).
			Once()

		result, err := adminsBusiness.FindAdmins()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})

	t.Run("valid - when SelectAdmins error", func(t *testing.T) {
		err := errors.E(errors.New("error"), errors.KindServerError)
		adminsData.
			On("SelectAdmins").
			Return([]admins.AdminCore{}, err).
			Once()

		result, err := adminsBusiness.FindAdmins()

		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
		assert.Equal(t, 0, len(result))
	})
}

func TestFindAdminById(t *testing.T) {
	t.Run("valid - when FindAdminById success", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, nil).
			Once()

		result, err := adminsBusiness.FindAdminById(1)

		assert.Nil(t, err)
		assert.Equal(t, adminValue, result)
	})

	t.Run("valid - when FindAdminById cannot find admin", func(t *testing.T) {
		err := errors.E(errors.New("not found"), errors.KindNotFound)
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, err).
			Once()

		result, err := adminsBusiness.FindAdminById(1)

		assert.Error(t, err)
		assert.Equal(t, errors.KindNotFound, errors.Kind(err))
		assert.Equal(t, result, admins.AdminCore{})
	})
}

func TestFindAdminByEmail(t *testing.T) {
	t.Run("valid - when FindAdminByEmail success", func(t *testing.T) {
		adminsData.
			On("SelectAdminByEmail", mock.AnythingOfType("string")).
			Return(adminValue, nil).
			Once()

		result, err := adminsBusiness.FindAdminByEmail("wow@mail.com")

		assert.Nil(t, err)
		assert.Equal(t, adminValue, result)
	})

	t.Run("valid - when FindAdminByEmail cannot find admin", func(t *testing.T) {
		err := errors.E(errors.New("not found"), errors.KindNotFound)
		adminsData.
			On("SelectAdminByEmail", mock.AnythingOfType("string")).
			Return(admins.AdminCore{}, err).
			Once()

		result, err := adminsBusiness.FindAdminByEmail("wow@mail.com")

		assert.Error(t, err)
		assert.Equal(t, errors.KindNotFound, errors.Kind(err))
		assert.Equal(t, admins.AdminCore{}, result)
	})
}

func TestCreateAdmin(t *testing.T) {
	t.Run("valid - when CreateAdmin success", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, nil).
			Once()

		adminsData.
			On("SelectAdminByEmail", mock.AnythingOfType("string")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		doctorsBusiness.
			On("FindDoctorByEmail", mock.AnythingOfType("string")).
			Return(doctors.DoctorCore{}, errNotFound).
			Once()

		nursesBusiness.
			On("FindNurseByEmail", mock.AnythingOfType("string")).
			Return(nurses.NurseCore{}, errNotFound).
			Once()

		adminsData.
			On("InsertAdmin", mock.AnythingOfType("admins.AdminCore")).
			Return(nil).
			Once()

		err := adminsBusiness.CreateAdmin(newAdmin)

		assert.Nil(t, err)
	})

	t.Run("valid - when admin who add new admin is not found", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		err := adminsBusiness.CreateAdmin(newAdmin)

		assert.Error(t, err)
		assert.Equal(t, errors.KindNotFound, errors.Kind(err))
	})

	t.Run("valid - when unkown SelectAdminById error occurs", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errServer).
			Once()

		err := adminsBusiness.CreateAdmin(newAdmin)

		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when email already registered", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			adminsData.
				On("SelectAdminById", mock.AnythingOfType("int")).
				Return(admins.AdminCore{}, nil).
				Once()

			switch i {
			case 0:
				adminsData.
					On("SelectAdminByEmail", mock.AnythingOfType("string")).
					Return(adminValue, nil).
					Once()
				doctorsBusiness.
					On("FindDoctorByEmail", mock.AnythingOfType("string")).
					Return(doctors.DoctorCore{}, errNotFound).
					Once()
				nursesBusiness.
					On("FindNurseByEmail", mock.AnythingOfType("string")).
					Return(nurses.NurseCore{}, errNotFound).
					Once()
			case 1:
				adminsData.
					On("SelectAdminByEmail", mock.AnythingOfType("string")).
					Return(admins.AdminCore{}, errNotFound).
					Once()
				doctorsBusiness.
					On("FindDoctorByEmail", mock.AnythingOfType("string")).
					Return(doctors.DoctorCore{ID: 1}, nil).
					Once()
				nursesBusiness.
					On("FindNurseByEmail", mock.AnythingOfType("string")).
					Return(nurses.NurseCore{}, errNotFound).
					Once()
			case 2:
				adminsData.
					On("SelectAdminByEmail", mock.AnythingOfType("string")).
					Return(admins.AdminCore{}, errNotFound).
					Once()
				doctorsBusiness.
					On("FindDoctorByEmail", mock.AnythingOfType("string")).
					Return(doctors.DoctorCore{}, errNotFound).
					Once()
				nursesBusiness.
					On("FindNurseByEmail", mock.AnythingOfType("string")).
					Return(nurses.NurseCore{ID: 1}, nil).
					Once()
			}

			err := adminsBusiness.CreateAdmin(newAdmin)

			assert.Error(t, err)
			assert.Equal(t, errors.KindUnprocessable, errors.Kind(err))
		}
	})

	t.Run("valid - when unknwon error occurs on InsertAdmin", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, nil).
			Once()

		adminsData.
			On("SelectAdminByEmail", mock.AnythingOfType("string")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		adminsData.
			On("InsertAdmin", mock.AnythingOfType("admins.AdminCore")).
			Return(errServer).
			Once()

		doctorsBusiness.
			On("FindDoctorByEmail", mock.AnythingOfType("string")).
			Return(doctors.DoctorCore{}, errNotFound).
			Once()

		nursesBusiness.
			On("FindNurseByEmail", mock.AnythingOfType("string")).
			Return(nurses.NurseCore{}, errNotFound).
			Once()

		err := adminsBusiness.CreateAdmin(newAdmin)

		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})
}

func TestEditAdmin(t *testing.T) {
	t.Run("valid - editing admin", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, nil).
			Twice()

		adminsData.
			On("UpdateAdmin", mock.AnythingOfType("admins.AdminCore")).
			Return(nil).
			Once()

		err := adminsBusiness.EditAdmin(adminValue)
		assert.Nil(t, err)
	})

	t.Run("valid - when existing admin not found", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		err := adminsBusiness.EditAdmin(adminValue)

		assert.Error(t, err)
		assert.Equal(t, errors.KindNotFound, errors.Kind(err))
	})

	t.Run("valid - when unknown error occurs on SelectAdminById", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errServer).
			Once()

		err := adminsBusiness.EditAdmin(adminValue)

		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when admin who updating data is not found", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, nil).
			Once()

		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		err := adminsBusiness.EditAdmin(adminValue)

		assert.Error(t, err)
		assert.Equal(t, errors.KindNotFound, errors.Kind(err))
	})

	t.Run("valid - when trying to find admin who wants to update data", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, nil).
			Once()

		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errServer).
			Once()

		err := adminsBusiness.EditAdmin(adminValue)

		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when fauled editing admin", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, nil).
			Twice()

		adminsData.
			On("UpdateAdmin", mock.AnythingOfType("admins.AdminCore")).
			Return(errServer).
			Once()

		err := adminsBusiness.EditAdmin(adminValue)
		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})
}

func TestEditAdminPassword(t *testing.T) {
	t.Run("valid - editing admin password", func(t *testing.T) {
		hashed, err := hash.Generate(adminValue.Password)
		if err != nil {
			panic(err)
		}
		adminValue.Password = hashed

		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, nil).
			Twice()

		adminsData.
			On("UpdateAdmin", mock.AnythingOfType("admins.AdminCore")).
			Return(nil).
			Once()

		err = adminsBusiness.EditAdminPassword(2, 1, "admin", "admin123")

		assert.Nil(t, err)
	})

	t.Run("valid - when existing admin not found", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		err := adminsBusiness.EditAdminPassword(2, 1, "old", "new")

		assert.Error(t, err)
		assert.Equal(t, errors.KindNotFound, errors.Kind(err))
	})

	t.Run("valid - when unknown error occurs on SelectAdminById", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errServer).
			Once()

		err := adminsBusiness.EditAdminPassword(2, 1, "old", "new")

		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when admin who updating data is not found", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, nil).
			Once()

		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		err := adminsBusiness.EditAdminPassword(2, 1, "old", "new")

		assert.Error(t, err)
		assert.Equal(t, errors.KindNotFound, errors.Kind(err))
	})

	t.Run("valid - when unknown error occurs on SelectAdminById", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, nil).
			Once()

		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errServer).
			Once()

		err := adminsBusiness.EditAdminPassword(2, 1, "old", "new")

		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when editing admin password failed", func(t *testing.T) {
		hashed, err := hash.Generate("admin")
		if err != nil {
			panic(err)
		}
		adminValue.Password = hashed

		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, nil).
			Twice()

		adminsData.
			On("UpdateAdmin", mock.AnythingOfType("admins.AdminCore")).
			Return(errServer).
			Once()

		err = adminsBusiness.EditAdminPassword(2, 1, "admin", "admin123")

		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when password doesnot match", func(t *testing.T) {
		copy := adminValue
		hashed, err := hash.Generate(adminValue.Password)
		if err != nil {
			panic(err)
		}
		copy.Password = hashed

		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(copy, nil).
			Twice()

		adminsData.
			On("UpdateAdmin", mock.AnythingOfType("admins.AdminCore")).
			Return(nil).
			Once()

		err = adminsBusiness.EditAdminPassword(2, 1, "admin-salah", "admin123")

		assert.Error(t, err)
		assert.Equal(t, errors.KindUnprocessable, errors.Kind(err))
	})
}

func TestEditAdminProfileImage(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, nil).
			Twice()

		adminsData.
			On("UpdateAdmin", mock.Anything).
			Return(nil).
			Once()

		err := adminsBusiness.EditAdminProfileImage(admins.AdminCore{})
		assert.Nil(t, err)
	})

	t.Run("valid - when SelectAdminById error", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errServer).
			Once()

		for i := 0; i < 2; i++ {
			err := adminsBusiness.EditAdminProfileImage(admins.AdminCore{})
			assert.Error(t, err)
		}
	})

	t.Run("valid - when second SelectAdminById error", func(t *testing.T) {
		for i := 1; i <= 2; i++ {
			adminsData.
				On("SelectAdminById", mock.AnythingOfType("int")).
				Return(admins.AdminCore{}, nil).
				Once()
			if i%2 != 0 {
				adminsData.
					On("SelectAdminById", mock.AnythingOfType("int")).
					Return(admins.AdminCore{}, errNotFound).
					Once()
			} else {
				adminsData.
					On("SelectAdminById", mock.AnythingOfType("int")).
					Return(admins.AdminCore{}, errServer).
					Once()
			}
			err := adminsBusiness.EditAdminProfileImage(admins.AdminCore{})
			assert.Error(t, err)
		}
	})

}

func TestRemoveAdminById(t *testing.T) {
	t.Run("valid - when RemoveAdminById success", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, nil).
			Once()

		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, nil).
			Once()

		adminsData.
			On("DeleteAdminById", mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(nil).
			Once()

		err := adminsBusiness.RemoveAdminById(2, 1)
		assert.Nil(t, err)
	})

	t.Run("valid - when admin who update data not found", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, errNotFound).
			Once()

		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, errServer).
			Once()

		for i := 0; i < 2; i++ {
			err := adminsBusiness.RemoveAdminById(2, 1)
			assert.Error(t, err)
		}
	})

	t.Run("valid - when admin who wants to be updated not found", func(t *testing.T) {

		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, nil).
			Once()

		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, errNotFound).
			Once()

		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, nil).
			Once()

		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, errServer).
			Once()

		for i := 0; i < 2; i++ {
			err := adminsBusiness.RemoveAdminById(2, 1)
			assert.Error(t, err)
		}
	})

	t.Run("valid - when DeleteAdminById error", func(t *testing.T) {
		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, nil).
			Once()

		adminsData.
			On("SelectAdminById", mock.AnythingOfType("int")).
			Return(adminValue, nil).
			Once()

		adminsData.
			On("DeleteAdminById", mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(errServer).
			Once()

		err := adminsBusiness.RemoveAdminById(2, 1)
		assert.Error(t, err)
	})
}
