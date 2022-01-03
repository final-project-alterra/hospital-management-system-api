package business_test

import (
	"os"
	"testing"

	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	amocks "github.com/final-project-alterra/hospital-management-system-api/features/admins/mocks"
	"github.com/final-project-alterra/hospital-management-system-api/features/patients"
	pb "github.com/final-project-alterra/hospital-management-system-api/features/patients/business"
	pmocks "github.com/final-project-alterra/hospital-management-system-api/features/patients/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	repo pmocks.IData

	adminBusiness amocks.IBusiness
	business      patients.IBusiness

	patient patients.PatientCore
	admin   admins.AdminCore

	errNotFound error
	errServer   error
)

func TestMain(m *testing.M) {
	business = pb.NewPatientBusinessBuilder().SetData(&repo).SetAdminBusiness(&adminBusiness).Build()

	patient = patients.PatientCore{
		ID:   1,
		NIK:  "123456789",
		Name: "John Doe",
	}
	admin = admins.AdminCore{
		ID:    1,
		Email: "admin@mail.com",
		Name:  "admin",
	}

	errNotFound = errors.E(errors.New("not found"), errors.KindNotFound)
	errServer = errors.E(errors.New("server error"), errors.KindServerError)

	os.Exit(m.Run())
}

func TestFindPatients(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectPatients").
			Return([]patients.PatientCore{patient}, nil).
			Once()

		result, err := business.FindPatients()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})

	t.Run("valid - when SelectPatients return error", func(t *testing.T) {
		repo.
			On("SelectPatients").
			Return([]patients.PatientCore{}, errServer).
			Once()

		result, err := business.FindPatients()

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})
}

func TestFindPatientsByIds(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectPatientsByIds", mock.AnythingOfType("[]int")).
			Return([]patients.PatientCore{patient}, nil).
			Once()

		result, err := business.FindPatientsByIds([]int{1})

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})

	t.Run("valid - when SelectPatientsByIds return error", func(t *testing.T) {
		repo.
			On("SelectPatientsByIds", mock.AnythingOfType("[]int")).
			Return([]patients.PatientCore{}, errServer).
			Once()

		result, err := business.FindPatientsByIds([]int{1})

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})
}

func TestFindPatientById(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectPatientById", mock.AnythingOfType("int")).
			Return(patient, nil).
			Once()

		result, err := business.FindPatientById(1)
		assert.NoError(t, err)
		assert.Equal(t, patient, result)
	})

	t.Run("valid - when SelectPatientById return error", func(t *testing.T) {
		repo.
			On("SelectPatientById", mock.AnythingOfType("int")).
			Return(patients.PatientCore{}, errServer).
			Once()

		result, err := business.FindPatientById(1)
		assert.Error(t, err)
		assert.Equal(t, patients.PatientCore{}, result)
	})
}

func TestCreatePatient(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin, nil).
			Once()

		repo.
			On("SelectPatientByNIK", mock.AnythingOfType("string")).
			Return(patients.PatientCore{}, errNotFound).
			Once()

		repo.
			On("InsertPatient", mock.AnythingOfType("patients.PatientCore")).
			Return(nil).
			Once()

		err := business.CreatePatient(patient)
		assert.NoError(t, err)
	})

	t.Run("valid - when FindAdminById return error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errServer).
			Once()

		err := business.CreatePatient(patient)
		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when duplicate NIK", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin, nil).
			Once()

		repo.
			On("SelectPatientByNIK", mock.AnythingOfType("string")).
			Return(patients.PatientCore{NIK: patient.NIK}, nil).
			Once()

		err := business.CreatePatient(patient)
		assert.Error(t, err)
		assert.Equal(t, errors.KindUnprocessable, errors.Kind(err))
	})

	t.Run("valid - when SelectPatientByNIK return error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin, nil).
			Once()

		repo.
			On("SelectPatientByNIK", mock.AnythingOfType("string")).
			Return(patients.PatientCore{}, errServer).
			Once()

		err := business.CreatePatient(patient)
		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when InsertPatient return error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin, nil).
			Once()

		repo.
			On("SelectPatientByNIK", mock.AnythingOfType("string")).
			Return(patients.PatientCore{}, errNotFound).
			Once()

		repo.
			On("InsertPatient", mock.AnythingOfType("patients.PatientCore")).
			Return(errServer).
			Once()

		err := business.CreatePatient(patient)
		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})
}

func TestEditPatient(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin, nil).
			Once()

		repo.
			On("SelectPatientById", mock.AnythingOfType("int")).
			Return(patient, nil).
			Once()

		repo.
			On("UpdatePatient", mock.AnythingOfType("patients.PatientCore")).
			Return(nil).
			Once()

		err := business.EditPatient(patient)
		assert.NoError(t, err)
	})

	t.Run("valid - when FindAdminById return error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errServer).
			Once()

		err := business.EditPatient(patient)
		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when SelectPatientById return error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin, nil).
			Once()

		repo.
			On("SelectPatientById", mock.AnythingOfType("int")).
			Return(patients.PatientCore{}, errNotFound).
			Once()

		err := business.EditPatient(patient)
		assert.Error(t, err)
		assert.Equal(t, errors.KindNotFound, errors.Kind(err))
	})

	t.Run("valid - UpdatePatient return error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin, nil).
			Once()

		repo.
			On("SelectPatientById", mock.AnythingOfType("int")).
			Return(patient, nil).
			Once()

		repo.
			On("UpdatePatient", mock.AnythingOfType("patients.PatientCore")).
			Return(errServer).
			Once()

		err := business.EditPatient(patient)
		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})
}

func TestRemovePatientById(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin, nil).
			Once()

		repo.
			On("DeletePatientById", mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(nil).
			Once()

		err := business.RemovePatientById(patient.ID, admin.ID)
		assert.NoError(t, err)
	})

	t.Run("valid - when FindAdminById return error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		err := business.RemovePatientById(patient.ID, admin.ID)
		assert.Error(t, err)
		assert.Equal(t, errors.KindNotFound, errors.Kind(err))
	})

	t.Run("valid - when DeletePatientById return error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admin, nil).
			Once()

		repo.
			On("DeletePatientById", mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(errServer).
			Once()

		err := business.RemovePatientById(patient.ID, admin.ID)
		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})
}
