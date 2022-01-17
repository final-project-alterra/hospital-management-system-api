package business_test

import (
	"os"
	"testing"

	"github.com/final-project-alterra/hospital-management-system-api/config"
	"github.com/final-project-alterra/hospital-management-system-api/errors"

	d "github.com/final-project-alterra/hospital-management-system-api/features/doctors"
	n "github.com/final-project-alterra/hospital-management-system-api/features/nurses"

	p "github.com/final-project-alterra/hospital-management-system-api/features/patients"
	s "github.com/final-project-alterra/hospital-management-system-api/features/schedules"

	dm "github.com/final-project-alterra/hospital-management-system-api/features/doctors/mocks"
	nm "github.com/final-project-alterra/hospital-management-system-api/features/nurses/mocks"
	pm "github.com/final-project-alterra/hospital-management-system-api/features/patients/mocks"
	sm "github.com/final-project-alterra/hospital-management-system-api/features/schedules/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	sb "github.com/final-project-alterra/hospital-management-system-api/features/schedules/business"
)

var (
	repo     sm.IData
	business s.IBusiness

	doctorBusiness  dm.IBusiness
	nurseBusiness   nm.IBusiness
	patientBusiness pm.IBusiness

	// emptyPrescription s.PrescriptionCore
	// emptyOutpatient   s.OutpatientCore
	// emptyWorkSchedule s.WorkScheduleCore
	// emptyPatient      s.PatientCore
	// emptyNurse        s.NurseCore
	// emptyDoctor       s.DoctorCore
	// emptyRoom         s.RoomCore

	doctorCore1  d.DoctorCore
	nurseCore1   n.NurseCore
	patientCore1 p.PatientCore

	q             s.ScheduleQuery
	prescription1 s.PrescriptionCore
	outpatient1   s.OutpatientCore
	workSchedule1 s.WorkScheduleCore
	patient1      s.PatientCore
	nurse1        s.NurseCore
	doctor1       s.DoctorCore
	room1         s.RoomCore

	any         string
	anyInt      mock.AnythingOfTypeArgument
	anySliceInt mock.AnythingOfTypeArgument

	errNotFound error
	errServer   error
)

func TestMain(m *testing.M) {
	config.LoadENV("../../../aws.env")
	config.InitTimeLoc("Asia/Jakarta")

	business = sb.NewScheduleBusinessBuilder().
		SetData(&repo).
		SetDoctorBusiness(&doctorBusiness).
		SetNurseBusiness(&nurseBusiness).
		SetPatientBusiness(&patientBusiness).
		Build()

	doctorCore1 = d.DoctorCore{ID: 1}
	nurseCore1 = n.NurseCore{ID: 1}
	patientCore1 = p.PatientCore{ID: 1}

	prescription1 = s.PrescriptionCore{ID: 1}
	room1 = s.RoomCore{ID: 1}

	doctor1 = s.DoctorCore{ID: 1, Room: room1}
	nurse1 = s.NurseCore{ID: 1}
	patient1 = s.PatientCore{ID: 1}

	outpatient1 = s.OutpatientCore{
		ID:            1,
		Patient:       patient1,
		Prescriptions: []s.PrescriptionCore{prescription1},
	}

	workSchedule1 = s.WorkScheduleCore{
		ID:          1,
		Doctor:      doctor1,
		Nurse:       nurse1,
		Date:        "2100-01-01",
		StartTime:   "00:00:00",
		EndTime:     "12:00:00",
		Outpatients: []s.OutpatientCore{outpatient1},
	}

	any = mock.Anything
	anyInt = mock.AnythingOfType("int")
	anySliceInt = mock.AnythingOfType("[]int")

	errNotFound = errors.E(errors.New("not found"), errors.KindNotFound)
	errServer = errors.E(errors.New("server"), errors.KindServerError)

	os.Exit(m.Run())
}

