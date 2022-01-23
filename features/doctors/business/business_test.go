package business_test

import (
	"os"
	"testing"

	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/admins"
	am "github.com/final-project-alterra/hospital-management-system-api/features/admins/mocks"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
	d "github.com/final-project-alterra/hospital-management-system-api/features/doctors/business"
	dm "github.com/final-project-alterra/hospital-management-system-api/features/doctors/mocks"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
	nm "github.com/final-project-alterra/hospital-management-system-api/features/nurses/mocks"
	sm "github.com/final-project-alterra/hospital-management-system-api/features/schedules/mocks"
	"github.com/final-project-alterra/hospital-management-system-api/utils/files"
	"github.com/final-project-alterra/hospital-management-system-api/utils/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	doctorData dm.IData

	adminBusiness    am.IBusiness
	doctorBusiness   doctors.IBusiness
	nurseBusiness    nm.IBusiness
	scheduleBusiness sm.IBusiness

	adminMaster admins.AdminCore
	doctorHan   doctors.DoctorCore
	room1       doctors.RoomCore
	speciality1 doctors.SpecialityCore

	errNotFound      error
	errUnprocessable error
	errServer        error
)

func TestMain(m *testing.M) {
	doctorBusiness = d.NewDoctorBusinessBuilder().
		SetData(&doctorData).
		SetAdminBusiness(&adminBusiness).
		SetNurseBusiness(&nurseBusiness).
		SetScheduleBusiness(&scheduleBusiness).
		Build()

	adminMaster = admins.AdminCore{
		ID:   1,
		Name: "Master admin",
	}

	room1 = doctors.RoomCore{ID: 1, Floor: "1", Code: "A"}
	speciality1 = doctors.SpecialityCore{ID: 1, Name: "Heart"}

	doctorHan = doctors.DoctorCore{
		ID:         1,
		Room:       room1,
		Speciality: speciality1,
		Name:       "Han Solo",
		Email:      "han@mail.com",
		Password:   "12345678",
	}

	errNotFound = errors.E(errors.New("not found"), errors.KindNotFound)
	errUnprocessable = errors.E(errors.New("not found"), errors.KindUnprocessable)
	errServer = errors.E(errors.New("server error"), errors.KindServerError)

	files.Remove = func(path string) error { return nil }

	os.Exit(m.Run())
}

func TestFindDoctors(t *testing.T) {
	t.Run("valid - find doctors", func(t *testing.T) {
		doctorData.
			On("SelectDoctors").
			Return([]doctors.DoctorCore{doctorHan}, nil).
			Once()

		result, err := doctorBusiness.FindDoctors()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})

	t.Run("valid - error occurs on find doctors", func(t *testing.T) {
		doctorData.
			On("SelectDoctors").
			Return([]doctors.DoctorCore{}, errServer).
			Once()

		result, err := doctorBusiness.FindDoctors()

		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
		assert.Equal(t, 0, len(result))
	})
}

func TestDoctosrByIds(t *testing.T) {
	t.Run("valid - find doctors by ids", func(t *testing.T) {
		doctorData.
			On("SelectDoctorsByIds", mock.AnythingOfType("[]int")).
			Return([]doctors.DoctorCore{doctorHan}, nil).
			Once()

		result, err := doctorBusiness.FindDoctorsByIds([]int{1})

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})

	t.Run("valid - error occurs on find doctors by ids", func(t *testing.T) {
		doctorData.
			On("SelectDoctorsByIds", mock.AnythingOfType("[]int")).
			Return([]doctors.DoctorCore{}, errServer).
			Once()

		result, err := doctorBusiness.FindDoctorsByIds([]int{1})

		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
		assert.Equal(t, 0, len(result))
	})
}

func TestFindDoctorById(t *testing.T) {
	t.Run("valid - find doctor by id", func(t *testing.T) {
		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctorHan, nil).
			Once()

		result, err := doctorBusiness.FindDoctorById(1)

		assert.Nil(t, err)
		assert.Equal(t, doctorHan, result)
	})

	t.Run("valid - error occurs on find doctor by id", func(t *testing.T) {
		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctors.DoctorCore{}, errServer).
			Once()

		result, err := doctorBusiness.FindDoctorById(1)

		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
		assert.Equal(t, doctors.DoctorCore{}, result)
	})
}

