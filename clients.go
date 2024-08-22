package tnkbroker

import (
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

// nolint: lll
//
//go:generate docker run --rm -v ${PWD}:/app -w /app vektra/mockery --name ordersServiceClient --inpackage --case snake
type ordersServiceClient interface {
	PostOrder(req *investgo.PostOrderRequest) (*investgo.PostOrderResponse, error)
	GetOrderState(accountID, orderID string, priceType pb.PriceType) (*investgo.GetOrderStateResponse, error)
}

// nolint: lll
//
//go:generate docker run --rm -v ${PWD}:/app -w /app vektra/mockery --name stopOrdersServiceClient --inpackage --case snake
type stopOrdersServiceClient interface {
	PostStopOrder(req *investgo.PostStopOrderRequest) (*investgo.PostStopOrderResponse, error)
	CancelStopOrder(accountID, stopOrderID string) (*investgo.CancelStopOrderResponse, error)
	GetStopOrders(accountID string) (*investgo.GetStopOrdersResponse, error)
}

// nolint: lll
//
//go:generate docker run --rm -v ${PWD}:/app -w /app vektra/mockery --name ordersStreamClient --inpackage --case snake
type ordersStreamClient interface {
	TradesStream(accounts []string) (*investgo.TradesStream, error)
}

// nolint: lll
//
//go:generate docker run --rm -v ${PWD}:/app -w /app vektra/mockery --name marketDataServiceClient --inpackage --case snake
type marketDataServiceClient interface {
	GetLastPrices(instrumentIds []string) (*investgo.GetLastPricesResponse, error)
}

// nolint: lll
//
//go:generate docker run --rm -v ${PWD}:/app -w /app vektra/mockery --name instrumentsServiceClient --inpackage --case snake
type instrumentsServiceClient interface {
	InstrumentByFigi(id string) (*investgo.InstrumentResponse, error)
}