func TestFindWorkSchedules(t *testing.T) {
	t.Run("valid - everything is fine", func(t *testing.T) {
		repo.
			On("SelectWorkSchedules", any).
			Return([]s.WorkScheduleCore{workSchedule1}, nil).
			Once()

		repo.
			On("SelectCountWorkSchedulesWaitings", anySliceInt).
			Return(map[int]int{}, nil).
			Once()

		doctorBusiness.
			On("FindDoctosrByIds", anySliceInt).
			Return([]d.DoctorCore{{ID: 1}}, nil).
			Once()

		nurseBusiness.
			On("FindNursesByIds", anySliceInt).
			Return([]n.NurseCore{{ID: 1}}, nil).
			Once()

		result, err := business.FindWorkSchedules(q)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})

	t.Run("valid - SelectWorkSchedules error", func(t *testing.T) {
		repo.
			On("SelectWorkSchedules", any).
			Return([]s.WorkScheduleCore{}, errServer).
			Once()

		result, err := business.FindWorkSchedules(q)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("valid - SelectCountWorkSchedulesWaitings error", func(t *testing.T) {
		repo.
			On("SelectWorkSchedules", any).
			Return([]s.WorkScheduleCore{workSchedule1}, nil).
			Once()

		repo.
			On("SelectCountWorkSchedulesWaitings", anySliceInt).
			Return(map[int]int{}, errServer).
			Once()

		result, err := business.FindWorkSchedules(q)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("valid - FindDoctosrByIds error", func(t *testing.T) {
		repo.
			On("SelectWorkSchedules", any).
			Return([]s.WorkScheduleCore{workSchedule1}, nil).
			Once()

		repo.
			On("SelectCountWorkSchedulesWaitings", anySliceInt).
			Return(map[int]int{}, nil).
			Once()

		doctorBusiness.
			On("FindDoctosrByIds", anySliceInt).
			Return([]d.DoctorCore{}, errServer).
			Once()

		result, err := business.FindWorkSchedules(q)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("valid - FindNursesByIds error", func(t *testing.T) {
		repo.
			On("SelectWorkSchedules", any).
			Return([]s.WorkScheduleCore{workSchedule1}, nil).
			Once()

		repo.
			On("SelectCountWorkSchedulesWaitings", anySliceInt).
			Return(map[int]int{}, nil).
			Once()

		doctorBusiness.
			On("FindDoctosrByIds", anySliceInt).
			Return([]d.DoctorCore{{ID: 1}}, nil).
			Once()

		nurseBusiness.
			On("FindNursesByIds", anySliceInt).
			Return([]n.NurseCore{}, errServer).
			Once()

		result, err := business.FindWorkSchedules(q)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})
}

func TestFindDoctorWorkSchedules(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectWorkSchedulesByDoctorId", anyInt, any).
			Return([]s.WorkScheduleCore{workSchedule1}, nil).
			Once()

		nurseBusiness.
			On("FindNursesByIds", anySliceInt).
			Return([]n.NurseCore{{ID: 1}}, nil).
			Once()

		result, err := business.FindDoctorWorkSchedules(1, q)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})

	t.Run("valid - SelectWorkSchedulesByDoctorId error", func(t *testing.T) {
		repo.
			On("SelectWorkSchedulesByDoctorId", anyInt, any).
			Return([]s.WorkScheduleCore{}, errServer).
			Once()

		result, err := business.FindDoctorWorkSchedules(1, q)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("valid - FindNursesByIds error", func(t *testing.T) {
		repo.
			On("SelectWorkSchedulesByDoctorId", anyInt, any).
			Return([]s.WorkScheduleCore{workSchedule1}, nil).
			Once()

		nurseBusiness.
			On("FindNursesByIds", anySliceInt).
			Return([]n.NurseCore{}, errServer).
			Once()

		result, err := business.FindDoctorWorkSchedules(1, q)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})
}

func TestFindNurseWorkSchedules(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectWorkSchedulesByNurseId", anyInt, q).
			Return([]s.WorkScheduleCore{workSchedule1}, nil).
			Once()

		doctorBusiness.
			On("FindDoctosrByIds", anySliceInt).
			Return([]d.DoctorCore{{ID: 1}}, nil).
			Once()

		result, err := business.FindNurseWorkSchedules(1, q)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})

	t.Run("valid - SelectWorkSchedulesByNurseId error", func(t *testing.T) {
		repo.
			On("SelectWorkSchedulesByNurseId", anyInt, q).
			Return([]s.WorkScheduleCore{}, errServer).
			Once()

		result, err := business.FindNurseWorkSchedules(1, q)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("valid - FindDoctosrByIds error", func(t *testing.T) {
		repo.
			On("SelectWorkSchedulesByNurseId", anyInt, q).
			Return([]s.WorkScheduleCore{workSchedule1}, nil).
			Once()

		doctorBusiness.
			On("FindDoctosrByIds", anySliceInt).
			Return([]d.DoctorCore{}, errServer).
			Once()

		result, err := business.FindNurseWorkSchedules(1, q)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})
}

func TestCreateWorkSchedule(t *testing.T) {
	t.Run("valid - InsertWorkSchedules error", func(t *testing.T) {
		doctorBusiness.
			On("FindDoctorById", anyInt).
			Return(doctorCore1, nil).
			Once()

		nurseBusiness.
			On("FindNurseById", anyInt).
			Return(nurseCore1, nil).
			Once()

		repo.
			On("InsertWorkSchedules", any).
			Return(errServer).
			Once()

		q := s.ScheduleQuery{Repeat: s.RepeatNoRepeat}
		err := business.CreateWorkSchedule(workSchedule1, q)
		assert.Error(t, err)
	})

	repeatTest := 9

	doctorBusiness.
		On("FindDoctorById", anyInt).
		Return(doctorCore1, nil).
		Times(repeatTest)

	nurseBusiness.
		On("FindNurseById", anyInt).
		Return(nurseCore1, nil).
		Times(repeatTest)

	repo.
		On("InsertWorkSchedules", any).
		Return(nil).
		Times(repeatTest)

	t.Run("valid - when everything is fine for no-repeat", func(t *testing.T) {
		q := s.ScheduleQuery{Repeat: s.RepeatNoRepeat}
		err := business.CreateWorkSchedule(workSchedule1, q)
		assert.Nil(t, err)
	})
	t.Run("valid - when everything is fine for daily", func(t *testing.T) {
		q := s.ScheduleQuery{
			Repeat:    s.RepeatDaily,
			StartDate: "2020-01-01",
			EndDate:   "2020-02-01",
		}
		err := business.CreateWorkSchedule(workSchedule1, q)
		assert.Nil(t, err)
	})
	t.Run("valid - when everything is fine for weekly", func(t *testing.T) {
		q := s.ScheduleQuery{
			Repeat:    s.RepeatWeekly,
			StartDate: "2020-01-01",
			EndDate:   "2020-04-01",
		}
		err := business.CreateWorkSchedule(workSchedule1, q)
		assert.Nil(t, err)
	})
	t.Run("valid - when everything is fine for monthly", func(t *testing.T) {
		q := s.ScheduleQuery{
			Repeat:    s.RepeatMonthly,
			StartDate: "2020-01-01",
			EndDate:   "2020-04-01",
		}
		err := business.CreateWorkSchedule(workSchedule1, q)
		assert.Nil(t, err)
	})

	t.Run("valid - for unknown repeat", func(t *testing.T) {
		q := s.ScheduleQuery{Repeat: "some-unknown-repeat"}
		err := business.CreateWorkSchedule(workSchedule1, q)
		assert.Error(t, err)
	})

	t.Run("valid - for invalid start date on repeat daily", func(t *testing.T) {
		q := s.ScheduleQuery{
			Repeat:    s.RepeatDaily,
			StartDate: "invalid start date",
			EndDate:   "2020-10-20",
		}
		err := business.CreateWorkSchedule(workSchedule1, q)
		assert.Error(t, err)
	})
	t.Run("valid - for invalid end date on repeat daily", func(t *testing.T) {
		q := s.ScheduleQuery{
			Repeat:    s.RepeatDaily,
			StartDate: "2020-10-20",
			EndDate:   "invalid end date",
		}
		err := business.CreateWorkSchedule(workSchedule1, q)
		assert.Error(t, err)
	})
	t.Run("valid - for invalid date on repeat weekly", func(t *testing.T) {
		q := s.ScheduleQuery{
			Repeat:    s.RepeatWeekly,
			StartDate: "2020-10-20",
			EndDate:   "invalid end date",
		}
		err := business.CreateWorkSchedule(workSchedule1, q)
		assert.Error(t, err)
	})
	t.Run("valid - for invalid date on repeat monthly", func(t *testing.T) {
		q := s.ScheduleQuery{
			Repeat:    s.RepeatMonthly,
			StartDate: "2020-10-20",
			EndDate:   "invalid end date",
		}
		err := business.CreateWorkSchedule(workSchedule1, q)
		assert.Error(t, err)
	})

	t.Run("valid - FindDoctorById error", func(t *testing.T) {
		doctorBusiness.
			On("FindDoctorById", anyInt).
			Return(d.DoctorCore{}, errNotFound).
			Once()

		err := business.CreateWorkSchedule(workSchedule1, q)
		assert.Error(t, err)
	})

	t.Run("valid - FindNurseById error", func(t *testing.T) {
		doctorBusiness.
			On("FindDoctorById", anyInt).
			Return(doctorCore1, nil).
			Once()

		nurseBusiness.
			On("FindNurseById", anyInt).
			Return(n.NurseCore{}, errNotFound).
			Once()

		err := business.CreateWorkSchedule(workSchedule1, q)
		assert.Error(t, err)
	})
}

func TestEditWorkSchedule(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectWorkScheduleById", anyInt).
			Return(workSchedule1, nil).
			Once()

		doctorBusiness.
			On("FindDoctorById", anyInt).
			Return(doctorCore1, nil).
			Once()

		nurseBusiness.
			On("FindNurseById", anyInt).
			Return(nurseCore1, nil).
			Once()

		repo.
			On("UpdateWorkSchedule", any).
			Return(nil).
			Once()

		err := business.EditWorkSchedule(workSchedule1)
		assert.Nil(t, err)
	})

	t.Run("valid - SelectWorkScheduleById error", func(t *testing.T) {
		repo.
			On("SelectWorkScheduleById", anyInt).
			Return(s.WorkScheduleCore{}, errServer).
			Once()

		err := business.EditWorkSchedule(workSchedule1)
		assert.Error(t, err)
	})

	t.Run("valid - FindDoctorById error", func(t *testing.T) {
		repo.
			On("SelectWorkScheduleById", anyInt).
			Return(workSchedule1, nil).
			Once()

		doctorBusiness.
			On("FindDoctorById", anyInt).
			Return(d.DoctorCore{}, errServer).
			Once()

		err := business.EditWorkSchedule(workSchedule1)
		assert.Error(t, err)
	})

	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectWorkScheduleById", anyInt).
			Return(workSchedule1, nil).
			Once()

		doctorBusiness.
			On("FindDoctorById", anyInt).
			Return(doctorCore1, nil).
			Once()

		nurseBusiness.
			On("FindNurseById", anyInt).
			Return(n.NurseCore{}, errServer).
			Once()

		err := business.EditWorkSchedule(workSchedule1)
		assert.Error(t, err)
	})

	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectWorkScheduleById", anyInt).
			Return(workSchedule1, nil).
			Once()

		doctorBusiness.
			On("FindDoctorById", anyInt).
			Return(doctorCore1, nil).
			Once()

		nurseBusiness.
			On("FindNurseById", anyInt).
			Return(nurseCore1, nil).
			Once()

		repo.
			On("UpdateWorkSchedule", any).
			Return(errServer).
			Once()

		err := business.EditWorkSchedule(workSchedule1)
		assert.Error(t, err)
	})
}

