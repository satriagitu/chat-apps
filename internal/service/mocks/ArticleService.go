// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	domain "chat-apps/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// ArticleService is an autogenerated mock type for the ArticleService type
type ArticleService struct {
	mock.Mock
}

// GetArticleList provides a mock function with given fields: search
func (_m *ArticleService) GetArticleList(search string) ([]domain.ArticleList, error) {
	ret := _m.Called(search)

	if len(ret) == 0 {
		panic("no return value specified for GetArticleList")
	}

	var r0 []domain.ArticleList
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]domain.ArticleList, error)); ok {
		return rf(search)
	}
	if rf, ok := ret.Get(0).(func(string) []domain.ArticleList); ok {
		r0 = rf(search)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.ArticleList)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(search)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewArticleService creates a new instance of ArticleService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewArticleService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ArticleService {
	mock := &ArticleService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
