// Package tnkbroker implements [trengin.Broker] using [Tinkoff Invest API].
// The implementation doesn't support multiple open positions at the same time.
// Commission in brokerPosition is approximate.
//
// [Tinkoff Invest API]: https://tinkoff.github.io/investAPI/
package tnkbroker

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/evsamsonov/trengin"
	"github.com/google/uuid"
	investapi "github.com/tinkoff/invest-api-go-sdk"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/evsamsonov/tinkoff-broker/internal/tnkposition"
)

var _ trengin.Broker = &Tinkoff{}

const (
	tinkoffHost                    = "invest-public-api.tinkoff.ru:443"
	defaultProtectiveSpread        = 1 // In percent
	defaultTradeStreamRetryTimeout = 1 * time.Minute
	defaultTradeStreamPingTimeout  = 6 * time.Minute
)

type Tinkoff struct {
	accountID                   string
	token                       string
	appName                     string
	orderClient                 investapi.OrdersServiceClient
	stopOrderClient             investapi.StopOrdersServiceClient
	tradeStreamClient           investapi.OrdersStreamServiceClient
	marketDataClient            investapi.MarketDataServiceClient
	instrumentClient            investapi.InstrumentsServiceClient
	tradeStreamRetryTimeout     time.Duration
	tradeStreamPingWaitDuration time.Duration
	protectiveSpread            float64
	positionStorage             *tnkposition.Storage
	logger                      *zap.Logger
}

type Option func(*Tinkoff)

// WithLogger returns Option which sets logger. The default logger is no-op Logger
func WithLogger(logger *zap.Logger) Option {
	return func(t *Tinkoff) {
		t.logger = logger
	}
}

// WithAppName returns Option which sets [x-app-name]
//
// [x-app-name]: https://tinkoff.github.io/investAPI/grpc/#appname
func WithAppName(appName string) Option {
	return func(t *Tinkoff) {
		t.appName = appName
	}
}

// WithProtectiveSpread returns Option which sets protective spread
// in percent for executing orders. The default value is 1%
func WithProtectiveSpread(protectiveSpread float64) Option {
	return func(t *Tinkoff) {
		t.protectiveSpread = protectiveSpread
	}
}

// WithTradeStreamRetryTimeout returns Option which defines retry timeout
// on trade stream error
func WithTradeStreamRetryTimeout(timeout time.Duration) Option {
	return func(t *Tinkoff) {
		t.tradeStreamRetryTimeout = timeout
	}
}

// WithTradeStreamPingWaitDuration returns Option which defines duration
// how long we wait for ping before reconnection
func WithTradeStreamPingWaitDuration(duration time.Duration) Option {
	return func(t *Tinkoff) {
		t.tradeStreamPingWaitDuration = duration
	}
}

// New creates a new Tinkoff object. It takes [full-access token],
// user account identifier.
//
// [full-access token]: https://tinkoff.github.io/investAPI/token/
func New(token, accountID string, opts ...Option) (*Tinkoff, error) {
	conn, err := grpc.Dial(
		tinkoffHost,
		grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true, //nolint: gosec
			}),
		),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc dial: %w", err)
	}

	tinkoff := &Tinkoff{
		accountID:                   accountID,
		token:                       token,
		protectiveSpread:            defaultProtectiveSpread,
		orderClient:                 investapi.NewOrdersServiceClient(conn),
		stopOrderClient:             investapi.NewStopOrdersServiceClient(conn),
		tradeStreamClient:           investapi.NewOrdersStreamServiceClient(conn),
		marketDataClient:            investapi.NewMarketDataServiceClient(conn),
		instrumentClient:            investapi.NewInstrumentsServiceClient(conn),
		tradeStreamRetryTimeout:     defaultTradeStreamRetryTimeout,
		tradeStreamPingWaitDuration: defaultTradeStreamPingTimeout,
		positionStorage:             tnkposition.NewStorage(),
		logger:                      zap.NewNop(),
	}
	for _, opt := range opts {
		opt(tinkoff)
	}
	return tinkoff, nil
}