func TestRemoveWorkScheduleById(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("DeleteWorkScheduleById", anyInt).
			Return(nil).
			Once()

		err := business.RemoveWorkScheduleById(1)
		assert.Nil(t, err)
	})

	t.Run("valid - DeleteWorkScheduleById error", func(t *testing.T) {
		repo.
			On("DeleteWorkScheduleById", anyInt).
			Return(errServer).
			Once()

		err := business.RemoveWorkScheduleById(1)
		assert.Error(t, err)
	})
}

func TestRemoveDoctorFutureWorkSchedules(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("DeleteWorkSchedulesByDoctorId", anyInt, any).
			Return(nil).
			Once()

		err := business.RemoveDoctorFutureWorkSchedules(1)
		assert.Nil(t, err)
	})

	t.Run("valid - DeleteWorkSchedulesByDoctorId error", func(t *testing.T) {
		repo.
			On("DeleteWorkSchedulesByDoctorId", anyInt, any).
			Return(errServer).
			Once()

		err := business.RemoveDoctorFutureWorkSchedules(1)
		assert.Error(t, err)
	})
}

func TestRemoveNurseFromNextWorkSchedules(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("DeleteNurseFromWorkSchedules", anyInt, any).
			Return(nil).
			Once()

		err := business.RemoveNurseFromNextWorkSchedules(1)
		assert.Nil(t, err)
	})

	t.Run("valid - DeleteNurseFromWorkSchedules error", func(t *testing.T) {
		repo.
			On("DeleteNurseFromWorkSchedules", anyInt, any).
			Return(errServer).
			Once()

		err := business.RemoveNurseFromNextWorkSchedules(1)
		assert.Error(t, err)
	})
}