func TestFindDoctorByEmail(t *testing.T) {
	t.Run("valid - FindDoctorByEmail", func(t *testing.T) {
		doctorData.
			On("SelectDoctorByEmail", mock.AnythingOfType("string")).
			Return(doctorHan, nil).
			Once()

		result, err := doctorBusiness.FindDoctorByEmail("han@mail.com")

		assert.Nil(t, err)
		assert.Equal(t, doctorHan, result)
	})

	t.Run("valid - when error occurs on FindDoctorByEmail", func(t *testing.T) {
		doctorData.
			On("SelectDoctorByEmail", mock.AnythingOfType("string")).
			Return(doctors.DoctorCore{}, errServer).
			Once()

		result, err := doctorBusiness.FindDoctorByEmail("han@mail.com")

		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
		assert.Equal(t, doctors.DoctorCore{}, result)
	})
}

func TestCreateDoctor(t *testing.T) {
	t.Run("valid - CreateDoctor", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(adminMaster, nil).
			Once()

		doctorData.
			On("SelectSpecialityById", mock.AnythingOfType("int")).
			Return(speciality1, nil).
			Once()

		doctorData.
			On("SelectRoomById", mock.AnythingOfType("int")).
			Return(room1, nil).
			Once()

		doctorData.
			On("SelectDoctorByEmail", mock.AnythingOfType("string")).
			Return(doctors.DoctorCore{}, errNotFound).
			Once()

		adminBusiness.
			On("FindAdminByEmail", mock.AnythingOfType("string")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		nurseBusiness.
			On("FindNurseByEmail", mock.AnythingOfType("string")).
			Return(nurses.NurseCore{}, errNotFound).
			Once()

		doctorData.
			On("InsertDoctor", mock.AnythingOfType("doctors.DoctorCore")).
			Return(nil).
			Once()

		newDoctor := doctorHan
		newDoctor.ID = 0
		newDoctor.CreatedBy = adminMaster.ID
		err := doctorBusiness.CreateDoctor(newDoctor)

		assert.Nil(t, err)
	})

	t.Run("valid - when creating admin not found", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		newDoctor := doctorHan
		newDoctor.ID = 0
		newDoctor.CreatedBy = adminMaster.ID
		err := doctorBusiness.CreateDoctor(newDoctor)

		assert.Error(t, err)
		assert.Equal(t, errors.KindNotFound, errors.Kind(err))
	})

	t.Run("valid - when error on finding creating admin", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errServer).
			Once()

		newDoctor := doctorHan
		newDoctor.ID = 0
		newDoctor.CreatedBy = adminMaster.ID
		err := doctorBusiness.CreateDoctor(newDoctor)

		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when SelectSpecialityById error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(adminMaster, nil).
			Once()

		doctorData.
			On("SelectSpecialityById", mock.AnythingOfType("int")).
			Return(doctors.SpecialityCore{}, errNotFound).
			Once()

		newDoctor := doctorHan
		newDoctor.ID = 0
		newDoctor.CreatedBy = adminMaster.ID
		err := doctorBusiness.CreateDoctor(newDoctor)

		assert.Equal(t, errors.KindNotFound, errors.Kind(err))
	})

	t.Run("valid - when SelectRoomById error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(adminMaster, nil).
			Once()

		doctorData.
			On("SelectSpecialityById", mock.AnythingOfType("int")).
			Return(speciality1, nil).
			Once()

		doctorData.
			On("SelectRoomById", mock.AnythingOfType("int")).
			Return(doctors.RoomCore{}, errNotFound).
			Once()

		newDoctor := doctorHan
		newDoctor.ID = 0
		newDoctor.CreatedBy = adminMaster.ID
		err := doctorBusiness.CreateDoctor(newDoctor)

		assert.Equal(t, errors.KindNotFound, errors.Kind(err))
	})

	t.Run("valid - when duplicate email happens", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(adminMaster, nil).
			Times(3)

		doctorData.
			On("SelectSpecialityById", mock.AnythingOfType("int")).
			Return(speciality1, nil).
			Times(3)

		doctorData.
			On("SelectRoomById", mock.AnythingOfType("int")).
			Return(room1, nil).
			Times(3)

		for i := 0; i < 3; i++ {
			switch i {
			case 0:
				adminBusiness.
					On("FindAdminByEmail", mock.AnythingOfType("string")).
					Return(admins.AdminCore{ID: 1}, nil).
					Once()
				doctorData.
					On("SelectDoctorByEmail", mock.AnythingOfType("string")).
					Return(doctors.DoctorCore{}, errNotFound).
					Once()
				nurseBusiness.
					On("FindNurseByEmail", mock.AnythingOfType("string")).
					Return(nurses.NurseCore{}, errNotFound).
					Once()
			case 1:
				adminBusiness.
					On("FindAdminByEmail", mock.AnythingOfType("string")).
					Return(admins.AdminCore{}, errNotFound).
					Once()
				doctorData.
					On("SelectDoctorByEmail", mock.AnythingOfType("string")).
					Return(doctors.DoctorCore{ID: 1}, nil).
					Once()
				nurseBusiness.
					On("FindNurseByEmail", mock.AnythingOfType("string")).
					Return(nurses.NurseCore{}, errNotFound).
					Once()
			case 2:
				adminBusiness.
					On("FindAdminByEmail", mock.AnythingOfType("string")).
					Return(admins.AdminCore{}, errNotFound).
					Once()
				doctorData.
					On("SelectDoctorByEmail", mock.AnythingOfType("string")).
					Return(doctors.DoctorCore{}, errNotFound).
					Once()
				nurseBusiness.
					On("FindNurseByEmail", mock.AnythingOfType("string")).
					Return(nurses.NurseCore{ID: 1}, nil).
					Once()
			}

			newDoctor := doctorHan
			newDoctor.ID = 0
			newDoctor.CreatedBy = adminMaster.ID
			err := doctorBusiness.CreateDoctor(newDoctor)

			assert.Equal(t, errors.KindUnprocessable, errors.Kind(err))
		}
	})

	t.Run("valid - when InsertDoctor error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(adminMaster, nil).
			Once()

		doctorData.
			On("SelectSpecialityById", mock.AnythingOfType("int")).
			Return(speciality1, nil).
			Once()

		doctorData.
			On("SelectRoomById", mock.AnythingOfType("int")).
			Return(room1, nil).
			Once()

		doctorData.
			On("SelectDoctorByEmail", mock.AnythingOfType("string")).
			Return(doctors.DoctorCore{}, errNotFound).
			Once()

		adminBusiness.
			On("FindAdminByEmail", mock.AnythingOfType("string")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		nurseBusiness.
			On("FindNurseByEmail", mock.AnythingOfType("string")).
			Return(nurses.NurseCore{}, errNotFound).
			Once()

		doctorData.
			On("InsertDoctor", mock.AnythingOfType("doctors.DoctorCore")).
			Return(errServer).
			Once()

		newDoctor := doctorHan
		newDoctor.ID = 0
		newDoctor.CreatedBy = adminMaster.ID
		err := doctorBusiness.CreateDoctor(newDoctor)

		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})
}

