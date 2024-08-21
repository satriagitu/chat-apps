// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	domain "chat-apps/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// FileRepository is an autogenerated mock type for the FileRepository type
type FileRepository struct {
	mock.Mock
}

// GetFileByID provides a mock function with given fields: id
func (_m *FileRepository) GetFileByID(id int) (domain.File, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetFileByID")
	}

	var r0 domain.File
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (domain.File, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) domain.File); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(domain.File)
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UploadFile provides a mock function with given fields: file
func (_m *FileRepository) UploadFile(file domain.File) (domain.File, error) {
	ret := _m.Called(file)

	if len(ret) == 0 {
		panic("no return value specified for UploadFile")
	}

	var r0 domain.File
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.File) (domain.File, error)); ok {
		return rf(file)
	}
	if rf, ok := ret.Get(0).(func(domain.File) domain.File); ok {
		r0 = rf(file)
	} else {
		r0 = ret.Get(0).(domain.File)
	}

	if rf, ok := ret.Get(1).(func(domain.File) error); ok {
		r1 = rf(file)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFileRepository creates a new instance of FileRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFileRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *FileRepository {
	mock := &FileRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
