// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	admins "github.com/final-project-alterra/hospital-management-system-api/features/admins"
	mock "github.com/stretchr/testify/mock"
)

// IBusiness is an autogenerated mock type for the IBusiness type
type IBusiness struct {
	mock.Mock
}

// CreateAdmin provides a mock function with given fields: admin
func (_m *IBusiness) CreateAdmin(admin admins.AdminCore) error {
	ret := _m.Called(admin)

	var r0 error
	if rf, ok := ret.Get(0).(func(admins.AdminCore) error); ok {
		r0 = rf(admin)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EditAdmin provides a mock function with given fields: admin
func (_m *IBusiness) EditAdmin(admin admins.AdminCore) error {
	ret := _m.Called(admin)

	var r0 error
	if rf, ok := ret.Get(0).(func(admins.AdminCore) error); ok {
		r0 = rf(admin)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EditAdminPassword provides a mock function with given fields: id, updatedBy, oldPassword, newPassword
func (_m *IBusiness) EditAdminPassword(id int, updatedBy int, oldPassword string, newPassword string) error {
	ret := _m.Called(id, updatedBy, oldPassword, newPassword)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, int, string, string) error); ok {
		r0 = rf(id, updatedBy, oldPassword, newPassword)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EditAdminProfileImage provides a mock function with given fields: admin
func (_m *IBusiness) EditAdminProfileImage(admin admins.AdminCore) error {
	ret := _m.Called(admin)

	var r0 error
	if rf, ok := ret.Get(0).(func(admins.AdminCore) error); ok {
		r0 = rf(admin)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAdminByEmail provides a mock function with given fields: email
func (_m *IBusiness) FindAdminByEmail(email string) (admins.AdminCore, error) {
	ret := _m.Called(email)

	var r0 admins.AdminCore
	if rf, ok := ret.Get(0).(func(string) admins.AdminCore); ok {
		r0 = rf(email)
	} else {
		r0 = ret.Get(0).(admins.AdminCore)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAdminById provides a mock function with given fields: id
func (_m *IBusiness) FindAdminById(id int) (admins.AdminCore, error) {
	ret := _m.Called(id)

	var r0 admins.AdminCore
	if rf, ok := ret.Get(0).(func(int) admins.AdminCore); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(admins.AdminCore)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAdmins provides a mock function with given fields:
func (_m *IBusiness) FindAdmins() ([]admins.AdminCore, error) {
	ret := _m.Called()

	var r0 []admins.AdminCore
	if rf, ok := ret.Get(0).(func() []admins.AdminCore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]admins.AdminCore)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveAdminById provides a mock function with given fields: id, updatedBy
func (_m *IBusiness) RemoveAdminById(id int, updatedBy int) error {
	ret := _m.Called(id, updatedBy)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, int) error); ok {
		r0 = rf(id, updatedBy)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
