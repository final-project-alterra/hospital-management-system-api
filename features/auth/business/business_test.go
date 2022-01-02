package business_test

import (
	"os"
	"testing"

	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	amock "github.com/final-project-alterra/hospital-management-system-api/features/admins/mocks"
	"github.com/final-project-alterra/hospital-management-system-api/features/auth"
	authBusiness "github.com/final-project-alterra/hospital-management-system-api/features/auth/business"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
	"github.com/final-project-alterra/hospital-management-system-api/utils/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	dmock "github.com/final-project-alterra/hospital-management-system-api/features/doctors/mocks"

	nmock "github.com/final-project-alterra/hospital-management-system-api/features/nurses/mocks"
)

var (
	business auth.IBusiness

	adminBusiness  amock.IBusiness
	doctorBusiness dmock.IBusiness
	nurseBusiness  nmock.IBusiness

	admin  admins.AdminCore
	doctor doctors.DoctorCore
	nurse  nurses.NurseCore

	errServer   error
	errNotFound error
)

func TestMain(m *testing.M) {
	business = authBusiness.NewAuthBusinessBuilder().
		SetAdminBusiness(&adminBusiness).
		SetDoctorBusiness(&doctorBusiness).
		SetNurseBusiness(&nurseBusiness).
		Build()

	password, err := hash.Generate("12345678")
	if err != nil {
		panic(err)
	}

	admin = admins.AdminCore{
		ID:       1,
		Email:    "admin@mail.com",
		Password: password,
	}

	doctor = doctors.DoctorCore{
		ID:       1,
		Email:    "doctor@mail.com",
		Password: password,
	}

	nurse = nurses.NurseCore{
		ID:       1,
		Email:    "nurse@mail.com",
		Password: password,
	}

	errNotFound = errors.E(errors.New("not found"), errors.KindNotFound)
	errServer = errors.E(errors.New("error"), errors.KindServerError)

	os.Exit(m.Run())
}

func TestLogin(t *testing.T) {
	t.Run("valid - when admin authentication success", func(t *testing.T) {
		adminBusiness.
			On("FindAdminByEmail", admin.Email).
			Return(admin, nil).
			Once()

		token, err := business.Login("admin@mail.com", "12345678")
		assert.Nil(t, err)
		assert.NotEqual(t, "", token)
	})

	t.Run("valid - when admin authentication failed", func(t *testing.T) {
		adminBusiness.
			On("FindAdminByEmail", admin.Email).
			Return(admin, nil).
			Once()

		token, err := business.Login("admin@mail.com", "wrong password")
		assert.Error(t, err)
		assert.Equal(t, "", token)
	})

	t.Run("valid - when FindAdminByEmail returne server error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminByEmail", admin.Email).
			Return(admins.AdminCore{}, errServer).
			Once()

		token, err := business.Login("admin@mail.com", "wrong password")
		assert.Error(t, err)
		assert.Equal(t, "", token)
	})

	t.Run("valid - when doctor authentication success", func(t *testing.T) {
		adminBusiness.
			On("FindAdminByEmail", mock.AnythingOfType("string")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		doctorBusiness.
			On("FindDoctorByEmail", doctor.Email).
			Return(doctor, nil).
			Once()

		token, err := business.Login("doctor@mail.com", "12345678")
		assert.Nil(t, err)
		assert.NotEqual(t, "", token)
	})

	t.Run("valid - when doctor authentication failed", func(t *testing.T) {
		adminBusiness.
			On("FindAdminByEmail", mock.AnythingOfType("string")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		doctorBusiness.
			On("FindDoctorByEmail", doctor.Email).
			Return(doctor, nil).
			Once()

		token, err := business.Login("doctor@mail.com", "wrong password")
		assert.Error(t, err)
		assert.Equal(t, "", token)
	})

	t.Run("valid - when FindDoctorByEmail return server error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminByEmail", mock.AnythingOfType("string")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		doctorBusiness.
			On("FindDoctorByEmail", mock.AnythingOfType("string")).
			Return(doctors.DoctorCore{}, errServer).
			Once()

		token, err := business.Login("doctor@mail.com", "wrong password")
		assert.Error(t, err)
		assert.Equal(t, "", token)
	})

	t.Run("valid - when nurse authentication success", func(t *testing.T) {
		adminBusiness.
			On("FindAdminByEmail", mock.AnythingOfType("string")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		doctorBusiness.
			On("FindDoctorByEmail", mock.AnythingOfType("string")).
			Return(doctors.DoctorCore{}, errNotFound).
			Once()

		nurseBusiness.
			On("FindNurseByEmail", nurse.Email).
			Return(nurse, nil).
			Once()

		token, err := business.Login("nurse@mail.com", "12345678")
		assert.Nil(t, err)
		assert.NotEqual(t, "", token)
	})

	t.Run("valid - when nurse authentication failed", func(t *testing.T) {
		adminBusiness.
			On("FindAdminByEmail", mock.AnythingOfType("string")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		doctorBusiness.
			On("FindDoctorByEmail", mock.AnythingOfType("string")).
			Return(doctors.DoctorCore{}, errNotFound).
			Once()

		nurseBusiness.
			On("FindNurseByEmail", nurse.Email).
			Return(nurse, nil).
			Once()

		token, err := business.Login("nurse@mail.com", "wrong password")
		assert.Error(t, err)
		assert.Equal(t, "", token)
	})

	t.Run("valid - when FindNurseByEmail return server error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminByEmail", mock.AnythingOfType("string")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		doctorBusiness.
			On("FindDoctorByEmail", mock.AnythingOfType("string")).
			Return(doctors.DoctorCore{}, errNotFound).
			Once()

		nurseBusiness.
			On("FindNurseByEmail", mock.AnythingOfType("string")).
			Return(nurses.NurseCore{}, errServer).
			Once()

		token, err := business.Login("nurse@mail.com", "wrong password")
		assert.Error(t, err)
		assert.Equal(t, "", token)
	})

	t.Run("valid - when no account found", func(t *testing.T) {
		adminBusiness.
			On("FindAdminByEmail", mock.AnythingOfType("string")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		doctorBusiness.
			On("FindDoctorByEmail", mock.AnythingOfType("string")).
			Return(doctors.DoctorCore{}, errNotFound).
			Once()

		nurseBusiness.
			On("FindNurseByEmail", mock.AnythingOfType("string")).
			Return(nurses.NurseCore{}, errNotFound).
			Once()

		token, err := business.Login("nurse@mail.com", "wrong password")
		assert.Error(t, err)
		assert.Equal(t, "", token)
	})

}