// Run starts to track an open positions
func (t *Tinkoff) Run(ctx context.Context) error {
	for {
		err := t.readTradesStream(ctx)
		if err == nil {
			return nil
		}
		t.logger.Error(
			"Failed to read trade stream. Retry after timeout",
			zap.Error(err),
			zap.Duration("timeout", t.tradeStreamRetryTimeout),
		)

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(t.tradeStreamRetryTimeout):
		}
	}
}

// OpenPosition opens brokerPosition, returns new brokerPosition and channel for tracking brokerPosition closing.
func (t *Tinkoff) OpenPosition(
	ctx context.Context,
	action trengin.OpenPositionAction,
) (trengin.Position, trengin.PositionClosed, error) {
	ctx = t.ctxWithMetadata(ctx)
	instrument, err := t.getInstrument(ctx, action.FIGI)
	if err != nil {
		return trengin.Position{}, nil, fmt.Errorf("get instrument: %w", err)
	}
	openPrice, commission, err := t.openMarketOrder(ctx, instrument, action.Type, action.Quantity)
	if err != nil {
		return trengin.Position{}, nil, fmt.Errorf("open market order: %w", err)
	}

	position, err := trengin.NewPosition(action, time.Now(), openPrice.ToFloat())
	if err != nil {
		return trengin.Position{}, nil, fmt.Errorf("new position: %w", err)
	}
	position.AddCommission(commission.ToFloat())

	var stopLossID, takeProfitID string
	if action.StopLossIndent != 0 {
		stopLoss := t.stopLossPriceByOpen(openPrice, action, instrument.MinPriceIncrement)
		stopLossID, err = t.setStopLoss(ctx, instrument, stopLoss, *position)
		if err != nil {
			return trengin.Position{}, nil, fmt.Errorf("set stop order: %w", err)
		}
	}
	if action.TakeProfitIndent != 0 {
		takeProfit := t.takeProfitPriceByOpen(openPrice, action, instrument.MinPriceIncrement)
		takeProfitID, err = t.setTakeProfit(ctx, instrument, takeProfit, *position)
		if err != nil {
			return trengin.Position{}, nil, fmt.Errorf("set stop order: %w", err)
		}
	}

	positionClosed := make(chan trengin.Position, 1)
	t.positionStorage.Store(tnkposition.NewPosition(position, instrument, stopLossID, takeProfitID, positionClosed))

	return *position, positionClosed, nil
}

// ChangeConditionalOrder changes stop loss and take profit of current brokerPosition.
// It returns updated brokerPosition.
func (t *Tinkoff) ChangeConditionalOrder(
	ctx context.Context,
	action trengin.ChangeConditionalOrderAction,
) (trengin.Position, error) {
	tinkoffPosition, unlockPosition, err := t.positionStorage.Load(action.PositionID)
	if err != nil {
		return trengin.Position{}, fmt.Errorf("load position: %w", err)
	}
	defer unlockPosition()

	ctx = t.ctxWithMetadata(ctx)
	instrument := tinkoffPosition.Instrument()
	if action.StopLoss != 0 {
		if err := t.cancelStopOrder(ctx, tinkoffPosition.StopLossID()); err != nil {
			return trengin.Position{}, fmt.Errorf("cancel stop loss: %w", err)
		}

		stopLoss := t.convertFloatToQuotation(instrument.MinPriceIncrement, action.StopLoss)
		stopLossID, err := t.setStopLoss(ctx, instrument, stopLoss, tinkoffPosition.Position())
		if err != nil {
			return trengin.Position{}, fmt.Errorf("set stop loss: %w", err)
		}
		tinkoffPosition.SetStopLoss(stopLossID, action.StopLoss)
	}

	if action.TakeProfit != 0 {
		if err := t.cancelStopOrder(ctx, tinkoffPosition.TakeProfitID()); err != nil {
			return trengin.Position{}, fmt.Errorf("cancel take profit: %w", err)
		}

		takeProfit := t.convertFloatToQuotation(instrument.MinPriceIncrement, action.TakeProfit)
		takeProfitID, err := t.setTakeProfit(ctx, instrument, takeProfit, tinkoffPosition.Position())
		if err != nil {
			return trengin.Position{}, fmt.Errorf("set take profit: %w", err)
		}
		tinkoffPosition.SetTakeProfitID(takeProfitID, action.TakeProfit)
	}
	return tinkoffPosition.Position(), nil
}

