// Package tnkbroker implements [trengin.Broker] using [Tinkoff Invest API].
// Supports multiple open positions at the same time.
// Commission in position is approximate.
//
// [Tinkoff Invest API]: https://tinkoff.github.io/investAPI/
package tnkbroker

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/evsamsonov/trengin/v2"
	"github.com/google/uuid"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/evsamsonov/tinkoff-broker/internal/tnkposition"
)

var _ trengin.Broker = &Tinkoff{}

const (
	defaultProtectiveSpread        = 1 // In percent
	defaultTradeStreamRetryTimeout = 1 * time.Minute
	defaultTradeStreamPingTimeout  = 6 * time.Minute
)

type Tinkoff struct {
	accountID                   string
	orderClient                 ordersServiceClient
	stopOrderClient             stopOrdersServiceClient
	tradeStreamClient           ordersStreamClient
	marketDataClient            marketDataServiceClient
	instrumentClient            instrumentsServiceClient
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

// New creates a new Tinkoff object. It takes [full-access token], todo?
// user account identifier.
//
// [full-access token]: https://tinkoff.github.io/investAPI/token/
func New(client *investgo.Client, accountID string, opts ...Option) (*Tinkoff, error) {
	tinkoff := &Tinkoff{
		accountID:                   accountID,
		protectiveSpread:            defaultProtectiveSpread,
		orderClient:                 client.NewOrdersServiceClient(),
		stopOrderClient:             client.NewStopOrdersServiceClient(),
		tradeStreamClient:           client.NewOrdersStreamClient(),
		marketDataClient:            client.NewMarketDataServiceClient(),
		instrumentClient:            client.NewInstrumentsServiceClient(),
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

// OpenPosition opens position, returns new position and channel for tracking position closing.
func (t *Tinkoff) OpenPosition(
	ctx context.Context,
	action trengin.OpenPositionAction,
) (trengin.Position, trengin.PositionClosed, error) {
	instrument, err := t.getInstrument(action.FIGI)
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
	if action.StopLossOffset != 0 {
		stopLoss := t.stopLossPriceByOpen(openPrice, action, instrument.MinPriceIncrement)
		stopLossID, err = t.setStopLoss(instrument, stopLoss, *position)
		if err != nil {
			return trengin.Position{}, nil, fmt.Errorf("set stop order: %w", err)
		}
	}
	if action.TakeProfitOffset != 0 {
		takeProfit := t.takeProfitPriceByOpen(openPrice, action, instrument.MinPriceIncrement)
		takeProfitID, err = t.setTakeProfit(instrument, takeProfit, *position)
		if err != nil {
			return trengin.Position{}, nil, fmt.Errorf("set stop order: %w", err)
		}
	}

	positionClosed := make(chan trengin.Position, 1)
	t.positionStorage.Store(tnkposition.NewPosition(position, instrument, stopLossID, takeProfitID, positionClosed))

	return *position, positionClosed, nil
}

// ChangeConditionalOrder changes stop loss and take profit of current position.
// It returns updated position.
func (t *Tinkoff) ChangeConditionalOrder(
	_ context.Context,
	action trengin.ChangeConditionalOrderAction,
) (trengin.Position, error) {
	tinkoffPosition, unlockPosition, err := t.positionStorage.Load(action.PositionID)
	if err != nil {
		return trengin.Position{}, fmt.Errorf("load position: %w", err)
	}
	defer unlockPosition()

	instrument := tinkoffPosition.Instrument()
	if action.StopLoss != 0 {
		if err := t.cancelStopOrder(tinkoffPosition.StopLossID()); err != nil {
			return trengin.Position{}, fmt.Errorf("cancel stop loss: %w", err)
		}

		stopLoss := t.convertFloatToQuotation(instrument.MinPriceIncrement, action.StopLoss)
		stopLossID, err := t.setStopLoss(instrument, stopLoss, tinkoffPosition.Position())
		if err != nil {
			return trengin.Position{}, fmt.Errorf("set stop loss: %w", err)
		}
		tinkoffPosition.SetStopLoss(stopLossID, action.StopLoss)
	}

	if action.TakeProfit != 0 {
		if err := t.cancelStopOrder(tinkoffPosition.TakeProfitID()); err != nil {
			return trengin.Position{}, fmt.Errorf("cancel take profit: %w", err)
		}

		takeProfit := t.convertFloatToQuotation(instrument.MinPriceIncrement, action.TakeProfit)
		takeProfitID, err := t.setTakeProfit(instrument, takeProfit, tinkoffPosition.Position())
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

	if err := t.cancelStopOrder(tinkoffPosition.StopLossID()); err != nil {
		return trengin.Position{}, fmt.Errorf("cancel stop loss: %w", err)
	}
	if err := t.cancelStopOrder(tinkoffPosition.TakeProfitID()); err != nil {
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

	tradesStream, err := t.tradeStreamClient.TradesStream([]string{t.accountID})
	if err != nil {
		return fmt.Errorf("trades stream: %w", err)
	}

	g := errgroup.Group{}
	g.Go(func() error {
		t.logger.Debug("Trades stream is connected")

		trades := tradesStream.Trades()
		for {
			select {
			case <-ctx.Done():
				return nil
			case orderTrades, ok := <-trades:
				if !ok {
					return nil
				}
				t.logger.Info("Order trades were received", zap.Any("orderTrades", orderTrades))

				if err := t.processOrderTrades(ctx, orderTrades); err != nil {
					return fmt.Errorf("process order trades: %w", err)
				}
			}
		}
	})

	g.Go(func() error {
		if err := tradesStream.Listen(); err != nil {
			return fmt.Errorf("listen: %w", err)
		}
		return nil
	})

	return g.Wait()
}

func (t *Tinkoff) processOrderTrades(ctx context.Context, orderTrades *pb.OrderTrades) error {
	if orderTrades.AccountId != t.accountID {
		return nil
	}
	err := t.positionStorage.ForEach(func(tinkoffPosition *tnkposition.Position) error {
		position := tinkoffPosition.Position()
		if orderTrades.Figi != position.FIGI {
			return nil
		}
		longClosed := position.IsLong() && orderTrades.Direction == pb.OrderDirection_ORDER_DIRECTION_SELL
		shortClosed := position.IsShort() && orderTrades.Direction == pb.OrderDirection_ORDER_DIRECTION_BUY
		if !longClosed && !shortClosed {
			return nil
		}

		conditionalOrdersFound, err := t.conditionalOrdersFound(tinkoffPosition)
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

		if err := t.cancelStopOrders(tinkoffPosition); err != nil {
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

func (t *Tinkoff) openMarketOrder(
	ctx context.Context,
	instrument *pb.Instrument,
	positionType trengin.PositionType,
	quantity int64,
) (*MoneyValue, *MoneyValue, error) {
	direction := pb.OrderDirection_ORDER_DIRECTION_BUY
	if positionType.IsShort() {
		direction = pb.OrderDirection_ORDER_DIRECTION_SELL
	}
	price, err := t.getLastPrice(instrument.Figi)
	if err != nil {
		return nil, nil, fmt.Errorf("get last price: %w", err)
	}
	orderRequest := &investgo.PostOrderRequest{
		InstrumentId: instrument.Figi,
		Quantity:     quantity,
		Direction:    direction,
		AccountId:    t.accountID,
		Price:        t.addProtectedSpread(positionType.Inverse(), price, instrument.MinPriceIncrement),
		OrderType:    pb.OrderType_ORDER_TYPE_LIMIT,
		OrderId:      uuid.New().String(),
	}

	order, err := t.orderClient.PostOrder(orderRequest)
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
	instrument *pb.Instrument,
	price *pb.Quotation,
	position trengin.Position,
) (string, error) {
	return t.setStopOrder(instrument, price, position, stopLossStopOrderType)
}

func (t *Tinkoff) setTakeProfit(
	instrument *pb.Instrument,
	price *pb.Quotation,
	position trengin.Position,
) (string, error) {
	return t.setStopOrder(instrument, price, position, takeProfitStopOrderType)
}

func (t *Tinkoff) setStopOrder(
	instrument *pb.Instrument,
	stopPrice *pb.Quotation,
	position trengin.Position,
	orderType stopOrderType,
) (string, error) {
	stopOrderDirection := pb.StopOrderDirection_STOP_ORDER_DIRECTION_BUY
	if position.Type.IsLong() {
		stopOrderDirection = pb.StopOrderDirection_STOP_ORDER_DIRECTION_SELL
	}
	reqStopOrderType := pb.StopOrderType_STOP_ORDER_TYPE_STOP_LIMIT
	if orderType == takeProfitStopOrderType {
		reqStopOrderType = pb.StopOrderType_STOP_ORDER_TYPE_TAKE_PROFIT
	}

	price := t.addProtectedSpread(position.Type, stopPrice, instrument.MinPriceIncrement)
	stopOrderRequest := &investgo.PostStopOrderRequest{
		InstrumentId:   position.FIGI,
		Quantity:       position.Quantity,
		Price:          price,
		StopPrice:      stopPrice,
		Direction:      stopOrderDirection,
		AccountId:      t.accountID,
		ExpirationType: pb.StopOrderExpirationType_STOP_ORDER_EXPIRATION_TYPE_GOOD_TILL_CANCEL,
		StopOrderType:  reqStopOrderType,
	}

	stopOrder, err := t.stopOrderClient.PostStopOrder(stopOrderRequest)
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

func (t *Tinkoff) cancelStopOrder(id string) error {
	if id == "" {
		return nil
	}
	logger := t.logger.With(zap.String("id", id))

	_, err := t.stopOrderClient.CancelStopOrder(t.accountID, id)
	if status.Code(err) == codes.NotFound {
		logger.Warn("Stop order is not found", zap.Error(err))
		return nil
	}
	if err != nil {
		logger.Error("Failed to cancel stop order", zap.Error(err))
		return fmt.Errorf("cancel stop order: %w", err)
	}

	t.logger.Info("Stop order was canceled", zap.String("id", id))
	return nil
}

func (t *Tinkoff) stopLossPriceByOpen(
	openPrice *MoneyValue,
	action trengin.OpenPositionAction,
	minPriceIncrement *pb.Quotation,
) *pb.Quotation {
	stopLoss := openPrice.ToFloat() - action.StopLossOffset*action.Type.Multiplier()
	return t.convertFloatToQuotation(minPriceIncrement, stopLoss)
}

func (t *Tinkoff) takeProfitPriceByOpen(
	openPrice *MoneyValue,
	action trengin.OpenPositionAction,
	minPriceIncrement *pb.Quotation,
) *pb.Quotation {
	takeProfit := openPrice.ToFloat() + action.TakeProfitOffset*action.Type.Multiplier()
	return t.convertFloatToQuotation(minPriceIncrement, takeProfit)
}

func (t *Tinkoff) convertFloatToQuotation(
	minPriceIncrement *pb.Quotation,
	stopLoss float64,
) *pb.Quotation {
	stopOrderUnits, stopOrderNano := math.Modf(stopLoss)

	var roundStopOrderNano int32
	if minPriceIncrement != nil {
		roundStopOrderNano = int32(math.Round(stopOrderNano*10e8/float64(minPriceIncrement.GetNano()))) *
			minPriceIncrement.GetNano()
	}
	return &pb.Quotation{
		Units: int64(stopOrderUnits),
		Nano:  roundStopOrderNano,
	}
}

func (t *Tinkoff) addProtectedSpread(
	positionType trengin.PositionType,
	price *pb.Quotation,
	minPriceIncrement *pb.Quotation,
) *pb.Quotation {
	priceFloat := NewMoneyValue(price).ToFloat()
	protectiveSpread := priceFloat * t.protectiveSpread / 100
	return t.convertFloatToQuotation(minPriceIncrement, priceFloat-positionType.Multiplier()*protectiveSpread)
}

func (t *Tinkoff) cancelStopOrders(tinkoffPosition *tnkposition.Position) error {
	resp, err := t.stopOrderClient.GetStopOrders(t.accountID)
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
		if err := t.cancelStopOrder(stopLossID); err != nil {
			return fmt.Errorf("cancel stop loss: %w", err)
		}
	}
	if _, ok := orders[takeProfitID]; ok {
		if err := t.cancelStopOrder(takeProfitID); err != nil {
			return fmt.Errorf("cancel take profit: %w", err)
		}
	}
	return nil
}

func (t *Tinkoff) getLastPrice(figi string) (*pb.Quotation, error) {
	prices, err := t.marketDataClient.GetLastPrices([]string{figi})
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
) (orderState *investgo.GetOrderStateResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	for {
		orderState, err = t.getOrderState(orderID)
		//nolint: lll
		isFullFilled := orderState.GetExecutionReportStatus() == pb.OrderExecutionReportStatus_EXECUTION_REPORT_STATUS_FILL &&
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

func (t *Tinkoff) getOrderState(orderID string) (orderState *investgo.GetOrderStateResponse, err error) {
	orderState, err = t.orderClient.GetOrderState(t.accountID, orderID, pb.PriceType_PRICE_TYPE_CURRENCY)
	if err != nil {
		t.logger.Error("Failed to get order state", zap.Error(err), zap.Any("orderID", orderID))
		return nil, fmt.Errorf("get order state: %w", err)
	}
	return orderState, nil
}

func (t *Tinkoff) getInstrument(figi string) (*pb.Instrument, error) {
	instrumentResponse, err := t.instrumentClient.InstrumentByFigi(figi)
	if err != nil {
		return nil, fmt.Errorf("get instrument by %s: %w", figi, err)
	}
	return instrumentResponse.GetInstrument(), nil
}

func (t *Tinkoff) conditionalOrdersFound(position *tnkposition.Position) (bool, error) {
	resp, err := t.stopOrderClient.GetStopOrders(t.accountID)
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
