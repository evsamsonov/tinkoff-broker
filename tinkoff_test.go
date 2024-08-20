package tnkbroker

import (
	"context"
	"testing"
	"time"

	"github.com/evsamsonov/trengin/v2"
	"github.com/google/uuid"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/undefinedlabs/go-mpatch"
	"go.uber.org/zap"

	"github.com/evsamsonov/tinkoff-broker/internal/tnkposition"
)

const (
	float64EqualityThreshold = 1e-6
)

func TestTinkoff_OpenPosition(t *testing.T) {
	type testWant struct {
		positionType       trengin.PositionType
		orderDirection     pb.OrderDirection
		stopOrderDirection pb.StopOrderDirection
		openPrice          *pb.MoneyValue
		stopLoss           *pb.Quotation
		takeProfit         *pb.Quotation
		stopLossID         string
		takeProfitID       string
	}

	tests := []struct {
		name               string
		openPositionAction trengin.OpenPositionAction
		want               testWant
	}{
		{
			name: "long with stop loss and take profit",
			openPositionAction: trengin.OpenPositionAction{
				FIGI:             "FUTSBRF06220",
				Type:             trengin.Long,
				Quantity:         2,
				StopLossOffset:   11.5,
				TakeProfitOffset: 20.1,
			},
			want: testWant{
				orderDirection:     pb.OrderDirection_ORDER_DIRECTION_BUY,
				stopOrderDirection: pb.StopOrderDirection_STOP_ORDER_DIRECTION_SELL,
				positionType:       trengin.Long,
				openPrice:          &pb.MoneyValue{Units: 148, Nano: 0.2 * 10e8},
				stopLoss:           &pb.Quotation{Units: 136, Nano: 0.7 * 10e8},
				takeProfit:         &pb.Quotation{Units: 168, Nano: 0.3 * 10e8},
				stopLossID:         "123",
				takeProfitID:       "321",
			},
		},
		{
			name: "short with stop loss and take profit",
			openPositionAction: trengin.OpenPositionAction{
				FIGI:             "FUTSBRF06220",
				Type:             trengin.Short,
				Quantity:         2,
				StopLossOffset:   11.5,
				TakeProfitOffset: 20.1,
			},
			want: testWant{
				orderDirection:     pb.OrderDirection_ORDER_DIRECTION_SELL,
				stopOrderDirection: pb.StopOrderDirection_STOP_ORDER_DIRECTION_BUY,
				positionType:       trengin.Short,
				openPrice:          &pb.MoneyValue{Units: 148, Nano: 0.2 * 10e8},
				stopLoss:           &pb.Quotation{Units: 159, Nano: 0.7 * 10e8},
				takeProfit:         &pb.Quotation{Units: 128, Nano: 0.1 * 10e8},
				stopLossID:         "123",
				takeProfitID:       "321",
			},
		},
		{
			name: "without stop loss and take profit",
			openPositionAction: trengin.OpenPositionAction{
				FIGI:             "FUTSBRF06220",
				Type:             trengin.Long,
				Quantity:         2,
				StopLossOffset:   0.0,
				TakeProfitOffset: 0.0,
			},
			want: testWant{
				orderDirection:     pb.OrderDirection_ORDER_DIRECTION_BUY,
				stopOrderDirection: pb.StopOrderDirection_STOP_ORDER_DIRECTION_SELL,
				positionType:       trengin.Long,
				openPrice:          &pb.MoneyValue{Units: 148, Nano: 0.2 * 10e8},
				stopLoss:           &pb.Quotation{},
				takeProfit:         &pb.Quotation{},
				stopLossID:         "",
				takeProfitID:       "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ordersServiceClient := &mockOrdersServiceClient{}
			stopOrdersServiceClient := &mockStopOrdersServiceClient{}
			marketDataServiceClient := &mockMarketDataServiceClient{}
			instrumentServiceClient := &mockInstrumentsServiceClient{}

			patch, err := mpatch.PatchMethod(uuid.New, func() uuid.UUID {
				return uuid.MustParse("8942e9ae-e4e1-11ec-8fea-0242ac120002")
			})
			defer func() { assert.NoError(t, patch.Unpatch()) }()
			assert.NoError(t, err)

			tinkoff := &Tinkoff{
				accountID:        "123",
				orderClient:      ordersServiceClient,
				stopOrderClient:  stopOrdersServiceClient,
				marketDataClient: marketDataServiceClient,
				instrumentClient: instrumentServiceClient,
				positionStorage:  tnkposition.NewStorage(),
				logger:           zap.NewNop(),
			}

			instrumentServiceClient.
				On("InstrumentByFigi", "FUTSBRF06220").
				Return(
					&investgo.InstrumentResponse{
						InstrumentResponse: &pb.InstrumentResponse{
							Instrument: &pb.Instrument{
								Figi:              "FUTSBRF06220",
								MinPriceIncrement: &pb.Quotation{Units: 0, Nano: 0.1 * 10e8},
							},
						},
					}, nil)

			marketDataServiceClient.
				On("GetLastPrices", []string{"FUTSBRF06220"}).
				Return(&investgo.GetLastPricesResponse{
					GetLastPricesResponse: &pb.GetLastPricesResponse{
						LastPrices: []*pb.LastPrice{{
							Figi:  "FUTSBRF06220",
							Price: &pb.Quotation{Units: 100},
						}},
					}}, nil)

			ordersServiceClient.
				On("PostOrder", &investgo.PostOrderRequest{
					InstrumentId: "FUTSBRF06220",
					Quantity:     2,
					Price:        &pb.Quotation{Units: 100},
					Direction:    tt.want.orderDirection,
					AccountId:    "123",
					OrderType:    pb.OrderType_ORDER_TYPE_LIMIT,
					OrderId:      "8942e9ae-e4e1-11ec-8fea-0242ac120002",
				}).Return(
				&investgo.PostOrderResponse{
					PostOrderResponse: &pb.PostOrderResponse{
						ExecutionReportStatus: pb.OrderExecutionReportStatus_EXECUTION_REPORT_STATUS_FILL,
						ExecutedOrderPrice:    tt.want.openPrice,
						OrderId:               "1953",
					}}, nil)

			ordersServiceClient.
				On("GetOrderState", "123", "1953", pb.PriceType_PRICE_TYPE_CURRENCY).
				Return(&investgo.GetOrderStateResponse{
					OrderState: &pb.OrderState{
						InitialCommission:     &pb.MoneyValue{Units: 12, Nano: 0.1 * 10e8},
						ExecutionReportStatus: pb.OrderExecutionReportStatus_EXECUTION_REPORT_STATUS_FILL,
						AveragePositionPrice:  tt.want.openPrice,
					}}, nil)

			if tt.openPositionAction.StopLossOffset != 0 {
				stopOrdersServiceClient.
					On("PostStopOrder", &investgo.PostStopOrderRequest{
						InstrumentId:   "FUTSBRF06220",
						Quantity:       2,
						Price:          tt.want.stopLoss,
						StopPrice:      tt.want.stopLoss,
						Direction:      tt.want.stopOrderDirection,
						AccountId:      "123",
						ExpirationType: pb.StopOrderExpirationType_STOP_ORDER_EXPIRATION_TYPE_GOOD_TILL_CANCEL,
						StopOrderType:  pb.StopOrderType_STOP_ORDER_TYPE_STOP_LIMIT,
					}).
					Return(&investgo.PostStopOrderResponse{
						PostStopOrderResponse: &pb.PostStopOrderResponse{
							StopOrderId: "123",
						},
					}, nil).Once()
			}

			if tt.openPositionAction.TakeProfitOffset != 0 {
				stopOrdersServiceClient.
					On("PostStopOrder", &investgo.PostStopOrderRequest{
						InstrumentId:   "FUTSBRF06220",
						Quantity:       2,
						Price:          tt.want.takeProfit,
						StopPrice:      tt.want.takeProfit,
						Direction:      tt.want.stopOrderDirection,
						AccountId:      "123",
						ExpirationType: pb.StopOrderExpirationType_STOP_ORDER_EXPIRATION_TYPE_GOOD_TILL_CANCEL,
						StopOrderType:  pb.StopOrderType_STOP_ORDER_TYPE_TAKE_PROFIT,
					}).
					Return(&investgo.PostStopOrderResponse{
						PostStopOrderResponse: &pb.PostStopOrderResponse{
							StopOrderId: "321",
						}}, nil).Once()
			}

			position, _, err := tinkoff.OpenPosition(context.Background(), tt.openPositionAction)
			assert.NoError(t, err)

			assert.Equal(t, tt.want.positionType, position.Type)
			assert.InEpsilon(t, NewMoneyValue(tt.want.openPrice).ToFloat(), position.OpenPrice, float64EqualityThreshold)

			wantStopLoss := NewMoneyValue(tt.want.stopLoss).ToFloat()
			if wantStopLoss != 0 {
				assert.InEpsilon(t, wantStopLoss, position.StopLoss, float64EqualityThreshold)
			} else {
				assert.Equal(t, 0., position.StopLoss)
			}

			wantTakeProfit := NewMoneyValue(tt.want.takeProfit).ToFloat()
			if wantTakeProfit != 0 {
				assert.InEpsilon(t, wantTakeProfit, position.TakeProfit, float64EqualityThreshold)
			} else {
				assert.Equal(t, 0., position.TakeProfit)
			}

			tinkoffPosition, _, err := tinkoff.positionStorage.Load(position.ID)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.stopLossID, tinkoffPosition.StopLossID())
			assert.Equal(t, tt.want.takeProfitID, tinkoffPosition.TakeProfitID())
		})
	}
}

