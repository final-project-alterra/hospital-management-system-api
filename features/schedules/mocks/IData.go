// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	schedules "github.com/final-project-alterra/hospital-management-system-api/features/schedules"
	mock "github.com/stretchr/testify/mock"
)

// IData is an autogenerated mock type for the IData type
type IData struct {
	mock.Mock
}

// DeleteNurseFromWorkSchedules provides a mock function with given fields: nurseId, q
func (_m *IData) DeleteNurseFromWorkSchedules(nurseId int, q schedules.ScheduleQuery) error {
	ret := _m.Called(nurseId, q)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, schedules.ScheduleQuery) error); ok {
		r0 = rf(nurseId, q)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteOutpatientById provides a mock function with given fields: outpatientId
func (_m *IData) DeleteOutpatientById(outpatientId int) error {
	ret := _m.Called(outpatientId)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(outpatientId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteWaitingOutpatientsByPatientId provides a mock function with given fields: patientId
func (_m *IData) DeleteWaitingOutpatientsByPatientId(patientId int) error {
	ret := _m.Called(patientId)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(patientId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteWorkScheduleById provides a mock function with given fields: workScheduleId
func (_m *IData) DeleteWorkScheduleById(workScheduleId int) error {
	ret := _m.Called(workScheduleId)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(workScheduleId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteWorkSchedulesByDoctorId provides a mock function with given fields: doctorId, q
func (_m *IData) DeleteWorkSchedulesByDoctorId(doctorId int, q schedules.ScheduleQuery) error {
	ret := _m.Called(doctorId, q)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, schedules.ScheduleQuery) error); ok {
		r0 = rf(doctorId, q)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertOutpatient provides a mock function with given fields: outpatient
func (_m *IData) InsertOutpatient(outpatient schedules.OutpatientCore) error {
	ret := _m.Called(outpatient)

	var r0 error
	if rf, ok := ret.Get(0).(func(schedules.OutpatientCore) error); ok {
		r0 = rf(outpatient)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertWorkSchedules provides a mock function with given fields: workSchedules
func (_m *IData) InsertWorkSchedules(workSchedules []schedules.WorkScheduleCore) error {
	ret := _m.Called(workSchedules)

	var r0 error
	if rf, ok := ret.Get(0).(func([]schedules.WorkScheduleCore) error); ok {
		r0 = rf(workSchedules)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SelectCountWorkSchedulesWaitings provides a mock function with given fields: ids
func (_m *IData) SelectCountWorkSchedulesWaitings(ids []int) (map[int]int, error) {
	ret := _m.Called(ids)

	var r0 map[int]int
	if rf, ok := ret.Get(0).(func([]int) map[int]int); ok {
		r0 = rf(ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[int]int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]int) error); ok {
		r1 = rf(ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SelectOutpatientById provides a mock function with given fields: outpatientId
func (_m *IData) SelectOutpatientById(outpatientId int) (schedules.OutpatientCore, error) {
	ret := _m.Called(outpatientId)

	var r0 schedules.OutpatientCore
	if rf, ok := ret.Get(0).(func(int) schedules.OutpatientCore); ok {
		r0 = rf(outpatientId)
	} else {
		r0 = ret.Get(0).(schedules.OutpatientCore)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(outpatientId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SelectOutpatients provides a mock function with given fields: q
func (_m *IData) SelectOutpatients(q schedules.ScheduleQuery) ([]schedules.OutpatientCore, error) {
	ret := _m.Called(q)

	var r0 []schedules.OutpatientCore
	if rf, ok := ret.Get(0).(func(schedules.ScheduleQuery) []schedules.OutpatientCore); ok {
		r0 = rf(q)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]schedules.OutpatientCore)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(schedules.ScheduleQuery) error); ok {
		r1 = rf(q)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SelectOutpatientsByPatientId provides a mock function with given fields: patientId, q
func (_m *IData) SelectOutpatientsByPatientId(patientId int, q schedules.ScheduleQuery) ([]schedules.OutpatientCore, error) {
	ret := _m.Called(patientId, q)

	var r0 []schedules.OutpatientCore
	if rf, ok := ret.Get(0).(func(int, schedules.ScheduleQuery) []schedules.OutpatientCore); ok {
		r0 = rf(patientId, q)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]schedules.OutpatientCore)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, schedules.ScheduleQuery) error); ok {
		r1 = rf(patientId, q)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SelectOutpatientsByWorkScheduleId provides a mock function with given fields: workScheduleId
func (_m *IData) SelectOutpatientsByWorkScheduleId(workScheduleId int) (schedules.WorkScheduleCore, error) {
	ret := _m.Called(workScheduleId)

	var r0 schedules.WorkScheduleCore
	if rf, ok := ret.Get(0).(func(int) schedules.WorkScheduleCore); ok {
		r0 = rf(workScheduleId)
	} else {
		r0 = ret.Get(0).(schedules.WorkScheduleCore)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(workScheduleId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SelectWorkScheduleById provides a mock function with given fields: workScheduleId
func (_m *IData) SelectWorkScheduleById(workScheduleId int) (schedules.WorkScheduleCore, error) {
	ret := _m.Called(workScheduleId)

	var r0 schedules.WorkScheduleCore
	if rf, ok := ret.Get(0).(func(int) schedules.WorkScheduleCore); ok {
		r0 = rf(workScheduleId)
	} else {
		r0 = ret.Get(0).(schedules.WorkScheduleCore)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(workScheduleId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SelectWorkSchedules provides a mock function with given fields: q
func (_m *IData) SelectWorkSchedules(q schedules.ScheduleQuery) ([]schedules.WorkScheduleCore, error) {
	ret := _m.Called(q)

	var r0 []schedules.WorkScheduleCore
	if rf, ok := ret.Get(0).(func(schedules.ScheduleQuery) []schedules.WorkScheduleCore); ok {
		r0 = rf(q)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]schedules.WorkScheduleCore)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(schedules.ScheduleQuery) error); ok {
		r1 = rf(q)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SelectWorkSchedulesByDoctorId provides a mock function with given fields: doctorId, q
func (_m *IData) SelectWorkSchedulesByDoctorId(doctorId int, q schedules.ScheduleQuery) ([]schedules.WorkScheduleCore, error) {
	ret := _m.Called(doctorId, q)

	var r0 []schedules.WorkScheduleCore
	if rf, ok := ret.Get(0).(func(int, schedules.ScheduleQuery) []schedules.WorkScheduleCore); ok {
		r0 = rf(doctorId, q)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]schedules.WorkScheduleCore)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, schedules.ScheduleQuery) error); ok {
		r1 = rf(doctorId, q)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SelectWorkSchedulesByNurseId provides a mock function with given fields: nurseId, q
func (_m *IData) SelectWorkSchedulesByNurseId(nurseId int, q schedules.ScheduleQuery) ([]schedules.WorkScheduleCore, error) {
	ret := _m.Called(nurseId, q)

	var r0 []schedules.WorkScheduleCore
	if rf, ok := ret.Get(0).(func(int, schedules.ScheduleQuery) []schedules.WorkScheduleCore); ok {
		r0 = rf(nurseId, q)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]schedules.WorkScheduleCore)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, schedules.ScheduleQuery) error); ok {
		r1 = rf(nurseId, q)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOutpatient provides a mock function with given fields: outpatient
func (_m *IData) UpdateOutpatient(outpatient schedules.OutpatientCore) error {
	ret := _m.Called(outpatient)

	var r0 error
	if rf, ok := ret.Get(0).(func(schedules.OutpatientCore) error); ok {
		r0 = rf(outpatient)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateWorkSchedule provides a mock function with given fields: workSchedule
func (_m *IData) UpdateWorkSchedule(workSchedule schedules.WorkScheduleCore) error {
	ret := _m.Called(workSchedule)

	var r0 error
	if rf, ok := ret.Get(0).(func(schedules.WorkScheduleCore) error); ok {
		r0 = rf(workSchedule)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