func TestFindOutpatietns(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectOutpatients", any).
			Return([]s.OutpatientCore{outpatient1}, nil).
			Once()

		patientBusiness.
			On("FindPatientsByIds", anySliceInt).
			Return([]p.PatientCore{patientCore1}, nil).
			Once()

		doctorBusiness.
			On("FindDoctosrByIds", anySliceInt).
			Return([]d.DoctorCore{doctorCore1}, nil).
			Once()

		nurseBusiness.
			On("FindNursesByIds", anySliceInt).
			Return([]n.NurseCore{nurseCore1}, nil).
			Once()

		result, err := business.FindOutpatietns(q)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})

	t.Run("valid - SelectOutpatients error", func(t *testing.T) {
		repo.
			On("SelectOutpatients", any).
			Return([]s.OutpatientCore{}, errServer).
			Once()

		result, err := business.FindOutpatietns(q)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("valid - FindPatientsByIds error", func(t *testing.T) {
		repo.
			On("SelectOutpatients", any).
			Return([]s.OutpatientCore{outpatient1}, nil).
			Once()

		patientBusiness.
			On("FindPatientsByIds", anySliceInt).
			Return([]p.PatientCore{}, errServer).
			Once()

		result, err := business.FindOutpatietns(q)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("valid - FindDoctosrByIds error", func(t *testing.T) {
		repo.
			On("SelectOutpatients", any).
			Return([]s.OutpatientCore{outpatient1}, nil).
			Once()

		patientBusiness.
			On("FindPatientsByIds", anySliceInt).
			Return([]p.PatientCore{patientCore1}, nil).
			Once()

		doctorBusiness.
			On("FindDoctosrByIds", anySliceInt).
			Return([]d.DoctorCore{}, errServer).
			Once()

		result, err := business.FindOutpatietns(q)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("valid - FindNursesByIds error", func(t *testing.T) {
		repo.
			On("SelectOutpatients", any).
			Return([]s.OutpatientCore{outpatient1}, nil).
			Once()

		patientBusiness.
			On("FindPatientsByIds", anySliceInt).
			Return([]p.PatientCore{patientCore1}, nil).
			Once()

		doctorBusiness.
			On("FindDoctosrByIds", anySliceInt).
			Return([]d.DoctorCore{doctorCore1}, nil).
			Once()

		nurseBusiness.
			On("FindNursesByIds", anySliceInt).
			Return([]n.NurseCore{}, errServer).
			Once()

		result, err := business.FindOutpatietns(q)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})
}

func TestFindOutpatietnsByWorkScheduleId(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectOutpatientsByWorkScheduleId", anyInt).
			Return(workSchedule1, nil).
			Once()

		doctorBusiness.
			On("FindDoctorById", anyInt).
			Return(doctorCore1, nil).
			Once()

		nurseBusiness.
			On("FindNurseById", anyInt).
			Return(nurseCore1, nil).
			Once()

		patientBusiness.
			On("FindPatientsByIds", anySliceInt).
			Return([]p.PatientCore{patientCore1}, nil).
			Once()

		result, err := business.FindOutpatietnsByWorkScheduleId(1)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result.Outpatients))
		assert.Equal(t, doctorCore1.ID, result.Doctor.ID)
		assert.Equal(t, nurseCore1.ID, result.Nurse.ID)
	})

	t.Run("valid - SelectOutpatientsByWorkScheduleId error", func(t *testing.T) {
		repo.
			On("SelectOutpatientsByWorkScheduleId", anyInt).
			Return(s.WorkScheduleCore{}, errNotFound).
			Once()

		result, err := business.FindOutpatietnsByWorkScheduleId(1)

		assert.Error(t, err)
		assert.Equal(t, 0, result.ID)
		assert.Equal(t, 0, len(result.Outpatients))
	})

	t.Run("valid - FindDoctorById error", func(t *testing.T) {
		repo.
			On("SelectOutpatientsByWorkScheduleId", anyInt).
			Return(workSchedule1, nil).
			Once()

		doctorBusiness.
			On("FindDoctorById", anyInt).
			Return(d.DoctorCore{}, errNotFound).
			Once()

		result, err := business.FindOutpatietnsByWorkScheduleId(1)

		assert.Error(t, err)
		assert.Equal(t, 0, result.ID)
		assert.Equal(t, 0, len(result.Outpatients))
	})

	t.Run("valid - FindNurseById error", func(t *testing.T) {
		repo.
			On("SelectOutpatientsByWorkScheduleId", anyInt).
			Return(workSchedule1, nil).
			Once()

		doctorBusiness.
			On("FindDoctorById", anyInt).
			Return(doctorCore1, nil).
			Once()

		nurseBusiness.
			On("FindNurseById", anyInt).
			Return(n.NurseCore{}, errServer).
			Once()

		result, err := business.FindOutpatietnsByWorkScheduleId(1)

		assert.Error(t, err)
		assert.Equal(t, 0, result.ID)
		assert.Equal(t, 0, len(result.Outpatients))
	})

	t.Run("valid - FindPatientsByIds error", func(t *testing.T) {
		repo.
			On("SelectOutpatientsByWorkScheduleId", anyInt).
			Return(workSchedule1, nil).
			Once()

		doctorBusiness.
			On("FindDoctorById", anyInt).
			Return(doctorCore1, nil).
			Once()

		nurseBusiness.
			On("FindNurseById", anyInt).
			Return(nurseCore1, nil).
			Once()

		patientBusiness.
			On("FindPatientsByIds", anySliceInt).
			Return([]p.PatientCore{}, errServer).
			Once()

		result, err := business.FindOutpatietnsByWorkScheduleId(1)

		assert.Error(t, err)
		assert.Equal(t, 0, result.ID)
		assert.Equal(t, 0, len(result.Outpatients))
	})
}

func TestFindOutpatientsByPatientId(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectOutpatientsByPatientId", anyInt, any).
			Return([]s.OutpatientCore{outpatient1}, nil).
			Once()

		doctorBusiness.
			On("FindDoctosrByIds", anySliceInt).
			Return([]d.DoctorCore{doctorCore1}, nil).
			Once()

		nurseBusiness.
			On("FindNursesByIds", anySliceInt).
			Return([]n.NurseCore{nurseCore1}, nil).
			Once()

		result, err := business.FindOutpatientsByPatientId(1, q)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})

	t.Run("valid - SelectOutpatientsByPatientId error", func(t *testing.T) {
		repo.
			On("SelectOutpatientsByPatientId", anyInt, any).
			Return([]s.OutpatientCore{}, errServer).
			Once()

		result, err := business.FindOutpatientsByPatientId(1, q)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("valid - FindDoctosrByIds error", func(t *testing.T) {
		repo.
			On("SelectOutpatientsByPatientId", anyInt, any).
			Return([]s.OutpatientCore{outpatient1}, nil).
			Once()

		doctorBusiness.
			On("FindDoctosrByIds", anySliceInt).
			Return([]d.DoctorCore{}, errServer).
			Once()

		result, err := business.FindOutpatientsByPatientId(1, q)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("valid - FindNursesByIds error", func(t *testing.T) {
		repo.
			On("SelectOutpatientsByPatientId", anyInt, any).
			Return([]s.OutpatientCore{outpatient1}, nil).
			Once()

		doctorBusiness.
			On("FindDoctosrByIds", anySliceInt).
			Return([]d.DoctorCore{doctorCore1}, nil).
			Once()

		nurseBusiness.
			On("FindNursesByIds", anySliceInt).
			Return([]n.NurseCore{}, errServer).
			Once()

		result, err := business.FindOutpatientsByPatientId(1, q)

		assert.Error(t, err)
		assert.Equal(t, 0, len(result))
	})
}

func TestFindOutpatientById(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(outpatient1, nil).
			Once()

		patientBusiness.
			On("FindPatientById", anyInt).
			Return(patientCore1, nil).
			Once()

		doctorBusiness.
			On("FindDoctorById", anyInt).
			Return(doctorCore1, nil).
			Once()

		nurseBusiness.
			On("FindNurseById", anyInt).
			Return(nurseCore1, nil).
			Once()

		result, err := business.FindOutpatientById(1)

		assert.Nil(t, err)
		assert.Equal(t, 1, result.ID)
	})

	t.Run("valid - SelectOutpatientById error", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(s.OutpatientCore{}, errServer).
			Once()

		result, err := business.FindOutpatientById(1)

		assert.Error(t, err)
		assert.Equal(t, 0, result.ID)
	})

	t.Run("valid - FindPatientById error", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(outpatient1, nil).
			Once()

		patientBusiness.
			On("FindPatientById", anyInt).
			Return(p.PatientCore{}, errServer).
			Once()

		result, err := business.FindOutpatientById(1)

		assert.Error(t, err)
		assert.Equal(t, 0, result.ID)
	})

	t.Run("valid - FindDoctorById error", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(outpatient1, nil).
			Once()

		patientBusiness.
			On("FindPatientById", anyInt).
			Return(patientCore1, nil).
			Once()

		doctorBusiness.
			On("FindDoctorById", anyInt).
			Return(d.DoctorCore{}, errServer).
			Once()

		result, err := business.FindOutpatientById(1)

		assert.Error(t, err)
		assert.Equal(t, 0, result.ID)
	})

	t.Run("valid - FindNurseById error", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(outpatient1, nil).
			Once()

		patientBusiness.
			On("FindPatientById", anyInt).
			Return(patientCore1, nil).
			Once()

		doctorBusiness.
			On("FindDoctorById", anyInt).
			Return(doctorCore1, nil).
			Once()

		nurseBusiness.
			On("FindNurseById", anyInt).
			Return(n.NurseCore{}, errServer).
			Once()

		result, err := business.FindOutpatientById(1)

		assert.Error(t, err)
		assert.Equal(t, 0, result.ID)
	})
}