func TestTinkoff_ChangeConditionalOrder_noOpenPosition(t *testing.T) {
	tinkoff := &Tinkoff{
		positionStorage: tnkposition.NewStorage(),
	}
	_, err := tinkoff.ChangeConditionalOrder(context.Background(), trengin.ChangeConditionalOrderAction{})
	assert.Errorf(t, err, "no open position")
}

func TestTinkoff_ChangeConditionalOrder(t *testing.T) {
	type testWant struct {
		stopLoss           *pb.Quotation
		takeProfit         *pb.Quotation
		stopOrderDirection pb.StopOrderDirection
		stopLossID         string
		takeProfitID       string
	}

	tests := []struct {
		name                       string
		changeConditionOrderAction trengin.ChangeConditionalOrderAction
		positionType               trengin.PositionType
		want                       testWant
	}{
		{
			name: "stop loss and take profit equal zero",
			changeConditionOrderAction: trengin.ChangeConditionalOrderAction{
				PositionID: trengin.PositionID{},
				StopLoss:   0,
				TakeProfit: 0,
			},
			want: testWant{
				stopLossID:   "1",
				takeProfitID: "3",
				stopLoss:     &pb.Quotation{},
				takeProfit:   &pb.Quotation{},
			},
		},
		{
			name: "long position, stop loss and take profit set are given",
			changeConditionOrderAction: trengin.ChangeConditionalOrderAction{
				PositionID: trengin.PositionID{},
				StopLoss:   123.43,
				TakeProfit: 156.31,
			},
			positionType: trengin.Long,
			want: testWant{
				stopLoss:           &pb.Quotation{Units: 123, Nano: 0.43 * 10e8},
				takeProfit:         &pb.Quotation{Units: 156, Nano: 0.31 * 10e8},
				stopOrderDirection: pb.StopOrderDirection_STOP_ORDER_DIRECTION_SELL,
				stopLossID:         "2",
				takeProfitID:       "4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ordersServiceClient := &mockOrdersServiceClient{}
			stopOrdersServiceClient := &mockStopOrdersServiceClient{}
			instrumentServiceClient := &mockInstrumentsServiceClient{}
			positionStorage := tnkposition.NewStorage()

			tinkoff := &Tinkoff{
				accountID:        "123",
				orderClient:      ordersServiceClient,
				stopOrderClient:  stopOrdersServiceClient,
				instrumentClient: instrumentServiceClient,
				positionStorage:  positionStorage,
				logger:           zap.NewNop(),
			}
			positionStorage.Store(tnkposition.NewPosition(
				&trengin.Position{Type: tt.positionType, Quantity: 2, FIGI: "FUTSBRF06220"},
				&pb.Instrument{
					Figi:              "FUTSBRF06220",
					MinPriceIncrement: &pb.Quotation{Units: 0, Nano: 0.01 * 10e8},
				},
				"1",
				"3",
				make(chan trengin.Position, 1),
			))

			instrumentServiceClient.On("GetInstrumentBy", mock.Anything, &pb.InstrumentRequest{
				IdType: pb.InstrumentIdType_INSTRUMENT_ID_TYPE_FIGI,
				Id:     "FUTSBRF06220",
			}).Return(&pb.InstrumentResponse{
				Instrument: &pb.Instrument{
					Figi:              "FUTSBRF06220",
					MinPriceIncrement: &pb.Quotation{Units: 0, Nano: 0.1 * 10e8},
				},
			}, nil)

			if tt.changeConditionOrderAction.StopLoss != 0 {
				stopOrdersServiceClient.
					On("CancelStopOrder", "123", "1").
					Return(
						&investgo.CancelStopOrderResponse{
							CancelStopOrderResponse: &pb.CancelStopOrderResponse{},
						},
						nil,
					).
					Once()

				stopOrdersServiceClient.On("PostStopOrder", &investgo.PostStopOrderRequest{
					InstrumentId:   "FUTSBRF06220",
					Quantity:       2,
					Price:          tt.want.stopLoss,
					StopPrice:      tt.want.stopLoss,
					Direction:      tt.want.stopOrderDirection,
					AccountId:      "123",
					ExpirationType: pb.StopOrderExpirationType_STOP_ORDER_EXPIRATION_TYPE_GOOD_TILL_CANCEL,
					StopOrderType:  pb.StopOrderType_STOP_ORDER_TYPE_STOP_LIMIT,
				}).Return(&investgo.PostStopOrderResponse{
					PostStopOrderResponse: &pb.PostStopOrderResponse{
						StopOrderId: "2",
					},
				}, nil).Once()
			}

			if tt.changeConditionOrderAction.TakeProfit != 0 {
				stopOrdersServiceClient.
					On("CancelStopOrder", "123", "3").
					Return(&investgo.CancelStopOrderResponse{}, nil).
					Once()

				stopOrdersServiceClient.
					On("PostStopOrder", &investgo.PostStopOrderRequest{
						InstrumentId:   "FUTSBRF06220",
						Quantity:       2,
						Price:          tt.want.takeProfit,
						StopPrice:      tt.want.takeProfit,
						Direction:      tt.want.stopOrderDirection,
						AccountId:      "123",
						ExpirationType: pb.StopOrderExpirationType_STOP_ORDER_EXPIRATION_TYPE_GOOD_TILL_CANCEL,
						StopOrderType:  pb.StopOrderType_STOP_ORDER_TYPE_TAKE_PROFIT,
					}).
					Return(&investgo.PostStopOrderResponse{
						PostStopOrderResponse: &pb.PostStopOrderResponse{
							StopOrderId: "4",
						}}, nil).
					Once()
			}

			position, err := tinkoff.ChangeConditionalOrder(context.Background(), trengin.ChangeConditionalOrderAction{
				PositionID: trengin.PositionID{},
				StopLoss:   tt.changeConditionOrderAction.StopLoss,
				TakeProfit: tt.changeConditionOrderAction.TakeProfit,
			})
			assert.NoError(t, err)

			wantStopLoss := NewMoneyValue(tt.want.stopLoss).ToFloat()
			if wantStopLoss == 0 {
				assert.Zero(t, position.StopLoss)
			} else {
				assert.InEpsilon(t, NewMoneyValue(tt.want.stopLoss).ToFloat(), position.StopLoss, float64EqualityThreshold)
			}

			wantTakeProfit := NewMoneyValue(tt.want.takeProfit).ToFloat()
			if wantTakeProfit == 0 {
				assert.Zero(t, position.StopLoss)
			} else {
				assert.InEpsilon(t, NewMoneyValue(tt.want.takeProfit).ToFloat(), position.TakeProfit, float64EqualityThreshold)
			}

			tinkoffPosition, _, err := tinkoff.positionStorage.Load(position.ID)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.stopLossID, tinkoffPosition.StopLossID())
			assert.Equal(t, tt.want.takeProfitID, tinkoffPosition.TakeProfitID())
		})
	}
}