func TestEditDoctor(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctorHan, nil).
			Once()

		doctorData.
			On("SelectSpecialityById", mock.AnythingOfType("int")).
			Return(speciality1, nil).
			Once()

		doctorData.
			On("SelectRoomById", mock.AnythingOfType("int")).
			Return(room1, nil).
			Once()

		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(adminMaster, nil).
			Once()

		doctorData.
			On("UpdateDoctor", mock.AnythingOfType("doctors.DoctorCore")).
			Return(nil).
			Once()

		updatedDoctor := doctorHan
		updatedDoctor.ID = 1
		updatedDoctor.CreatedBy = adminMaster.ID
		updatedDoctor.UpdatedBy = adminMaster.ID
		err := doctorBusiness.EditDoctor(updatedDoctor)

		assert.Nil(t, err)
	})

	t.Run("valid - when SelectDoctorById error", func(t *testing.T) {
		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctors.DoctorCore{}, errServer).
			Once()

		updatedDoctor := doctorHan
		updatedDoctor.ID = 1
		updatedDoctor.CreatedBy = adminMaster.ID
		updatedDoctor.UpdatedBy = adminMaster.ID
		err := doctorBusiness.EditDoctor(updatedDoctor)

		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when SelectSpecialityById error", func(t *testing.T) {
		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctorHan, nil).
			Once()

		doctorData.
			On("SelectSpecialityById", mock.AnythingOfType("int")).
			Return(doctors.SpecialityCore{}, errServer).
			Once()

		updatedDoctor := doctorHan
		updatedDoctor.ID = 1
		updatedDoctor.CreatedBy = adminMaster.ID
		updatedDoctor.UpdatedBy = adminMaster.ID
		err := doctorBusiness.EditDoctor(updatedDoctor)

		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when SelectRoomById error", func(t *testing.T) {
		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctorHan, nil).
			Once()

		doctorData.
			On("SelectSpecialityById", mock.AnythingOfType("int")).
			Return(speciality1, nil).
			Once()

		doctorData.
			On("SelectRoomById", mock.AnythingOfType("int")).
			Return(doctors.RoomCore{}, errServer).
			Once()

		updatedDoctor := doctorHan
		updatedDoctor.ID = 1
		updatedDoctor.CreatedBy = adminMaster.ID
		updatedDoctor.UpdatedBy = adminMaster.ID
		err := doctorBusiness.EditDoctor(updatedDoctor)

		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when FindAdminById does not find admin", func(t *testing.T) {
		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctorHan, nil).
			Once()

		doctorData.
			On("SelectSpecialityById", mock.AnythingOfType("int")).
			Return(speciality1, nil).
			Once()

		doctorData.
			On("SelectRoomById", mock.AnythingOfType("int")).
			Return(room1, nil).
			Once()

		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		updatedDoctor := doctorHan
		updatedDoctor.ID = 1
		updatedDoctor.CreatedBy = adminMaster.ID
		updatedDoctor.UpdatedBy = adminMaster.ID
		err := doctorBusiness.EditDoctor(updatedDoctor)

		assert.Equal(t, errors.KindNotFound, errors.Kind(err))
	})

	t.Run("valid - when FindAdminById error", func(t *testing.T) {
		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctorHan, nil).
			Once()

		doctorData.
			On("SelectSpecialityById", mock.AnythingOfType("int")).
			Return(speciality1, nil).
			Once()

		doctorData.
			On("SelectRoomById", mock.AnythingOfType("int")).
			Return(room1, nil).
			Once()

		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errServer).
			Once()

		updatedDoctor := doctorHan
		updatedDoctor.ID = 1
		updatedDoctor.CreatedBy = adminMaster.ID
		updatedDoctor.UpdatedBy = adminMaster.ID
		err := doctorBusiness.EditDoctor(updatedDoctor)

		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when UpdateDoctor error", func(t *testing.T) {
		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctorHan, nil).
			Once()

		doctorData.
			On("SelectSpecialityById", mock.AnythingOfType("int")).
			Return(speciality1, nil).
			Once()

		doctorData.
			On("SelectRoomById", mock.AnythingOfType("int")).
			Return(room1, nil).
			Once()

		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(adminMaster, nil).
			Once()

		doctorData.
			On("UpdateDoctor", mock.AnythingOfType("doctors.DoctorCore")).
			Return(errServer).
			Once()

		updatedDoctor := doctorHan
		updatedDoctor.ID = 1
		updatedDoctor.CreatedBy = adminMaster.ID
		updatedDoctor.UpdatedBy = adminMaster.ID
		err := doctorBusiness.EditDoctor(updatedDoctor)

		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})
}

