// Code generated by mockery v2.43.2. DO NOT EDIT.

package tnkbroker

import mock "github.com/stretchr/testify/mock"

// mockOrdersServiceClient is an autogenerated mock type for the ordersServiceClient type
type mockOrdersServiceClient struct {
	mock.Mock
}

// newMockOrdersServiceClient creates a new instance of mockOrdersServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockOrdersServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockOrdersServiceClient {
	mock := &mockOrdersServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