func TestCreateOutpatient(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectWorkScheduleById", anyInt).
			Return(workSchedule1, nil).
			Once()

		patientBusiness.
			On("FindPatientById", anyInt).
			Return(patientCore1, nil).
			Once()

		repo.
			On("InsertOutpatient", any).
			Return(nil).
			Once()

		err := business.CreateOutpatient(outpatient1)

		assert.Nil(t, err)
	})

	t.Run("valid - FindPatientById error", func(t *testing.T) {

		patientBusiness.
			On("FindPatientById", anyInt).
			Return(p.PatientCore{}, errNotFound).
			Once()

		err := business.CreateOutpatient(outpatient1)

		assert.Error(t, err)
	})

	t.Run("valid - SelectWorkScheduleById error", func(t *testing.T) {
		patientBusiness.
			On("FindPatientById", anyInt).
			Return(patientCore1, nil).
			Once()

		repo.
			On("SelectWorkScheduleById", anyInt).
			Return(s.WorkScheduleCore{}, errNotFound).
			Once()

		err := business.CreateOutpatient(outpatient1)

		assert.Error(t, err)
	})

	t.Run("valid - time.ParseInLocation error", func(t *testing.T) {
		invalidWorkschedule := workSchedule1
		invalidWorkschedule.StartTime = ""
		invalidWorkschedule.EndTime = ""

		repo.
			On("SelectWorkScheduleById", anyInt).
			Return(invalidWorkschedule, nil).
			Once()

		patientBusiness.
			On("FindPatientById", anyInt).
			Return(patientCore1, nil).
			Once()

		err := business.CreateOutpatient(outpatient1)

		assert.Error(t, err)
	})

	t.Run("valid - when work schedule has ended", func(t *testing.T) {
		endedSchedule := workSchedule1
		endedSchedule.Date = "1990-01-01"

		repo.
			On("SelectWorkScheduleById", anyInt).
			Return(endedSchedule, nil).
			Once()

		patientBusiness.
			On("FindPatientById", anyInt).
			Return(patientCore1, nil).
			Once()

		err := business.CreateOutpatient(outpatient1)

		assert.Error(t, err)
	})

	t.Run("valid - InsertOutpatient error", func(t *testing.T) {
		repo.
			On("SelectWorkScheduleById", anyInt).
			Return(workSchedule1, nil).
			Once()

		patientBusiness.
			On("FindPatientById", anyInt).
			Return(patientCore1, nil).
			Once()

		repo.
			On("InsertOutpatient", any).
			Return(errServer).
			Once()

		err := business.CreateOutpatient(outpatient1)

		assert.Error(t, err)
	})
}

