package business_test

import (
	"os"
	"testing"

	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	am "github.com/final-project-alterra/hospital-management-system-api/features/admins/mocks"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
	dm "github.com/final-project-alterra/hospital-management-system-api/features/doctors/mocks"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
	nb "github.com/final-project-alterra/hospital-management-system-api/features/nurses/business"
	nm "github.com/final-project-alterra/hospital-management-system-api/features/nurses/mocks"
	sm "github.com/final-project-alterra/hospital-management-system-api/features/schedules/mocks"
	"github.com/final-project-alterra/hospital-management-system-api/utils/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	repo nm.IData

	business         nurses.IBusiness
	adminBusiness    am.IBusiness
	doctorBusiness   dm.IBusiness
	scheduleBusiness sm.IBusiness

	admin1 admins.AdminCore
	nurse1 nurses.NurseCore

	errServer   error
	errNotFound error
)

func TestMain(t *testing.M) {
	business = nb.NewNurseBusinessBuilder().
		SetData(&repo).
		SetAdminBusiness(&adminBusiness).
		SetDoctorBusiness(&doctorBusiness).
		SetScheduleBusiness(&scheduleBusiness).
		Build()

	admin1 = admins.AdminCore{
		ID:   1,
		Name: "admin1",
	}

	nurse1 = nurses.NurseCore{
		ID:        1,
		CreatedBy: 1,
		UpdatedBy: 1,
		Name:      "Nurse 1",
		Email:     "example@mail.com",
		Password:  "password",
	}

	errNotFound = errors.E(errors.New("not found"), errors.KindNotFound)
	errServer = errors.E(errors.New("server error"), errors.KindServerError)

	os.Exit(t.Run())
}

func TestFindNurses(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectNurses").
			Return([]nurses.NurseCore{nurse1}, nil).
			Once()

		result, err := business.FindNurses()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})

	t.Run("valid - when SelectNurses return error", func(t *testing.T) {
		repo.
			On("SelectNurses").
			Return([]nurses.NurseCore{}, errServer).
			Once()

		result, err := business.FindNurses()

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})
}

func TestFindNursesByIds(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectNursesByIds", []int{1}).
			Return([]nurses.NurseCore{nurse1}, nil).
			Once()

		result, err := business.FindNursesByIds([]int{1})

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})

	t.Run("valid - when SelectNursesByIds return error", func(t *testing.T) {
		repo.
			On("SelectNursesByIds", []int{1}).
			Return([]nurses.NurseCore{}, errServer).
			Once()

		result, err := business.FindNursesByIds([]int{1})

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})
}

func TestFindNurseById(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectNurseById", 1).
			Return(nurse1, nil).
			Once()

		result, err := business.FindNurseById(1)

		assert.Nil(t, err)
		assert.Equal(t, nurse1, result)
	})

	t.Run("valid - SelectNurseById return error", func(t *testing.T) {
		repo.
			On("SelectNurseById", 1).
			Return(nurses.NurseCore{}, errNotFound).
			Once()

		result, err := business.FindNurseById(1)

		assert.Error(t, err)
		assert.Equal(t, nurses.NurseCore{}, result)
	})
}

func TestFindNurseByEmail(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectNurseByEmail", mock.AnythingOfType("string")).
			Return(nurse1, nil).
			Once()

		result, err := business.FindNurseByEmail("hello@mail.com")

		assert.Nil(t, err)
		assert.Equal(t, nurse1, result)
	})

	t.Run("valid - SelectNurseByEmail return error", func(t *testing.T) {
		repo.
			On("SelectNurseByEmail", mock.AnythingOfType("string")).
			Return(nurses.NurseCore{}, errNotFound).
			Once()

		result, err := business.FindNurseByEmail("e@mail.com")

		assert.Error(t, err)
		assert.Equal(t, nurses.NurseCore{}, result)
	})
}

