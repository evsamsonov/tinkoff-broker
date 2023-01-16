// Code generated by mockery v2.13.0-beta.1. DO NOT EDIT.

package tnkbroker

import (
	context "context"

	grpc "google.golang.org/grpc"

	investapi "github.com/tinkoff/invest-api-go-sdk"

	mock "github.com/stretchr/testify/mock"
)

// mockInstrumentServiceClient is an autogenerated mock type for the instrumentServiceClient type
type mockInstrumentServiceClient struct {
	mock.Mock
}

// BondBy provides a mock function with given fields: ctx, in, opts
func (_m *mockInstrumentServiceClient) BondBy(ctx context.Context, in *investapi.InstrumentRequest, opts ...grpc.CallOption) (*investapi.BondResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.BondResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.InstrumentRequest, ...grpc.CallOption) *investapi.BondResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.BondResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.InstrumentRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Bonds provides a mock function with given fields: ctx, in, opts
func (_m *mockInstrumentServiceClient) Bonds(ctx context.Context, in *investapi.InstrumentsRequest, opts ...grpc.CallOption) (*investapi.BondsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.BondsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.InstrumentsRequest, ...grpc.CallOption) *investapi.BondsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.BondsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.InstrumentsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Currencies provides a mock function with given fields: ctx, in, opts
func (_m *mockInstrumentServiceClient) Currencies(ctx context.Context, in *investapi.InstrumentsRequest, opts ...grpc.CallOption) (*investapi.CurrenciesResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.CurrenciesResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.InstrumentsRequest, ...grpc.CallOption) *investapi.CurrenciesResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.CurrenciesResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.InstrumentsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CurrencyBy provides a mock function with given fields: ctx, in, opts
func (_m *mockInstrumentServiceClient) CurrencyBy(ctx context.Context, in *investapi.InstrumentRequest, opts ...grpc.CallOption) (*investapi.CurrencyResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.CurrencyResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.InstrumentRequest, ...grpc.CallOption) *investapi.CurrencyResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.CurrencyResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.InstrumentRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EtfBy provides a mock function with given fields: ctx, in, opts
func (_m *mockInstrumentServiceClient) EtfBy(ctx context.Context, in *investapi.InstrumentRequest, opts ...grpc.CallOption) (*investapi.EtfResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.EtfResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.InstrumentRequest, ...grpc.CallOption) *investapi.EtfResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.EtfResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.InstrumentRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Etfs provides a mock function with given fields: ctx, in, opts
func (_m *mockInstrumentServiceClient) Etfs(ctx context.Context, in *investapi.InstrumentsRequest, opts ...grpc.CallOption) (*investapi.EtfsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.EtfsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.InstrumentsRequest, ...grpc.CallOption) *investapi.EtfsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.EtfsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.InstrumentsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FutureBy provides a mock function with given fields: ctx, in, opts
func (_m *mockInstrumentServiceClient) FutureBy(ctx context.Context, in *investapi.InstrumentRequest, opts ...grpc.CallOption) (*investapi.FutureResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.FutureResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.InstrumentRequest, ...grpc.CallOption) *investapi.FutureResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.FutureResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.InstrumentRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Futures provides a mock function with given fields: ctx, in, opts
func (_m *mockInstrumentServiceClient) Futures(ctx context.Context, in *investapi.InstrumentsRequest, opts ...grpc.CallOption) (*investapi.FuturesResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.FuturesResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.InstrumentsRequest, ...grpc.CallOption) *investapi.FuturesResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.FuturesResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.InstrumentsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccruedInterests provides a mock function with given fields: ctx, in, opts
func (_m *mockInstrumentServiceClient) GetAccruedInterests(ctx context.Context, in *investapi.GetAccruedInterestsRequest, opts ...grpc.CallOption) (*investapi.GetAccruedInterestsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.GetAccruedInterestsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.GetAccruedInterestsRequest, ...grpc.CallOption) *investapi.GetAccruedInterestsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.GetAccruedInterestsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.GetAccruedInterestsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAssetBy provides a mock function with given fields: ctx, in, opts
func (_m *mockInstrumentServiceClient) GetAssetBy(ctx context.Context, in *investapi.AssetRequest, opts ...grpc.CallOption) (*investapi.AssetResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.AssetResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.AssetRequest, ...grpc.CallOption) *investapi.AssetResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.AssetResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.AssetRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAssets provides a mock function with given fields: ctx, in, opts
func (_m *mockInstrumentServiceClient) GetAssets(ctx context.Context, in *investapi.AssetsRequest, opts ...grpc.CallOption) (*investapi.AssetsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.AssetsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.AssetsRequest, ...grpc.CallOption) *investapi.AssetsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.AssetsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.AssetsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBondCoupons provides a mock function with given fields: ctx, in, opts
func (_m *mockInstrumentServiceClient) GetBondCoupons(ctx context.Context, in *investapi.GetBondCouponsRequest, opts ...grpc.CallOption) (*investapi.GetBondCouponsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.GetBondCouponsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.GetBondCouponsRequest, ...grpc.CallOption) *investapi.GetBondCouponsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.GetBondCouponsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.GetBondCouponsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDividends provides a mock function with given fields: ctx, in, opts
func (_m *mockInstrumentServiceClient) GetDividends(ctx context.Context, in *investapi.GetDividendsRequest, opts ...grpc.CallOption) (*investapi.GetDividendsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.GetDividendsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.GetDividendsRequest, ...grpc.CallOption) *investapi.GetDividendsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.GetDividendsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.GetDividendsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFuturesMargin provides a mock function with given fields: ctx, in, opts
func (_m *mockInstrumentServiceClient) GetFuturesMargin(ctx context.Context, in *investapi.GetFuturesMarginRequest, opts ...grpc.CallOption) (*investapi.GetFuturesMarginResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.GetFuturesMarginResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.GetFuturesMarginRequest, ...grpc.CallOption) *investapi.GetFuturesMarginResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.GetFuturesMarginResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.GetFuturesMarginRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetInstrumentBy provides a mock function with given fields: ctx, in, opts
func (_m *mockInstrumentServiceClient) GetInstrumentBy(ctx context.Context, in *investapi.InstrumentRequest, opts ...grpc.CallOption) (*investapi.InstrumentResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.InstrumentResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.InstrumentRequest, ...grpc.CallOption) *investapi.InstrumentResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.InstrumentResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.InstrumentRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ShareBy provides a mock function with given fields: ctx, in, opts
func (_m *mockInstrumentServiceClient) ShareBy(ctx context.Context, in *investapi.InstrumentRequest, opts ...grpc.CallOption) (*investapi.ShareResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.ShareResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.InstrumentRequest, ...grpc.CallOption) *investapi.ShareResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.ShareResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.InstrumentRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Shares provides a mock function with given fields: ctx, in, opts
func (_m *mockInstrumentServiceClient) Shares(ctx context.Context, in *investapi.InstrumentsRequest, opts ...grpc.CallOption) (*investapi.SharesResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.SharesResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.InstrumentsRequest, ...grpc.CallOption) *investapi.SharesResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.SharesResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.InstrumentsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TradingSchedules provides a mock function with given fields: ctx, in, opts
func (_m *mockInstrumentServiceClient) TradingSchedules(ctx context.Context, in *investapi.TradingSchedulesRequest, opts ...grpc.CallOption) (*investapi.TradingSchedulesResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *investapi.TradingSchedulesResponse
	if rf, ok := ret.Get(0).(func(context.Context, *investapi.TradingSchedulesRequest, ...grpc.CallOption) *investapi.TradingSchedulesResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*investapi.TradingSchedulesResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *investapi.TradingSchedulesRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type newMockInstrumentServiceClientT interface {
	mock.TestingT
	Cleanup(func())
}

// newMockInstrumentServiceClient creates a new instance of mockInstrumentServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockInstrumentServiceClient(t newMockInstrumentServiceClientT) *mockInstrumentServiceClient {
	mock := &mockInstrumentServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}