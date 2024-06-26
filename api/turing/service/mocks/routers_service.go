// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	models "github.com/caraml-dev/turing/api/turing/models"
	mock "github.com/stretchr/testify/mock"
)

// RoutersService is an autogenerated mock type for the RoutersService type
type RoutersService struct {
	mock.Mock
}

// Delete provides a mock function with given fields: router
func (_m *RoutersService) Delete(router *models.Router) error {
	ret := _m.Called(router)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Router) error); ok {
		r0 = rf(router)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByID provides a mock function with given fields: routerID
func (_m *RoutersService) FindByID(routerID models.ID) (*models.Router, error) {
	ret := _m.Called(routerID)

	var r0 *models.Router
	var r1 error
	if rf, ok := ret.Get(0).(func(models.ID) (*models.Router, error)); ok {
		return rf(routerID)
	}
	if rf, ok := ret.Get(0).(func(models.ID) *models.Router); ok {
		r0 = rf(routerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Router)
		}
	}

	if rf, ok := ret.Get(1).(func(models.ID) error); ok {
		r1 = rf(routerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByProjectAndEnvironmentAndName provides a mock function with given fields: projectID, environmentName, routerName
func (_m *RoutersService) FindByProjectAndEnvironmentAndName(projectID models.ID, environmentName string, routerName string) (*models.Router, error) {
	ret := _m.Called(projectID, environmentName, routerName)

	var r0 *models.Router
	var r1 error
	if rf, ok := ret.Get(0).(func(models.ID, string, string) (*models.Router, error)); ok {
		return rf(projectID, environmentName, routerName)
	}
	if rf, ok := ret.Get(0).(func(models.ID, string, string) *models.Router); ok {
		r0 = rf(projectID, environmentName, routerName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Router)
		}
	}

	if rf, ok := ret.Get(1).(func(models.ID, string, string) error); ok {
		r1 = rf(projectID, environmentName, routerName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByProjectAndName provides a mock function with given fields: projectID, routerName
func (_m *RoutersService) FindByProjectAndName(projectID models.ID, routerName string) (*models.Router, error) {
	ret := _m.Called(projectID, routerName)

	var r0 *models.Router
	var r1 error
	if rf, ok := ret.Get(0).(func(models.ID, string) (*models.Router, error)); ok {
		return rf(projectID, routerName)
	}
	if rf, ok := ret.Get(0).(func(models.ID, string) *models.Router); ok {
		r0 = rf(projectID, routerName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Router)
		}
	}

	if rf, ok := ret.Get(1).(func(models.ID, string) error); ok {
		r1 = rf(projectID, routerName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListRouters provides a mock function with given fields: projectID, environmentName
func (_m *RoutersService) ListRouters(projectID models.ID, environmentName string) ([]*models.Router, error) {
	ret := _m.Called(projectID, environmentName)

	var r0 []*models.Router
	var r1 error
	if rf, ok := ret.Get(0).(func(models.ID, string) ([]*models.Router, error)); ok {
		return rf(projectID, environmentName)
	}
	if rf, ok := ret.Get(0).(func(models.ID, string) []*models.Router); ok {
		r0 = rf(projectID, environmentName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Router)
		}
	}

	if rf, ok := ret.Get(1).(func(models.ID, string) error); ok {
		r1 = rf(projectID, environmentName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: router
func (_m *RoutersService) Save(router *models.Router) (*models.Router, error) {
	ret := _m.Called(router)

	var r0 *models.Router
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.Router) (*models.Router, error)); ok {
		return rf(router)
	}
	if rf, ok := ret.Get(0).(func(*models.Router) *models.Router); ok {
		r0 = rf(router)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Router)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.Router) error); ok {
		r1 = rf(router)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRoutersService interface {
	mock.TestingT
	Cleanup(func())
}

// NewRoutersService creates a new instance of RoutersService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRoutersService(t mockConstructorTestingTNewRoutersService) *RoutersService {
	mock := &RoutersService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
