// Code generated by mockery v2.43.2. DO NOT EDIT.

package tnkbroker

import (
	investgo "github.com/russianinvestments/invest-api-go-sdk/investgo"
	mock "github.com/stretchr/testify/mock"
)

// mockStopOrdersServiceClient is an autogenerated mock type for the stopOrdersServiceClient type
type mockStopOrdersServiceClient struct {
	mock.Mock
}

// CancelStopOrder provides a mock function with given fields: accountId, stopOrderId
func (_m *mockStopOrdersServiceClient) CancelStopOrder(accountId string, stopOrderId string) (*investgo.CancelStopOrderResponse, error) {
	ret := _m.Called(accountId, stopOrderId)

	if len(ret) == 0 {
		panic("no return value specified for CancelStopOrder")
	}

	var r0 *investgo.CancelStopOrderResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*investgo.CancelStopOrderResponse, error)); ok {
		return rf(accountId, stopOrderId)
	}
	if rf, ok := ret.Get(0).(func(string, string) *investgo.CancelStopOrderResponse); ok {
		r0 = rf(accountId, stopOrderId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investgo.CancelStopOrderResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(accountId, stopOrderId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStopOrders provides a mock function with given fields: accountId
func (_m *mockStopOrdersServiceClient) GetStopOrders(accountId string) (*investgo.GetStopOrdersResponse, error) {
	ret := _m.Called(accountId)

	if len(ret) == 0 {
		panic("no return value specified for GetStopOrders")
	}

	var r0 *investgo.GetStopOrdersResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*investgo.GetStopOrdersResponse, error)); ok {
		return rf(accountId)
	}
	if rf, ok := ret.Get(0).(func(string) *investgo.GetStopOrdersResponse); ok {
		r0 = rf(accountId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investgo.GetStopOrdersResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(accountId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PostStopOrder provides a mock function with given fields: req
func (_m *mockStopOrdersServiceClient) PostStopOrder(req *investgo.PostStopOrderRequest) (*investgo.PostStopOrderResponse, error) {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for PostStopOrder")
	}

	var r0 *investgo.PostStopOrderResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*investgo.PostStopOrderRequest) (*investgo.PostStopOrderResponse, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*investgo.PostStopOrderRequest) *investgo.PostStopOrderResponse); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investgo.PostStopOrderResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*investgo.PostStopOrderRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// newMockStopOrdersServiceClient creates a new instance of mockStopOrdersServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockStopOrdersServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockStopOrdersServiceClient {
	mock := &mockStopOrdersServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
