// Code generated by mockery v2.43.2. DO NOT EDIT.

package tnkbroker

import (
	investgo "github.com/russianinvestments/invest-api-go-sdk/investgo"
	mock "github.com/stretchr/testify/mock"
)

// mockInstrumentsServiceClient is an autogenerated mock type for the instrumentsServiceClient type
type mockInstrumentsServiceClient struct {
	mock.Mock
}

// InstrumentByFigi provides a mock function with given fields: id
func (_m *mockInstrumentsServiceClient) InstrumentByFigi(id string) (*investgo.InstrumentResponse, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for InstrumentByFigi")
	}

	var r0 *investgo.InstrumentResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*investgo.InstrumentResponse, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *investgo.InstrumentResponse); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investgo.InstrumentResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// newMockInstrumentsServiceClient creates a new instance of mockInstrumentsServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockInstrumentsServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockInstrumentsServiceClient {
	mock := &mockInstrumentsServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