func TestEditOutpatient(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(outpatient1, nil).
			Once()

		repo.
			On("UpdateOutpatient", any).
			Return(nil).
			Once()

		err := business.EditOutpatient(outpatient1)

		assert.Nil(t, err)
	})

	t.Run("valid - SelectOutpatientById error", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(s.OutpatientCore{}, errNotFound).
			Once()

		err := business.EditOutpatient(outpatient1)

		assert.Error(t, err)
	})

	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(outpatient1, nil).
			Once()

		repo.
			On("UpdateOutpatient", any).
			Return(errServer).
			Once()

		err := business.EditOutpatient(outpatient1)

		assert.Error(t, err)
	})
}

func TestExamineOutpatient(t *testing.T) {
	onprogress := s.OutpatientCore{
		Status: s.StatusOnprogress,
		WorkSchedule: s.WorkScheduleCore{
			ID:     2,
			Doctor: s.DoctorCore{ID: 2},
			Nurse:  s.NurseCore{ID: 2},
		},
	}
	waiting := s.OutpatientCore{
		Status: s.StatusWaiting,
		WorkSchedule: s.WorkScheduleCore{
			ID:     1,
			Doctor: doctor1,
			Nurse:  nurse1,
		},
	}

	t.Run("valid - for doctor when everything is fine", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(waiting, nil).
			Once()

		w := s.WorkScheduleCore{
			Outpatients: []s.OutpatientCore{waiting},
		}
		repo.
			On("SelectOutpatientsByWorkScheduleId", anyInt).
			Return(w, nil).
			Once()

		repo.
			On("UpdateOutpatient", any).
			Return(nil).
			Once()

		err := business.ExamineOutpatient(outpatient1.ID, doctor1.ID, "doctor")
		assert.Nil(t, err)
	})

	t.Run("valid - for nurse when everything is fine", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(waiting, nil).
			Once()

		w := s.WorkScheduleCore{
			Outpatients: []s.OutpatientCore{waiting},
		}
		repo.
			On("SelectOutpatientsByWorkScheduleId", anyInt).
			Return(w, nil).
			Once()

		repo.
			On("UpdateOutpatient", any).
			Return(nil).
			Once()

		err := business.ExamineOutpatient(outpatient1.ID, nurse1.ID, "nurse")
		assert.Nil(t, err)
	})

	t.Run("valid - when role is unknown", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(waiting, nil).
			Once()

		err := business.ExamineOutpatient(outpatient1.ID, 1, "unknown")
		assert.Error(t, err)
	})

	t.Run("valid - for doctor that is not his schedule", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(waiting, nil).
			Once()

		err := business.ExamineOutpatient(outpatient1.ID, 2, "doctor")
		assert.Error(t, err)
	})

	t.Run("valid - for nurse that is not her schedule", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(waiting, nil).
			Once()

		err := business.ExamineOutpatient(outpatient1.ID, 2, "nurse")
		assert.Error(t, err)
	})

	t.Run("valid - SelectOutpatientById error", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(s.OutpatientCore{}, errNotFound).
			Once()

		err := business.ExamineOutpatient(outpatient1.ID, doctor1.ID, "doctor")
		assert.Error(t, err)
	})

	t.Run("valid - when outpatient is not on waiting state", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(onprogress, nil).
			Once()

		err := business.ExamineOutpatient(outpatient1.ID, doctor1.ID, "doctor")
		assert.Error(t, err)
	})

	t.Run("valid - SelectOutpatientsByWorkScheduleId error", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(waiting, nil).
			Once()

		repo.
			On("SelectOutpatientsByWorkScheduleId", anyInt).
			Return(s.WorkScheduleCore{ID: 100}, errNotFound).
			Once()

		err := business.ExamineOutpatient(outpatient1.ID, doctor1.ID, "doctor")
		assert.Error(t, err)
	})

	t.Run("valid - when there's an ongoing examination", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(waiting, nil).
			Once()

		w := s.WorkScheduleCore{
			Outpatients: []s.OutpatientCore{onprogress},
		}
		repo.
			On("SelectOutpatientsByWorkScheduleId", anyInt).
			Return(w, nil).
			Once()

		err := business.ExamineOutpatient(outpatient1.ID, doctor1.ID, "doctor")
		assert.Error(t, err)
	})

	t.Run("valid - UpdateOutpatient error", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(waiting, nil).
			Once()

		w := s.WorkScheduleCore{
			Outpatients: []s.OutpatientCore{waiting},
		}
		repo.
			On("SelectOutpatientsByWorkScheduleId", anyInt).
			Return(w, nil).
			Once()

		repo.
			On("UpdateOutpatient", any).
			Return(errServer).
			Once()

		err := business.ExamineOutpatient(outpatient1.ID, doctor1.ID, "doctor")
		assert.Error(t, err)
	})
}