// ClosePosition closes current position and returns closed position.
func (t *Tinkoff) ClosePosition(ctx context.Context, action trengin.ClosePositionAction) (trengin.Position, error) {
	tinkoffPosition, unlockPosition, err := t.positionStorage.Load(action.PositionID)
	if err != nil {
		return trengin.Position{}, fmt.Errorf("load position: %w", err)
	}
	defer unlockPosition()

	ctx = t.ctxWithMetadata(ctx)
	if err := t.cancelStopOrder(ctx, tinkoffPosition.StopLossID()); err != nil {
		return trengin.Position{}, fmt.Errorf("cancel stop loss: %w", err)
	}
	if err := t.cancelStopOrder(ctx, tinkoffPosition.TakeProfitID()); err != nil {
		return trengin.Position{}, fmt.Errorf("cancel take profit: %w", err)
	}

	instrument := tinkoffPosition.Instrument()
	position := tinkoffPosition.Position()
	closePrice, commission, err := t.openMarketOrder(ctx, instrument, position.Type.Inverse(), position.Quantity)
	if err != nil {
		return trengin.Position{}, fmt.Errorf("open market order: %w", err)
	}

	position.AddCommission(commission.ToFloat())
	if err := tinkoffPosition.Close(closePrice.ToFloat()); err != nil {
		return trengin.Position{}, fmt.Errorf("close: %w", err)
	}

	t.logger.Info("Position was closed", zap.Any("tinkoffPosition", tinkoffPosition))
	return tinkoffPosition.Position(), nil
}

func (t *Tinkoff) readTradesStream(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctx = t.ctxWithMetadata(ctx)
	stream, err := t.tradeStreamClient.TradesStream(ctx, &investapi.TradesStreamRequest{})
	if err != nil {
		return fmt.Errorf("trades stream: %w", err)
	}

	g := errgroup.Group{}
	heartbeat := make(chan struct{})
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-heartbeat:
				continue
			case <-time.After(t.tradeStreamPingWaitDuration):
				cancel()
				return fmt.Errorf("ping timed out")
			}
		}
	})

	g.Go(func() error {
		t.logger.Debug("Trades stream is connected")
		for {
			resp, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					t.logger.Info("Trades stream connection is closed")
					break
				}
				if status.Code(err) == codes.Canceled {
					t.logger.Info("Trades stream connection is canceled")
					break
				}
				return fmt.Errorf("stream recv: %w", err)
			}

			switch v := resp.Payload.(type) {
			case *investapi.TradesStreamResponse_Ping:
				t.logger.Debug("Trade stream ping was received", zap.Any("ping", v))

				select {
				case heartbeat <- struct{}{}:
				default:
				}
			case *investapi.TradesStreamResponse_OrderTrades:
				t.logger.Info("Order trades were received", zap.Any("orderTrades", v))

				if err := t.processOrderTrades(ctx, v.OrderTrades); err != nil {
					return fmt.Errorf("process order trades: %w", err)
				}
			default:
				return errors.New("unexpected payload")
			}
		}
		return nil
	})
	return g.Wait()
}