func TestEditDoctorPassword(t *testing.T) {
	t.Run("valid - when eveerything is fine", func(t *testing.T) {
		doctorHanRecord := doctorHan
		hashed, err := hash.Generate(doctorHan.Password)
		if err != nil {
			panic(err)
		}
		doctorHanRecord.Password = hashed

		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(adminMaster, nil).
			Once()

		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctorHanRecord, nil).
			Once()

		doctorData.
			On("UpdateDoctor", mock.AnythingOfType("doctors.DoctorCore")).
			Return(nil).
			Once()

		err = doctorBusiness.EditDoctorPassword(doctorHan.ID, adminMaster.ID, doctorHan.Password, "new password")

		assert.Nil(t, err)
	})

	t.Run("valid - when UpdateDoctor error", func(t *testing.T) {
		doctorHanRecord := doctorHan
		hashed, err := hash.Generate(doctorHan.Password)
		if err != nil {
			panic(err)
		}
		doctorHanRecord.Password = hashed

		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(adminMaster, nil).
			Once()

		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctorHanRecord, nil).
			Once()

		doctorData.
			On("UpdateDoctor", mock.AnythingOfType("doctors.DoctorCore")).
			Return(errServer).
			Once()

		err = doctorBusiness.EditDoctorPassword(doctorHan.ID, adminMaster.ID, doctorHan.Password, "new password")

		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when FindAdminById didnt find admin", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		err := doctorBusiness.EditDoctorPassword(doctorHan.ID, adminMaster.ID, doctorHan.Password, "new password")

		assert.Equal(t, errors.KindNotFound, errors.Kind(err))
	})

	t.Run("valid - when FindAdminById error", func(t *testing.T) {

		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errServer).
			Once()

		err := doctorBusiness.EditDoctorPassword(doctorHan.ID, adminMaster.ID, doctorHan.Password, "new password")

		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when SelectDoctorById didnt find doctor", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(adminMaster, nil).
			Once()

		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctors.DoctorCore{}, errNotFound).
			Once()

		err := doctorBusiness.EditDoctorPassword(doctorHan.ID, adminMaster.ID, doctorHan.Password, "new password")

		assert.Equal(t, errors.KindNotFound, errors.Kind(err))
	})

	t.Run("valid - when SelectDoctorById error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(adminMaster, nil).
			Once()

		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctors.DoctorCore{}, errServer).
			Once()

		err := doctorBusiness.EditDoctorPassword(doctorHan.ID, adminMaster.ID, doctorHan.Password, "new password")

		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when old password does not match", func(t *testing.T) {
		doctorHanRecord := doctorHan
		hashed, err := hash.Generate(doctorHan.Password)
		if err != nil {
			panic(err)
		}
		doctorHanRecord.Password = hashed

		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(adminMaster, nil).
			Once()

		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctorHanRecord, nil).
			Once()

		doctorData.
			On("UpdateDoctor", mock.AnythingOfType("doctors.DoctorCore")).
			Return(nil).
			Once()

		err = doctorBusiness.EditDoctorPassword(doctorHan.ID, adminMaster.ID, "wrong old password", "new password")

		assert.Equal(t, errors.KindUnprocessable, errors.Kind(err))
	})
}