func TestCreateNurse(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin1, nil).
			Once()

		repo.
			On("SelectNurseByEmail", mock.AnythingOfType("string")).
			Return(nurses.NurseCore{}, errNotFound).
			Once()

		adminBusiness.
			On("FindAdminByEmail", mock.AnythingOfType("string")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		doctorBusiness.
			On("FindDoctorByEmail", mock.AnythingOfType("string")).
			Return(doctors.DoctorCore{}, errNotFound).
			Once()

		repo.
			On("InsertNurse", mock.AnythingOfType("nurses.NurseCore")).
			Return(nil).
			Once()

		err := business.CreateNurse(nurse1)
		assert.Nil(t, err)
	})

	t.Run("valid - when FindAdminById return not found", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		err := business.CreateNurse(nurse1)

		assert.Error(t, err)
	})

	t.Run("valid - when FindAdminById return server error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errServer).
			Once()

		err := business.CreateNurse(nurse1)

		assert.Error(t, err)
	})

	t.Run("valid - when email already exists", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin1, nil).
			Times(3)

		for i := 0; i < 3; i++ {
			switch i {
			case 0:
				repo.
					On("SelectNurseByEmail", mock.AnythingOfType("string")).
					Return(nurses.NurseCore{ID: 2, Email: nurse1.Email}, nil).
					Once()

				adminBusiness.
					On("FindAdminByEmail", mock.AnythingOfType("string")).
					Return(admins.AdminCore{}, errNotFound).
					Once()

				doctorBusiness.
					On("FindDoctorByEmail", mock.AnythingOfType("string")).
					Return(doctors.DoctorCore{}, errNotFound).
					Once()
			case 1:
				repo.
					On("SelectNurseByEmail", mock.AnythingOfType("string")).
					Return(nurses.NurseCore{}, errNotFound).
					Once()

				adminBusiness.
					On("FindAdminByEmail", mock.AnythingOfType("string")).
					Return(admins.AdminCore{ID: 1}, nil).
					Once()

				doctorBusiness.
					On("FindDoctorByEmail", mock.AnythingOfType("string")).
					Return(doctors.DoctorCore{}, errNotFound).
					Once()
			case 2:
				repo.
					On("SelectNurseByEmail", mock.AnythingOfType("string")).
					Return(nurses.NurseCore{}, errNotFound).
					Once()

				adminBusiness.
					On("FindAdminByEmail", mock.AnythingOfType("string")).
					Return(admins.AdminCore{}, errNotFound).
					Once()

				doctorBusiness.
					On("FindDoctorByEmail", mock.AnythingOfType("string")).
					Return(doctors.DoctorCore{ID: 1}, nil).
					Once()
			}

			err := business.CreateNurse(nurse1)
			assert.Error(t, err)
		}

	})

	t.Run("valid - when InsertNurse return error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin1, nil).
			Once()

		repo.
			On("SelectNurseByEmail", mock.AnythingOfType("string")).
			Return(nurses.NurseCore{}, errNotFound).
			Once()

		adminBusiness.
			On("FindAdminByEmail", mock.AnythingOfType("string")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		doctorBusiness.
			On("FindDoctorByEmail", mock.AnythingOfType("string")).
			Return(doctors.DoctorCore{}, errNotFound).
			Once()

		repo.
			On("InsertNurse", mock.AnythingOfType("nurses.NurseCore")).
			Return(errServer).
			Once()

		err := business.CreateNurse(nurse1)
		assert.Error(t, err)
	})

}

func TestEditNurse(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin1, nil).
			Once()

		repo.
			On("SelectNurseById", mock.AnythingOfType("int")).
			Return(nurse1, nil).
			Once()

		repo.
			On("UpdateNurse", mock.AnythingOfType("nurses.NurseCore")).
			Return(nil).
			Once()

		err := business.EditNurse(nurse1)
		assert.Nil(t, err)
	})

	t.Run("valid - when FindAdminById return error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errServer).
			Once()

		err := business.EditNurse(nurse1)
		assert.Error(t, err)
	})

	t.Run("valid - when SelectNurseById return error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin1, nil).
			Once()

		repo.
			On("SelectNurseById", mock.AnythingOfType("int")).
			Return(nurses.NurseCore{}, errServer).
			Once()

		err := business.EditNurse(nurse1)
		assert.Error(t, err)
	})

	t.Run("valid - when UpdateNurse return error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin1, nil).
			Once()

		repo.
			On("SelectNurseById", mock.AnythingOfType("int")).
			Return(nurse1, nil).
			Once()

		repo.
			On("UpdateNurse", mock.AnythingOfType("nurses.NurseCore")).
			Return(errServer).
			Once()

		err := business.EditNurse(nurse1)
		assert.Error(t, err)
	})
}