func TestTinkoff_ClosePosition_noOpenPosition(t *testing.T) {
	tinkoff := &Tinkoff{
		positionStorage: tnkposition.NewStorage(),
	}
	_, err := tinkoff.ClosePosition(context.Background(), trengin.ClosePositionAction{})
	assert.Errorf(t, err, "no open position")
}

func TestTinkoff_ClosePosition(t *testing.T) {
	tests := []struct {
		name               string
		positionType       trengin.PositionType
		wantOrderDirection pb.OrderDirection
		wantClosePrice     *pb.MoneyValue
	}{
		{
			name:               "long",
			positionType:       trengin.Long,
			wantOrderDirection: pb.OrderDirection_ORDER_DIRECTION_SELL,
			wantClosePrice:     &pb.MoneyValue{Units: 148, Nano: 0.2 * 10e8},
		},
		{
			name:               "short",
			positionType:       trengin.Short,
			wantOrderDirection: pb.OrderDirection_ORDER_DIRECTION_BUY,
			wantClosePrice:     &pb.MoneyValue{Units: 256, Nano: 0.3 * 10e8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ordersServiceClient := &mockOrdersServiceClient{}
			stopOrdersServiceClient := &mockStopOrdersServiceClient{}
			marketDataServiceClient := &mockMarketDataServiceClient{}
			instrumentServiceClient := &mockInstrumentsServiceClient{}

			patch, err := mpatch.PatchMethod(uuid.New, func() uuid.UUID {
				return uuid.MustParse("8942e9ae-e4e1-11ec-8fea-0242ac120002")
			})
			defer func() { assert.NoError(t, patch.Unpatch()) }()
			assert.NoError(t, err)

			positionStorage := tnkposition.NewStorage()
			tinkoff := &Tinkoff{
				accountID:        "123",
				orderClient:      ordersServiceClient,
				stopOrderClient:  stopOrdersServiceClient,
				marketDataClient: marketDataServiceClient,
				instrumentClient: instrumentServiceClient,
				positionStorage:  positionStorage,
				logger:           zap.NewNop(),
			}
			pos, err := trengin.NewPosition(
				trengin.NewOpenPositionAction("FUTSBRF06220", tt.positionType, 2, 0, 0),
				time.Now(),
				150,
			)
			assert.NoError(t, err)

			positionStorage.Store(tnkposition.NewPosition(
				pos,
				&pb.Instrument{
					Figi:              "FUTSBRF06220",
					MinPriceIncrement: &pb.Quotation{Units: 0, Nano: 0.01 * 10e8},
				},
				"1",
				"3",
				make(chan trengin.Position, 1),
			))

			stopOrdersServiceClient.
				On("CancelStopOrder", "123", "1").
				Return(&investgo.CancelStopOrderResponse{}, nil).
				Once()

			stopOrdersServiceClient.
				On("CancelStopOrder", "123", "3").
				Return(&investgo.CancelStopOrderResponse{}, nil).
				Once()

			ordersServiceClient.On("PostOrder", &investgo.PostOrderRequest{
				InstrumentId: "FUTSBRF06220",
				Quantity:     2,
				Direction:    tt.wantOrderDirection,
				AccountId:    "123",
				OrderType:    pb.OrderType_ORDER_TYPE_LIMIT,
				OrderId:      "8942e9ae-e4e1-11ec-8fea-0242ac120002",
				Price:        &pb.Quotation{Units: 100},
			}).Return(&investgo.PostOrderResponse{
				PostOrderResponse: &pb.PostOrderResponse{
					ExecutionReportStatus: pb.OrderExecutionReportStatus_EXECUTION_REPORT_STATUS_FILL,
					ExecutedOrderPrice:    tt.wantClosePrice,
					OrderId:               "1953",
				},
			}, nil)

			marketDataServiceClient.
				On("GetLastPrices", []string{"FUTSBRF06220"}).
				Return(&investgo.GetLastPricesResponse{
					GetLastPricesResponse: &pb.GetLastPricesResponse{
						LastPrices: []*pb.LastPrice{{
							Figi:  "FUTSBRF06220",
							Price: &pb.Quotation{Units: 100},
						}},
					}}, nil)

			ordersServiceClient.
				On("GetOrderState", "123", "1953", pb.PriceType_PRICE_TYPE_CURRENCY).
				Return(&investgo.GetOrderStateResponse{
					OrderState: &pb.OrderState{
						InitialCommission:     &pb.MoneyValue{Units: 12, Nano: 0.1 * 10e8},
						ExecutionReportStatus: pb.OrderExecutionReportStatus_EXECUTION_REPORT_STATUS_FILL,
						AveragePositionPrice:  tt.wantClosePrice,
					},
				}, nil)

			position, err := tinkoff.ClosePosition(context.Background(), trengin.ClosePositionAction{
				PositionID: pos.ID,
			})
			assert.NoError(t, err)
			assert.InEpsilon(t, NewMoneyValue(tt.wantClosePrice).ToFloat(), position.ClosePrice, float64EqualityThreshold)
		})
	}
}