func TestEditDoctorImageProfile(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(adminMaster, nil).
			Once()

		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctorHan, nil).
			Once()

		doctorData.
			On("UpdateDoctor", mock.Anything).
			Return(nil).
			Once()

		err := doctorBusiness.EditDoctorImageProfile(doctorHan)
		assert.Nil(t, err)
	})

	t.Run("valid - when FindAdminById error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errServer).
			Once()

		err := doctorBusiness.EditDoctorImageProfile(doctorHan)
		assert.Error(t, err)
	})

	t.Run("valid - when SelectDoctorById error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(adminMaster, nil).
			Once()

		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctors.DoctorCore{}, errServer).
			Once()

		err := doctorBusiness.EditDoctorImageProfile(doctorHan)
		assert.Error(t, err)
	})

}

func TestRemoveDoctorById(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(adminMaster, nil).
			Once()

		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctorHan, nil).
			Once()

		scheduleBusiness.
			On("RemoveDoctorFutureWorkSchedules", mock.AnythingOfType("int")).
			Return(nil).
			Once()

		doctorData.
			On("DeleteDoctorById", mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(nil).
			Once()

		err := doctorBusiness.RemoveDoctorById(doctorHan.ID, adminMaster.ID)

		assert.Nil(t, err)
	})

	t.Run("valid - when FindAdminById didnt find admin", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errNotFound).
			Once()

		err := doctorBusiness.RemoveDoctorById(doctorHan.ID, adminMaster.ID)

		assert.Equal(t, errors.KindNotFound, errors.Kind(err))
	})

	t.Run("valid - when FindAdminById error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(admins.AdminCore{}, errServer).
			Once()

		err := doctorBusiness.RemoveDoctorById(doctorHan.ID, adminMaster.ID)

		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when SelectDoctorById error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(adminMaster, nil).
			Once()

		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctorHan, errServer).
			Once()

		err := doctorBusiness.RemoveDoctorById(doctorHan.ID, adminMaster.ID)

		assert.Error(t, err)
	})

	t.Run("valid - when RemoveDoctorFutureWorkSchedules error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(adminMaster, nil).
			Once()

		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctorHan, nil).
			Once()

		scheduleBusiness.
			On("RemoveDoctorFutureWorkSchedules", mock.AnythingOfType("int")).
			Return(errServer).
			Once()

		err := doctorBusiness.RemoveDoctorById(doctorHan.ID, adminMaster.ID)

		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when DeleteDoctorById error", func(t *testing.T) {
		adminBusiness.
			On("FindAdminById", mock.AnythingOfType("int")).
			Return(adminMaster, nil).
			Once()

		doctorData.
			On("SelectDoctorById", mock.AnythingOfType("int")).
			Return(doctorHan, nil).
			Once()

		scheduleBusiness.
			On("RemoveDoctorFutureWorkSchedules", mock.AnythingOfType("int")).
			Return(nil).
			Once()

		doctorData.
			On("DeleteDoctorById", mock.AnythingOfType("int"), mock.AnythingOfType("int")).
			Return(errServer).
			Once()

		err := doctorBusiness.RemoveDoctorById(doctorHan.ID, adminMaster.ID)

		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

}

