// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import context "context"

import huandufacebook "github.com/huandu/facebook"
import mock "github.com/stretchr/testify/mock"
import oauth2 "golang.org/x/oauth2"

// FacebookClient is an autogenerated mock type for the FacebookClient type
type FacebookClient struct {
	mock.Mock
}

// Exchange provides a mock function with given fields: config, ctx, code
func (_m *FacebookClient) Exchange(config *oauth2.Config, ctx context.Context, code string) (*oauth2.Token, error) {
	ret := _m.Called(config, ctx, code)

	var r0 *oauth2.Token
	if rf, ok := ret.Get(0).(func(*oauth2.Config, context.Context, string) *oauth2.Token); ok {
		r0 = rf(config, ctx, code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth2.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*oauth2.Config, context.Context, string) error); ok {
		r1 = rf(config, ctx, code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetInfo provides a mock function with given fields: accessToken
func (_m *FacebookClient) GetInfo(accessToken string) (huandufacebook.Result, error) {
	ret := _m.Called(accessToken)

	var r0 huandufacebook.Result
	if rf, ok := ret.Get(0).(func(string) huandufacebook.Result); ok {
		r0 = rf(accessToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(huandufacebook.Result)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(accessToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}