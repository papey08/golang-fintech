// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import ads "homework_extra/internal/model/ads"

import context "context"
import filter "homework_extra/internal/model/filter"
import mock "github.com/stretchr/testify/mock"
import users "homework_extra/internal/model/users"

// App is an autogenerated mock type for the App type
type App struct {
	mock.Mock
}

// ChangeAdStatus provides a mock function with given fields: ctx, adID, userID, published
func (_m *App) ChangeAdStatus(ctx context.Context, adID int64, userID int64, published bool) (ads.Ad, error) {
	ret := _m.Called(ctx, adID, userID, published)

	var r0 ads.Ad
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, bool) ads.Ad); ok {
		r0 = rf(ctx, adID, userID, published)
	} else {
		r0 = ret.Get(0).(ads.Ad)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, int64, bool) error); ok {
		r1 = rf(ctx, adID, userID, published)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateAd provides a mock function with given fields: ctx, title, text, userID
func (_m *App) CreateAd(ctx context.Context, title string, text string, userID int64) (ads.Ad, error) {
	ret := _m.Called(ctx, title, text, userID)

	var r0 ads.Ad
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int64) ads.Ad); ok {
		r0 = rf(ctx, title, text, userID)
	} else {
		r0 = ret.Get(0).(ads.Ad)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, int64) error); ok {
		r1 = rf(ctx, title, text, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: ctx, nickname, email
func (_m *App) CreateUser(ctx context.Context, nickname string, email string) (users.User, error) {
	ret := _m.Called(ctx, nickname, email)

	var r0 users.User
	if rf, ok := ret.Get(0).(func(context.Context, string, string) users.User); ok {
		r0 = rf(ctx, nickname, email)
	} else {
		r0 = ret.Get(0).(users.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, nickname, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAd provides a mock function with given fields: ctx, userID, adID
func (_m *App) DeleteAd(ctx context.Context, userID int64, adID int64) error {
	ret := _m.Called(ctx, userID, adID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) error); ok {
		r0 = rf(ctx, userID, adID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteUser provides a mock function with given fields: ctx, userID
func (_m *App) DeleteUser(ctx context.Context, userID int64) error {
	ret := _m.Called(ctx, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAdByID provides a mock function with given fields: ctx, adID
func (_m *App) GetAdByID(ctx context.Context, adID int64) (ads.Ad, error) {
	ret := _m.Called(ctx, adID)

	var r0 ads.Ad
	if rf, ok := ret.Get(0).(func(context.Context, int64) ads.Ad); ok {
		r0 = rf(ctx, adID)
	} else {
		r0 = ret.Get(0).(ads.Ad)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, adID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAdsList provides a mock function with given fields: ctx, f
func (_m *App) GetAdsList(ctx context.Context, f filter.Filter) ([]ads.Ad, error) {
	ret := _m.Called(ctx, f)

	var r0 []ads.Ad
	if rf, ok := ret.Get(0).(func(context.Context, filter.Filter) []ads.Ad); ok {
		r0 = rf(ctx, f)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ads.Ad)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, filter.Filter) error); ok {
		r1 = rf(ctx, f)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByID provides a mock function with given fields: ctx, userID
func (_m *App) GetUserByID(ctx context.Context, userID int64) (users.User, error) {
	ret := _m.Called(ctx, userID)

	var r0 users.User
	if rf, ok := ret.Get(0).(func(context.Context, int64) users.User); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(users.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchAds provides a mock function with given fields: ctx, pattern
func (_m *App) SearchAds(ctx context.Context, pattern string) ([]ads.Ad, error) {
	ret := _m.Called(ctx, pattern)

	var r0 []ads.Ad
	if rf, ok := ret.Get(0).(func(context.Context, string) []ads.Ad); ok {
		r0 = rf(ctx, pattern)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ads.Ad)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, pattern)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAd provides a mock function with given fields: ctx, adID, userID, title, text
func (_m *App) UpdateAd(ctx context.Context, adID int64, userID int64, title string, text string) (ads.Ad, error) {
	ret := _m.Called(ctx, adID, userID, title, text)

	var r0 ads.Ad
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, string, string) ads.Ad); ok {
		r0 = rf(ctx, adID, userID, title, text)
	} else {
		r0 = ret.Get(0).(ads.Ad)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, int64, string, string) error); ok {
		r1 = rf(ctx, adID, userID, title, text)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: ctx, userID, nickname, email
func (_m *App) UpdateUser(ctx context.Context, userID int64, nickname string, email string) (users.User, error) {
	ret := _m.Called(ctx, userID, nickname, email)

	var r0 users.User
	if rf, ok := ret.Get(0).(func(context.Context, int64, string, string) users.User); ok {
		r0 = rf(ctx, userID, nickname, email)
	} else {
		r0 = ret.Get(0).(users.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, string, string) error); ok {
		r1 = rf(ctx, userID, nickname, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}