func TestFindSpecialities(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		doctorData.
			On("SelectSpecialities").
			Return([]doctors.SpecialityCore{speciality1}, nil).
			Once()

		specialities, err := doctorBusiness.FindSpecialities()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(specialities))
	})

	t.Run("valid - when SelectSpecialities error", func(t *testing.T) {
		doctorData.
			On("SelectSpecialities").
			Return([]doctors.SpecialityCore{}, errServer).
			Once()

		specialities, err := doctorBusiness.FindSpecialities()

		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
		assert.Equal(t, 0, len(specialities))
	})
}

func TestFindSpecialityById(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		doctorData.
			On("SelectSpecialityById", mock.AnythingOfType("int")).
			Return(speciality1, nil).
			Once()

		result, err := doctorBusiness.FindSpecialityById(1)

		assert.Nil(t, err)
		assert.Equal(t, speciality1, result)
	})

	t.Run("valid - when SelectSpecialityById error", func(t *testing.T) {
		doctorData.
			On("SelectSpecialityById", mock.AnythingOfType("int")).
			Return(doctors.SpecialityCore{}, errNotFound).
			Once()

		result, err := doctorBusiness.FindSpecialityById(1)

		assert.Error(t, err)
		assert.Equal(t, errors.KindNotFound, errors.Kind(err))
		assert.Equal(t, doctors.SpecialityCore{}, result)
	})
}

func TestCreateSpeciality(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		doctorData.
			On("InsertSpeciality", mock.AnythingOfType("doctors.SpecialityCore")).
			Return(nil).
			Once()

		err := doctorBusiness.CreateSpeciality(speciality1)

		assert.Nil(t, err)
	})

	t.Run("valid - when InsertSpeciality error", func(t *testing.T) {
		doctorData.
			On("InsertSpeciality", mock.AnythingOfType("doctors.SpecialityCore")).
			Return(errServer).
			Once()

		err := doctorBusiness.CreateSpeciality(speciality1)

		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})
}

func TestEditSpeciality(t *testing.T) {
	t.Run("valid - when everyhing is fine", func(t *testing.T) {
		doctorData.
			On("SelectSpecialityById", mock.AnythingOfType("int")).
			Return(doctors.SpecialityCore{}, nil).
			Once()

		doctorData.
			On("UpdateSpeciality", mock.AnythingOfType("doctors.SpecialityCore")).
			Return(nil).
			Once()

		err := doctorBusiness.EditSpeciality(speciality1)

		assert.Nil(t, err)
	})

	t.Run("valid - when SelectSpecialityById error", func(t *testing.T) {
		doctorData.
			On("SelectSpecialityById", mock.AnythingOfType("int")).
			Return(doctors.SpecialityCore{}, errNotFound).
			Once()

		err := doctorBusiness.EditSpeciality(speciality1)

		assert.Error(t, err)
	})

	t.Run("valid - when UpdateSpeciality error", func(t *testing.T) {
		doctorData.
			On("SelectSpecialityById", mock.AnythingOfType("int")).
			Return(doctors.SpecialityCore{}, nil).
			Once()

		doctorData.
			On("UpdateSpeciality", mock.AnythingOfType("doctors.SpecialityCore")).
			Return(errServer).
			Once()

		err := doctorBusiness.EditSpeciality(speciality1)

		assert.Error(t, err)
	})
}

