// Code generated by mockery v2.13.0-beta.1. DO NOT EDIT.

package tnkbroker

import (
	context "context"

	grpc "google.golang.org/grpc"

	investapi "github.com/tinkoff/invest-api-go-sdk"

	mock "github.com/stretchr/testify/mock"
)

// mockMarketDataServiceClient is an autogenerated mock type for the marketDataServiceClient type
type mockMarketDataServiceClient struct {
	mock.Mock
}

// GetCandles provides a mock function with given fields: ctx, in, opts
func (_m *mockMarketDataServiceClient) GetCandles(ctx context.Context, in *investapi.GetCandlesRequest, opts ...grpc.CallOption) (*investapi.GetCandlesResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.GetCandlesResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.GetCandlesRequest, ...grpc.CallOption) *investapi.GetCandlesResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.GetCandlesResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.GetCandlesRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastPrices provides a mock function with given fields: ctx, in, opts
func (_m *mockMarketDataServiceClient) GetLastPrices(ctx context.Context, in *investapi.GetLastPricesRequest, opts ...grpc.CallOption) (*investapi.GetLastPricesResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.GetLastPricesResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.GetLastPricesRequest, ...grpc.CallOption) *investapi.GetLastPricesResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.GetLastPricesResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.GetLastPricesRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastTrades provides a mock function with given fields: ctx, in, opts
func (_m *mockMarketDataServiceClient) GetLastTrades(ctx context.Context, in *investapi.GetLastTradesRequest, opts ...grpc.CallOption) (*investapi.GetLastTradesResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.GetLastTradesResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.GetLastTradesRequest, ...grpc.CallOption) *investapi.GetLastTradesResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.GetLastTradesResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.GetLastTradesRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrderBook provides a mock function with given fields: ctx, in, opts
func (_m *mockMarketDataServiceClient) GetOrderBook(ctx context.Context, in *investapi.GetOrderBookRequest, opts ...grpc.CallOption) (*investapi.GetOrderBookResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.GetOrderBookResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.GetOrderBookRequest, ...grpc.CallOption) *investapi.GetOrderBookResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.GetOrderBookResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.GetOrderBookRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTradingStatus provides a mock function with given fields: ctx, in, opts
func (_m *mockMarketDataServiceClient) GetTradingStatus(ctx context.Context, in *investapi.GetTradingStatusRequest, opts ...grpc.CallOption) (*investapi.GetTradingStatusResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.GetTradingStatusResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.GetTradingStatusRequest, ...grpc.CallOption) *investapi.GetTradingStatusResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.GetTradingStatusResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.GetTradingStatusRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type newMockMarketDataServiceClientT interface {
	mock.TestingT
	Cleanup(func())
}

// newMockMarketDataServiceClient creates a new instance of mockMarketDataServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockMarketDataServiceClient(t newMockMarketDataServiceClientT) *mockMarketDataServiceClient {
	mock := &mockMarketDataServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}