func (t *Tinkoff) processOrderTrades(ctx context.Context, orderTrades *investapi.OrderTrades) error {
	if orderTrades.AccountId != t.accountID {
		return nil
	}
	err := t.positionStorage.ForEach(func(tinkoffPosition *tnkposition.Position) error {
		position := tinkoffPosition.Position()
		if orderTrades.Figi != position.FIGI {
			return nil
		}
		longClosed := position.IsLong() && orderTrades.Direction == investapi.OrderDirection_ORDER_DIRECTION_SELL
		shortClosed := position.IsShort() && orderTrades.Direction == investapi.OrderDirection_ORDER_DIRECTION_BUY
		if !longClosed && !shortClosed {
			return nil
		}

		conditionalOrdersFound, err := t.conditionalOrdersFound(ctx, tinkoffPosition)
		if err != nil {
			return fmt.Errorf("conditional orders found: %w", err)
		}
		if conditionalOrdersFound {
			return nil
		}

		logger := t.logger.With(zap.Any("position", position))
		tinkoffPosition.AddOrderTrade(orderTrades.GetTrades()...)

		var executedQuantity int64
		for _, trade := range tinkoffPosition.OrderTrades() {
			executedQuantity += trade.GetQuantity() / int64(tinkoffPosition.Instrument().Lot)
		}
		if executedQuantity < position.Quantity {
			logger.Info("Position partially closed", zap.Any("executedQuantity", executedQuantity))
			return nil
		}

		if err := t.cancelStopOrders(ctx, tinkoffPosition); err != nil {
			return fmt.Errorf("cancel stop orders: %w", err)
		}

		orderState, err := t.getExecutedOrderState(ctx, orderTrades.OrderId)
		if err != nil {
			return fmt.Errorf("get executed order state: %w", err)
		}

		closePrice := NewMoneyValue(orderState.AveragePositionPrice)
		commission := NewMoneyValue(orderState.InitialCommission)
		tinkoffPosition.AddCommission(commission.ToFloat())

		if err := tinkoffPosition.Close(closePrice.ToFloat()); err != nil {
			if errors.Is(err, trengin.ErrAlreadyClosed) {
				logger.Info("Position already closed")
				return nil
			}
			return fmt.Errorf("close: %w", err)
		}
		logger.Info("Position was closed by order trades", zap.Any("orderTrades", orderTrades))
		return nil
	})
	return err
}

func (t *Tinkoff) ctxWithMetadata(ctx context.Context) context.Context {
	md := metadata.New(map[string]string{
		"Authorization": "Bearer " + t.token,
		"x-app-name":    t.appName,
	})
	return metadata.NewOutgoingContext(ctx, md)
}

func (t *Tinkoff) openMarketOrder(
	ctx context.Context,
	instrument *investapi.Instrument,
	positionType trengin.PositionType,
	quantity int64,
) (*MoneyValue, *MoneyValue, error) {
	direction := investapi.OrderDirection_ORDER_DIRECTION_BUY
	if positionType.IsShort() {
		direction = investapi.OrderDirection_ORDER_DIRECTION_SELL
	}
	price, err := t.getLastPrice(ctx, instrument.Figi)
	if err != nil {
		return nil, nil, fmt.Errorf("get last price: %w", err)
	}
	orderRequest := &investapi.PostOrderRequest{
		Figi:      instrument.Figi,
		Quantity:  quantity,
		Direction: direction,
		AccountId: t.accountID,
		Price:     t.addProtectedSpread(positionType.Inverse(), price, instrument.MinPriceIncrement),
		OrderType: investapi.OrderType_ORDER_TYPE_LIMIT,
		OrderId:   uuid.New().String(),
	}

	order, err := t.orderClient.PostOrder(ctx, orderRequest)
	if err != nil {
		t.logger.Error("Failed to execute order", zap.Error(err), zap.Any("orderRequest", orderRequest))
		return nil, nil, fmt.Errorf("post order: %w", err)
	}

	orderState, err := t.getExecutedOrderState(ctx, order.OrderId)
	if err != nil {
		return nil, nil, fmt.Errorf("wait fill order state: %w", err)
	}
	t.logger.Info("Order was executed",
		zap.Any("orderRequest", orderRequest),
		zap.Any("order", order),
		zap.Any("orderState", orderState),
	)
	return NewMoneyValue(orderState.AveragePositionPrice), NewMoneyValue(orderState.InitialCommission), nil
}

type stopOrderType int

const (
	stopLossStopOrderType stopOrderType = iota + 1
	takeProfitStopOrderType
)

func (t *Tinkoff) setStopLoss(
	ctx context.Context,
	instrument *investapi.Instrument,
	price *investapi.Quotation,
	position trengin.Position,
) (string, error) {
	return t.setStopOrder(ctx, instrument, price, position, stopLossStopOrderType)
}

func (t *Tinkoff) setTakeProfit(
	ctx context.Context,
	instrument *investapi.Instrument,
	price *investapi.Quotation,
	position trengin.Position,
) (string, error) {
	return t.setStopOrder(ctx, instrument, price, position, takeProfitStopOrderType)
}