func TestRemoveSpeciality(t *testing.T) {
	t.Run("valid - when everyhing is fine", func(t *testing.T) {
		doctorData.
			On("SelectDoctorsBySpecialityId", mock.AnythingOfType("int")).
			Return([]doctors.DoctorCore{}, nil).
			Once()

		doctorData.
			On("DeleteSpecialityId", mock.AnythingOfType("int")).
			Return(nil).
			Once()

		err := doctorBusiness.RemoveSpeciality(1)

		assert.Nil(t, err)
	})

	t.Run("valid - SelectDoctorsBySpecialityId return error", func(t *testing.T) {
		doctorData.
			On("SelectDoctorsBySpecialityId", mock.AnythingOfType("int")).
			Return([]doctors.DoctorCore{}, errServer).
			Once()

		err := doctorBusiness.RemoveSpeciality(1)

		assert.Error(t, err)
	})

	t.Run("valid - when there exists doctor with that speciality", func(t *testing.T) {
		doctorData.
			On("SelectDoctorsBySpecialityId", mock.AnythingOfType("int")).
			Return([]doctors.DoctorCore{doctorHan}, nil).
			Once()

		err := doctorBusiness.RemoveSpeciality(1)

		assert.Error(t, err)
	})

	t.Run("valid - when DeleteSpecialityId error", func(t *testing.T) {
		doctorData.
			On("SelectDoctorsBySpecialityId", mock.AnythingOfType("int")).
			Return([]doctors.DoctorCore{}, nil).
			Once()

		doctorData.
			On("DeleteSpecialityId", mock.AnythingOfType("int")).
			Return(errServer).
			Once()

		err := doctorBusiness.RemoveSpeciality(1)

		assert.Error(t, err)
	})

}

func TestFindRooms(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		doctorData.
			On("SelectRooms").
			Return([]doctors.RoomCore{room1}, nil).
			Once()

		rooms, err := doctorBusiness.FindRooms()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(rooms))
	})

	t.Run("valid - when everything is fine", func(t *testing.T) {
		doctorData.
			On("SelectRooms").
			Return([]doctors.RoomCore{}, errServer).
			Once()

		rooms, err := doctorBusiness.FindRooms()

		assert.Error(t, err)
		assert.Equal(t, 0, len(rooms))
	})
}

func TestCreateRoom(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		doctorData.
			On("SelectRoomByCode", mock.AnythingOfType("string")).
			Return(doctors.RoomCore{}, errNotFound).
			Once()

		doctorData.
			On("InsertRoom", mock.AnythingOfType("doctors.RoomCore")).
			Return(nil).
			Once()

		err := doctorBusiness.CreateRoom(room1)

		assert.Nil(t, err)
	})

	t.Run("valid - when found a room with the same code", func(t *testing.T) {
		doctorData.
			On("SelectRoomByCode", mock.AnythingOfType("string")).
			Return(doctors.RoomCore{ID: 2, Code: room1.Code}, nil).
			Once()

		err := doctorBusiness.CreateRoom(room1)

		assert.Error(t, err)
		assert.Equal(t, errors.KindUnprocessable, errors.Kind(err))
	})

	t.Run("valid - when SelectRoomByCode error", func(t *testing.T) {
		doctorData.
			On("SelectRoomByCode", mock.AnythingOfType("string")).
			Return(doctors.RoomCore{}, errServer).
			Once()

		err := doctorBusiness.CreateRoom(room1)

		assert.Error(t, err)
		assert.Equal(t, errors.KindServerError, errors.Kind(err))
	})

	t.Run("valid - when InsertRoom error", func(t *testing.T) {
		doctorData.
			On("SelectRoomByCode", mock.AnythingOfType("string")).
			Return(doctors.RoomCore{}, errNotFound).
			Once()

		doctorData.
			On("InsertRoom", mock.AnythingOfType("doctors.RoomCore")).
			Return(errServer).
			Once()

		err := doctorBusiness.CreateRoom(room1)

		assert.Error(t, err)
	})
}