func TestTinkoff_stopLossPriceByOpen(t *testing.T) {
	tests := []struct {
		name      string
		openPrice *pb.MoneyValue
		action    trengin.OpenPositionAction
		want      *pb.Quotation
	}{
		{
			name: "long nano is zero",
			openPrice: &pb.MoneyValue{
				Units: 123,
				Nano:  0,
			},
			action: trengin.OpenPositionAction{
				Type:           trengin.Long,
				StopLossOffset: 5,
			},
			want: &pb.Quotation{
				Units: 118,
				Nano:  0,
			},
		},
		{
			name: "long nano without overflow",
			openPrice: &pb.MoneyValue{
				Units: 123,
				Nano:  430000000,
			},
			action: trengin.OpenPositionAction{
				Type:           trengin.Long,
				StopLossOffset: 50.5,
			},
			want: &pb.Quotation{
				Units: 72,
				Nano:  930000000,
			},
		},
		{
			name: "long nano with overflow",
			openPrice: &pb.MoneyValue{
				Units: 123,
				Nano:  530000000,
			},
			action: trengin.OpenPositionAction{
				Type:           trengin.Long,
				StopLossOffset: 50.556,
			},
			want: &pb.Quotation{
				Units: 72,
				Nano:  974000000,
			},
		},
		{
			name: "short nano is zero",
			openPrice: &pb.MoneyValue{
				Units: 123,
				Nano:  0,
			},
			action: trengin.OpenPositionAction{
				Type:           trengin.Short,
				StopLossOffset: 5,
			},
			want: &pb.Quotation{
				Units: 128,
				Nano:  0,
			},
		},
		{
			name: "short nano without overflow",
			openPrice: &pb.MoneyValue{
				Units: 123,
				Nano:  430000000,
			},
			action: trengin.OpenPositionAction{
				Type:           trengin.Short,
				StopLossOffset: 50.4,
			},
			want: &pb.Quotation{
				Units: 173,
				Nano:  830000000,
			},
		},
		{
			name: "short nano with overflow",
			openPrice: &pb.MoneyValue{
				Units: 123,
				Nano:  530000000,
			},
			action: trengin.OpenPositionAction{
				Type:           trengin.Short,
				StopLossOffset: 50.556,
			},
			want: &pb.Quotation{
				Units: 174,
				Nano:  86000000,
			},
		},
	}

	tinkoff := Tinkoff{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			openPrice := NewMoneyValue(tt.openPrice)
			quotation := tinkoff.stopLossPriceByOpen(openPrice, tt.action, &pb.Quotation{
				Nano: 1000000,
			})
			assert.Equal(t, tt.want, quotation)
		})
	}
}