func TestFinishOutpatient(t *testing.T) {
	onprogress := s.OutpatientCore{
		Status:       s.StatusOnprogress,
		WorkSchedule: workSchedule1,
	}
	waiting := s.OutpatientCore{Status: s.StatusWaiting}

	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(onprogress, nil).
			Once()

		repo.
			On("UpdateOutpatient", any).
			Return(nil).
			Once()

		err := business.FinishOutpatient(onprogress, doctor1.ID, "doctor")
		assert.Nil(t, err)
	})

	t.Run("valid - when the doctor is not his work schedule", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(onprogress, nil).
			Once()

		err := business.FinishOutpatient(onprogress, 2, "doctor")
		assert.Error(t, err)
	})

	t.Run("valid - when role is not doctor", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(onprogress, nil).
			Once()

		err := business.FinishOutpatient(onprogress, 2, "admin")
		assert.Error(t, err)
	})

	t.Run("valid - SelectOutpatientById error", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(s.OutpatientCore{}, errNotFound).
			Once()

		err := business.FinishOutpatient(onprogress, doctor1.ID, "doctor")
		assert.Error(t, err)
	})

	t.Run("valid - when outpatient status is not onprogress", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(waiting, nil).
			Once()

		err := business.FinishOutpatient(waiting, doctor1.ID, "doctor")
		assert.Error(t, err)
	})

	t.Run("valid - UpdateOutpatient error", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(onprogress, nil).
			Once()

		repo.
			On("UpdateOutpatient", any).
			Return(errServer).
			Once()

		err := business.FinishOutpatient(onprogress, doctor1.ID, "doctor")
		assert.Error(t, err)
	})

}