func TestEditRoom(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		doctorData.
			On("SelectRoomById", mock.AnythingOfType("int")).
			Return(room1, nil).
			Once()

		doctorData.
			On("SelectRoomByCode", mock.AnythingOfType("string")).
			Return(doctors.RoomCore{}, errNotFound).
			Once()

		doctorData.
			On("UpdateRoom", mock.AnythingOfType("doctors.RoomCore")).
			Return(nil).
			Once()

		err := doctorBusiness.EditRoom(room1)

		assert.Nil(t, err)
	})

	t.Run("valid - when cheching existing room error", func(t *testing.T) {
		doctorData.
			On("SelectRoomById", mock.AnythingOfType("int")).
			Return(doctors.RoomCore{}, errServer).
			Once()

		err := doctorBusiness.EditRoom(room1)

		assert.Error(t, err)
	})

	t.Run("valid - when code already used by other room", func(t *testing.T) {
		doctorData.
			On("SelectRoomById", mock.AnythingOfType("int")).
			Return(room1, nil).
			Once()

		doctorData.
			On("SelectRoomByCode", mock.AnythingOfType("string")).
			Return(doctors.RoomCore{ID: 2, Code: room1.Code}, nil).
			Once()

		err := doctorBusiness.EditRoom(room1)

		assert.Error(t, err)
	})

	t.Run("valid - when SelectRoomByCode error", func(t *testing.T) {
		doctorData.
			On("SelectRoomById", mock.AnythingOfType("int")).
			Return(room1, nil).
			Once()

		doctorData.
			On("SelectRoomByCode", mock.AnythingOfType("string")).
			Return(doctors.RoomCore{}, errServer).
			Once()

		err := doctorBusiness.EditRoom(room1)

		assert.Error(t, err)
	})

	t.Run("valid - when everything is fine", func(t *testing.T) {
		doctorData.
			On("SelectRoomById", mock.AnythingOfType("int")).
			Return(room1, nil).
			Once()

		doctorData.
			On("SelectRoomByCode", mock.AnythingOfType("string")).
			Return(doctors.RoomCore{}, errNotFound).
			Once()

		doctorData.
			On("UpdateRoom", mock.AnythingOfType("doctors.RoomCore")).
			Return(errServer).
			Once()

		err := doctorBusiness.EditRoom(room1)

		assert.Error(t, err)
	})

}

func TestDeleteRoom(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		doctorData.
			On("SelectDoctorsByRoomId", mock.AnythingOfType("int")).
			Return([]doctors.DoctorCore{}, nil).
			Once()

		doctorData.
			On("DeleteRoomById", mock.AnythingOfType("int")).
			Return(nil).
			Once()

		err := doctorBusiness.RemoveRoomById(1)
		assert.Nil(t, err)
	})

	t.Run("valid - when SelectDoctorsByRoomId return error", func(t *testing.T) {
		doctorData.
			On("SelectDoctorsByRoomId", mock.AnythingOfType("int")).
			Return([]doctors.DoctorCore{}, errServer).
			Once()

		err := doctorBusiness.RemoveRoomById(1)
		assert.Error(t, err)
	})

	t.Run("valid - when there are still doctor using the room", func(t *testing.T) {
		doctorData.
			On("SelectDoctorsByRoomId", mock.AnythingOfType("int")).
			Return([]doctors.DoctorCore{doctorHan}, nil).
			Once()

		err := doctorBusiness.RemoveRoomById(1)
		assert.Error(t, err)
	})

	t.Run("valid - when DeleteRoomById return error", func(t *testing.T) {
		doctorData.
			On("SelectDoctorsByRoomId", mock.AnythingOfType("int")).
			Return([]doctors.DoctorCore{}, nil).
			Once()

		doctorData.
			On("DeleteRoomById", mock.AnythingOfType("int")).
			Return(errServer).
			Once()

		err := doctorBusiness.RemoveRoomById(1)
		assert.Error(t, err)
	})

}
