package tnkbroker

import investapi "github.com/tinkoff/invest-api-go-sdk"

// nolint: lll,unused
//
//go:generate docker run --rm -v ${PWD}/../../:/app -w /app vektra/mockery --dir=/app/broker/tinkoff --name ordersServiceClient --inpackage --case snake
type ordersServiceClient interface {
	investapi.OrdersServiceClient
}

// nolint: lll,unused
//
//go:generate docker run --rm -v ${PWD}/../../:/app -w /app vektra/mockery --dir=/app/broker/tinkoff --name stopOrdersServiceClient --inpackage --case snake
type stopOrdersServiceClient interface {
	investapi.StopOrdersServiceClient
}