func TestEditNursePassword(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin1, nil).
			Once()

		record := nurse1
		hashed, err := hash.Generate(nurse1.Password)
		if err != nil {
			panic(err)
		}
		record.Password = hashed

		repo.
			On("SelectNurseById", mock.AnythingOfType("int")).
			Return(record, nil).
			Once()

		repo.
			On("UpdateNurse", mock.AnythingOfType("nurses.NurseCore")).
			Return(nil).
			Once()

		err = business.EditNursePassword(1, 2, nurse1.Password, "new password")
		assert.Nil(t, err)
	})

	t.Run("valid - when FindAdminById return error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errServer).
			Once()

		err := business.EditNursePassword(1, 2, nurse1.Password, "new password")
		assert.Error(t, err)
	})

	t.Run("valid - when SelectNurseById return error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin1, nil).
			Once()

		repo.
			On("SelectNurseById", mock.AnythingOfType("int")).
			Return(nurses.NurseCore{}, errNotFound).
			Once()

		err := business.EditNursePassword(1, 2, nurse1.Password, "new password")
		assert.Error(t, err)
	})

	t.Run("valid - when old password does not match", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin1, nil).
			Once()

		record := nurse1
		hashed, err := hash.Generate(nurse1.Password)
		if err != nil {
			panic(err)
		}
		record.Password = hashed

		repo.
			On("SelectNurseById", mock.AnythingOfType("int")).
			Return(record, nil).
			Once()

		err = business.EditNursePassword(1, 2, "wrong old password", "new password")
		assert.Error(t, err)
	})

	t.Run("valid - when UpdateNurse return error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin1, nil).
			Once()

		record := nurse1
		hashed, err := hash.Generate(nurse1.Password)
		if err != nil {
			panic(err)
		}
		record.Password = hashed

		repo.
			On("SelectNurseById", mock.AnythingOfType("int")).
			Return(record, nil).
			Once()

		repo.
			On("UpdateNurse", mock.AnythingOfType("nurses.NurseCore")).
			Return(errServer).
			Once()

		err = business.EditNursePassword(1, 2, nurse1.Password, "new password")
		assert.Error(t, err)
	})
}

func TestRemoveNurseById(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin1, nil).
			Once()

		scheduleBusiness.
			On("RemoveNurseFromNextWorkSchedules", mock.AnythingOfType("int")).
			Return(nil).
			Once()

		repo.
			On("DeleteNurseById", mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(nil).
			Once()

		err := business.RemoveNurseById(1, 2)
		assert.Nil(t, err)
	})

	t.Run("valid - when FindAdminById return error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errServer).
			Once()

		err := business.RemoveNurseById(1, 2)
		assert.Error(t, err)
	})

	t.Run("valid - when RemoveNurseFromNextWorkSchedules error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin1, nil).
			Once()

		scheduleBusiness.
			On("RemoveNurseFromNextWorkSchedules", mock.AnythingOfType("int")).
			Return(errServer).
			Once()

		err := business.RemoveNurseById(1, 2)
		assert.Error(t, err)
	})

	t.Run("valid - when DeleteNurseById return error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin1, nil).
			Once()

		scheduleBusiness.
			On("RemoveNurseFromNextWorkSchedules", mock.AnythingOfType("int")).
			Return(nil).
			Once()

		repo.
			On("DeleteNurseById", mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(errServer).
			Once()

		err := business.RemoveNurseById(1, 2)
		assert.Error(t, err)
	})

}