func TestCancelOutpatient(t *testing.T) {
	onprogress := s.OutpatientCore{Status: s.StatusOnprogress}
	waiting := s.OutpatientCore{Status: s.StatusWaiting, WorkSchedule: workSchedule1}

	t.Run("valid - for admin when everything is fine", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(waiting, nil).
			Once()

		repo.
			On("UpdateOutpatient", any).
			Return(nil).
			Once()

		err := business.CancelOutpatient(waiting.ID, 1, "admin")
		assert.Nil(t, err)
	})

	t.Run("valid - for doctor when everything is fine", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(waiting, nil).
			Once()

		repo.
			On("UpdateOutpatient", any).
			Return(nil).
			Once()

		err := business.CancelOutpatient(waiting.ID, doctorCore1.ID, "doctor")
		assert.Nil(t, err)
	})

	t.Run("valid - for nurse when everything is fine", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(waiting, nil).
			Once()

		repo.
			On("UpdateOutpatient", any).
			Return(nil).
			Once()

		err := business.CancelOutpatient(waiting.ID, nurse1.ID, "nurse")
		assert.Nil(t, err)
	})

	t.Run("valid - when doctor is not his work schedule", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(waiting, nil).
			Once()

		err := business.CancelOutpatient(waiting.ID, 2, "doctor")
		assert.Error(t, err)
	})

	t.Run("valid - when nurse is not her work schedule", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(waiting, nil).
			Once()

		err := business.CancelOutpatient(waiting.ID, 2, "nurse")
		assert.Error(t, err)
	})

	t.Run("valid - when role is unknown", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(waiting, nil).
			Once()

		err := business.CancelOutpatient(waiting.ID, 2, "unkown")
		assert.Error(t, err)
	})

	t.Run("valid - SelectOutpatientById error", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(s.OutpatientCore{}, errNotFound).
			Once()

		err := business.CancelOutpatient(waiting.ID, doctor1.ID, "doctor")
		assert.Error(t, err)
	})

	t.Run("valid - when outpatient status is not waiting", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(onprogress, nil).
			Once()

		err := business.CancelOutpatient(onprogress.ID, doctor1.ID, "doctor")
		assert.Error(t, err)
	})

	t.Run("valid - UpdateOutpatient error", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(waiting, nil).
			Once()

		repo.
			On("UpdateOutpatient", any).
			Return(errServer).
			Once()

		err := business.CancelOutpatient(waiting.ID, doctor1.ID, "doctor")
		assert.Error(t, err)
	})
}

func TestRemoveOutpatientById(t *testing.T) {
	onprogress := s.OutpatientCore{Status: s.StatusOnprogress}
	waiting := s.OutpatientCore{Status: s.StatusWaiting}

	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(waiting, nil).
			Once()

		repo.
			On("DeleteOutpatientById", anyInt).
			Return(nil).
			Once()

		err := business.RemoveOutpatientById(waiting.ID)
		assert.Nil(t, err)
	})

	t.Run("valid - SelectOutpatientById error", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(s.OutpatientCore{}, errNotFound).
			Once()

		err := business.RemoveOutpatientById(waiting.ID)
		assert.Error(t, err)
	})

	t.Run("valid - when outpatient status is onprogress", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(onprogress, nil).
			Once()

		err := business.RemoveOutpatientById(onprogress.ID)
		assert.Error(t, err)
	})

	t.Run("valid - DeleteOutpatientById error", func(t *testing.T) {
		repo.
			On("SelectOutpatientById", anyInt).
			Return(waiting, nil).
			Once()

		repo.
			On("DeleteOutpatientById", anyInt).
			Return(errServer).
			Once()

		err := business.RemoveOutpatientById(waiting.ID)
		assert.Error(t, err)
	})
}

func TestRemovePatientWaitingOutpatients(t *testing.T) {
	t.Run("valid - when everything is fine", func(t *testing.T) {
		repo.
			On("DeleteWaitingOutpatientsByPatientId", anyInt).
			Return(nil).
			Once()

		err := business.RemovePatientWaitingOutpatients(1)
		assert.Nil(t, err)
	})

	t.Run("valid - DeleteWaitingOutpatientsByPatientId error", func(t *testing.T) {
		repo.
			On("DeleteWaitingOutpatientsByPatientId", anyInt).
			Return(errServer).
			Once()

		err := business.RemovePatientWaitingOutpatients(1)
		assert.Error(t, err)
	})
}