func (t *Tinkoff) setStopOrder(
	ctx context.Context,
	instrument *investapi.Instrument,
	stopPrice *investapi.Quotation,
	position trengin.Position,
	orderType stopOrderType,
) (string, error) {
	stopOrderDirection := investapi.StopOrderDirection_STOP_ORDER_DIRECTION_BUY
	if position.Type.IsLong() {
		stopOrderDirection = investapi.StopOrderDirection_STOP_ORDER_DIRECTION_SELL
	}
	reqStopOrderType := investapi.StopOrderType_STOP_ORDER_TYPE_STOP_LIMIT
	if orderType == takeProfitStopOrderType {
		reqStopOrderType = investapi.StopOrderType_STOP_ORDER_TYPE_TAKE_PROFIT
	}

	price := t.addProtectedSpread(position.Type, stopPrice, instrument.MinPriceIncrement)
	stopOrderRequest := &investapi.PostStopOrderRequest{
		Figi:           position.FIGI,
		Quantity:       position.Quantity,
		Price:          price,
		StopPrice:      stopPrice,
		Direction:      stopOrderDirection,
		AccountId:      t.accountID,
		ExpirationType: investapi.StopOrderExpirationType_STOP_ORDER_EXPIRATION_TYPE_GOOD_TILL_CANCEL,
		StopOrderType:  reqStopOrderType,
	}

	stopOrder, err := t.stopOrderClient.PostStopOrder(ctx, stopOrderRequest)
	if err != nil {
		t.logger.Info(
			"Failed to set stop order",
			zap.Any("stopOrderRequest", stopOrderRequest),
			zap.Error(err),
		)
		return "", fmt.Errorf("post stop order: %w", err)
	}

	t.logger.Info(
		"Stop order was set",
		zap.Any("stopOrderRequest", stopOrderRequest),
		zap.Any("stopOrder", stopOrder),
	)
	return stopOrder.StopOrderId, nil
}

func (t *Tinkoff) cancelStopOrder(ctx context.Context, id string) error {
	if id == "" {
		return nil
	}
	cancelStopOrderRequest := &investapi.CancelStopOrderRequest{
		AccountId:   t.accountID,
		StopOrderId: id,
	}
	loggerWithRequest := t.logger.With(zap.Any("cancelStopOrderRequest", cancelStopOrderRequest))

	_, err := t.stopOrderClient.CancelStopOrder(ctx, cancelStopOrderRequest)
	if status.Code(err) == codes.NotFound {
		loggerWithRequest.Warn("Stop order is not found", zap.Error(err))
		return nil
	}
	if err != nil {
		t.logger.Error("Failed to cancel stop order", zap.Error(err))
		return fmt.Errorf("cancel stop order: %w", err)
	}

	t.logger.Info("Stop order was canceled", zap.String("id", id))
	return nil
}

func (t *Tinkoff) stopLossPriceByOpen(
	openPrice *MoneyValue,
	action trengin.OpenPositionAction,
	minPriceIncrement *investapi.Quotation,
) *investapi.Quotation {
	stopLoss := openPrice.ToFloat() - action.StopLossIndent*action.Type.Multiplier()
	return t.convertFloatToQuotation(minPriceIncrement, stopLoss)
}

func (t *Tinkoff) takeProfitPriceByOpen(
	openPrice *MoneyValue,
	action trengin.OpenPositionAction,
	minPriceIncrement *investapi.Quotation,
) *investapi.Quotation {
	takeProfit := openPrice.ToFloat() + action.TakeProfitIndent*action.Type.Multiplier()
	return t.convertFloatToQuotation(minPriceIncrement, takeProfit)
}

func (t *Tinkoff) convertFloatToQuotation(
	minPriceIncrement *investapi.Quotation,
	stopLoss float64,
) *investapi.Quotation {
	stopOrderUnits, stopOrderNano := math.Modf(stopLoss)

	var roundStopOrderNano int32
	if minPriceIncrement != nil {
		roundStopOrderNano = int32(math.Round(stopOrderNano*10e8/float64(minPriceIncrement.GetNano()))) *
			minPriceIncrement.GetNano()
	}
	return &investapi.Quotation{
		Units: int64(stopOrderUnits),
		Nano:  roundStopOrderNano,
	}
}

func (t *Tinkoff) addProtectedSpread(
	positionType trengin.PositionType,
	price *investapi.Quotation,
	minPriceIncrement *investapi.Quotation,
) *investapi.Quotation {
	priceFloat := NewMoneyValue(price).ToFloat()
	protectiveSpread := priceFloat * t.protectiveSpread / 100
	return t.convertFloatToQuotation(minPriceIncrement, priceFloat-positionType.Multiplier()*protectiveSpread)
}

