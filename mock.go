package tnkbroker

import investapi "github.com/russianinvestments/invest-api-go-sdk/investgo"

// nolint: lll,unused
//
//go:generate docker run --rm -v ${PWD}:/app -w /app vektra/mockery --name ordersServiceClient --inpackage --case snake
type ordersServiceClient interface {
	investapi.OrdersServiceClient
}

// nolint: lll,unused
//
//go:generate docker run --rm -v ${PWD}:/app -w /app vektra/mockery --name stopOrdersServiceClient --inpackage --case snake
type stopOrdersServiceClient interface {
	investapi.StopOrdersServiceClient
}

// nolint: lll,unused
//
//go:generate docker run --rm -v ${PWD}:/app -w /app vektra/mockery --name marketDataServiceClient --inpackage --case snake
type marketDataServiceClient interface {
	investapi.MarketDataServiceClient
}

// nolint: lll,unused
//
//go:generate docker run --rm -v ${PWD}:/app -w /app vektra/mockery --name instrumentServiceClient --inpackage --case snake
type instrumentServiceClient interface {
	investapi.InstrumentsServiceClient
}