func TestTinkoff_processOrderTrades(t *testing.T) {
	ordersServiceClient := &mockOrdersServiceClient{}
	stopOrdersServiceClient := &mockStopOrdersServiceClient{}

	closed := make(chan trengin.Position, 1)

	stopOrdersServiceClient.
		On("GetStopOrders", "123").
		Return(&investgo.GetStopOrdersResponse{
			GetStopOrdersResponse: &pb.GetStopOrdersResponse{
				StopOrders: []*pb.StopOrder{
					{StopOrderId: "10"},
					{StopOrderId: "30"},
				},
			},
		}, nil)

	stopOrdersServiceClient.
		On("CancelStopOrder", "123", "1").
		Return(&investgo.CancelStopOrderResponse{}, nil)

	stopOrdersServiceClient.
		On("CancelStopOrder", "123", "3").
		Return(&investgo.CancelStopOrderResponse{}, nil)

	ordersServiceClient.
		On("GetOrderState", "123", "1953465028754600565", pb.PriceType_PRICE_TYPE_CURRENCY).
		Return(&investgo.GetOrderStateResponse{
			OrderState: &pb.OrderState{
				InitialCommission:     &pb.MoneyValue{Units: 125, Nano: 0.6 * 10e8},
				ExecutionReportStatus: pb.OrderExecutionReportStatus_EXECUTION_REPORT_STATUS_FILL,
				AveragePositionPrice:  &pb.MoneyValue{Units: 174, Nano: 0.7 * 10e8},
			},
		}, nil)

	positionStorage := tnkposition.NewStorage()
	tinkoff := &Tinkoff{
		accountID:       "123",
		orderClient:     ordersServiceClient,
		stopOrderClient: stopOrdersServiceClient,
		positionStorage: positionStorage,
		logger:          zap.NewNop(),
	}
	pos, err := trengin.NewPosition(
		trengin.NewOpenPositionAction("FUTSBRF06220", trengin.Long, 4, 0, 0),
		time.Now(),
		150,
	)
	assert.NoError(t, err)
	positionStorage.Store(tnkposition.NewPosition(
		pos,
		&pb.Instrument{
			Figi:              "FUTSBRF06220",
			MinPriceIncrement: &pb.Quotation{Units: 0, Nano: 0.01 * 10e8},
			Lot:               1,
		},
		"1",
		"3",
		closed,
	))

	ot := &pb.OrderTrades{
		OrderId:   "1953465028754600565",
		Direction: pb.OrderDirection_ORDER_DIRECTION_SELL,
		Figi:      "FUTSBRF06220",
		Trades: []*pb.OrderTrade{
			{Price: &pb.Quotation{Units: 112, Nano: 0.3 * 10e8}, Quantity: 2},
			{Price: &pb.Quotation{Units: 237, Nano: 0.1 * 10e8}, Quantity: 1},
		},
		AccountId: "123",
	}

	err = tinkoff.processOrderTrades(context.Background(), ot)
	assert.NoError(t, err)

	ot2 := &pb.OrderTrades{
		OrderId:   "1953465028754600565",
		Direction: pb.OrderDirection_ORDER_DIRECTION_SELL,
		Figi:      "FUTSBRF06220",
		Trades: []*pb.OrderTrade{
			{Price: &pb.Quotation{Units: 237, Nano: 0.1 * 10e8}, Quantity: 1},
		},
		AccountId: "123",
	}

	err = tinkoff.processOrderTrades(context.Background(), ot2)
	assert.NoError(t, err)

	select {
	case position := <-closed:
		assert.InEpsilon(t, 174.7, position.ClosePrice, float64EqualityThreshold)
		assert.InEpsilon(t, 125.6, position.Commission, float64EqualityThreshold)
	default:
		assert.Fail(t, "Failed to get closed position")
	}
}

func TestTinkoff_addProtectedSpread(t *testing.T) {
	var tests = []struct {
		name  string
		pType trengin.PositionType
		price *pb.Quotation
		want  *pb.Quotation
	}{
		{
			name:  "long",
			pType: trengin.Long,
			price: &pb.Quotation{
				Units: 237,
				Nano:  0.1 * 10e8,
			},
			want: &pb.Quotation{
				Units: 225,
				Nano:  0.2 * 10e8,
			},
		},
		{
			name:  "short",
			pType: trengin.Short,
			price: &pb.Quotation{
				Units: 237,
				Nano:  0.1 * 10e8,
			},
			want: &pb.Quotation{
				Units: 248,
				Nano:  1 * 10e8,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tinkoff := Tinkoff{
				protectiveSpread: 5,
			}
			result := tinkoff.addProtectedSpread(tt.pType, tt.price, &pb.Quotation{
				Units: 0,
				Nano:  0.1 * 10e8,
			})
			assert.Equal(t, tt.want, result)
		})
	}
}