func (t *Tinkoff) cancelStopOrders(ctx context.Context, tinkoffPosition *tnkposition.Position) error {
	ctx = t.ctxWithMetadata(ctx)

	resp, err := t.stopOrderClient.GetStopOrders(ctx, &investapi.GetStopOrdersRequest{
		AccountId: t.accountID,
	})
	if err != nil {
		return err
	}

	orders := make(map[string]struct{})
	for _, order := range resp.StopOrders {
		orders[order.StopOrderId] = struct{}{}
	}

	stopLossID := tinkoffPosition.StopLossID()
	takeProfitID := tinkoffPosition.TakeProfitID()
	if _, ok := orders[stopLossID]; ok {
		if err := t.cancelStopOrder(ctx, stopLossID); err != nil {
			return fmt.Errorf("cancel stop loss: %w", err)
		}
	}
	if _, ok := orders[takeProfitID]; ok {
		if err := t.cancelStopOrder(ctx, takeProfitID); err != nil {
			return fmt.Errorf("cancel take profit: %w", err)
		}
	}
	return nil
}

func (t *Tinkoff) getLastPrice(ctx context.Context, figi string) (*investapi.Quotation, error) {
	prices, err := t.marketDataClient.GetLastPrices(ctx, &investapi.GetLastPricesRequest{
		Figi: []string{figi},
	})
	if err != nil {
		t.logger.Error("Failed to get last prices", zap.Error(err), zap.Any("figi", figi))
		return nil, fmt.Errorf("last prices: %w", err)
	}
	t.logger.Debug("Last prices were received", zap.Any("prices", prices))

	for _, p := range prices.GetLastPrices() {
		if p.Figi == figi {
			return p.Price, nil
		}
	}
	return nil, errors.New("figi not found")
}

func (t *Tinkoff) getExecutedOrderState(
	ctx context.Context,
	orderID string,
) (orderState *investapi.OrderState, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	for {
		orderState, err = t.getOrderState(ctx, orderID)
		//nolint: lll
		isFullFilled := orderState.GetExecutionReportStatus() == investapi.OrderExecutionReportStatus_EXECUTION_REPORT_STATUS_FILL &&
			orderState.LotsExecuted == orderState.LotsRequested
		if isFullFilled {
			break
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(250 * time.Millisecond):
		}
	}
	return orderState, nil
}

func (t *Tinkoff) getOrderState(ctx context.Context, orderID string) (orderState *investapi.OrderState, err error) {
	orderStateRequest := &investapi.GetOrderStateRequest{
		AccountId: t.accountID,
		OrderId:   orderID,
	}
	orderState, err = t.orderClient.GetOrderState(ctx, orderStateRequest)
	if err != nil {
		t.logger.Error(
			"Failed to get order state",
			zap.Error(err),
			zap.Any("orderStateRequest", orderStateRequest),
		)
		return nil, fmt.Errorf("get order state: %w", err)
	}
	return orderState, nil
}

func (t *Tinkoff) getInstrument(ctx context.Context, figi string) (*investapi.Instrument, error) {
	instrumentResponse, err := t.instrumentClient.GetInstrumentBy(t.ctxWithMetadata(ctx), &investapi.InstrumentRequest{
		IdType: investapi.InstrumentIdType_INSTRUMENT_ID_TYPE_FIGI,
		Id:     figi,
	})
	if err != nil {
		return nil, fmt.Errorf("get instrument by %s: %w", figi, err)
	}
	return instrumentResponse.GetInstrument(), nil
}

func (t *Tinkoff) conditionalOrdersFound(ctx context.Context, position *tnkposition.Position) (bool, error) {
	resp, err := t.stopOrderClient.GetStopOrders(ctx, &investapi.GetStopOrdersRequest{AccountId: t.accountID})
	if err != nil {
		return false, fmt.Errorf("get stop orders: %w", err)
	}
	for _, order := range resp.StopOrders {
		if order.StopOrderId == position.TakeProfitID() || order.StopOrderId == position.StopLossID() {
			return true, nil
		}
	}
	return false, nil
}
