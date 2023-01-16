package tnkbroker

import (
	"sync"
	"time"

	investapi "github.com/tinkoff/invest-api-go-sdk"

	"github.com/evsamsonov/trengin"
)

type openPosition struct {
	mtx          *sync.Mutex
	position     *trengin.Position
	stopLossID   string
	takeProfitID string
	orderTrades  []*investapi.OrderTrade
	closed       chan trengin.Position
	instrument   *investapi.Instrument
}

func newOpenPosition(
	pos *trengin.Position,
	instrument *investapi.Instrument,
	stopLossID string,
	takeProfitID string,
	closed chan trengin.Position,
) *openPosition {
	return &openPosition{
		position:     pos,
		stopLossID:   stopLossID,
		takeProfitID: takeProfitID,
		closed:       closed,
		instrument:   instrument,
	}
}

func (p *openPosition) SetStopLoss(id string, stopLoss float64) {
	p.stopLossID = id
	p.position.StopLoss = stopLoss
}

func (p *openPosition) SetTakeProfitID(id string, takeProfit float64) {
	p.takeProfitID = id
	p.position.TakeProfit = takeProfit
}

//func (p *brokerPosition) AddCommission(commission float64) {
//	p.mtx.Lock()
//	defer p.mtx.Unlock()
//	if p.position == nil {
//		return
//	}
//	p.position.AddCommission(commission)
//}

func (p *openPosition) AddOrderTrade(orderTrades ...*investapi.OrderTrade) {
	p.orderTrades = append(p.orderTrades, orderTrades...)
}

func (p *openPosition) OrderTrades() []*investapi.OrderTrade {
	result := make([]*investapi.OrderTrade, len(p.orderTrades))
	copy(result, p.orderTrades)
	return p.orderTrades
}

func (p *openPosition) Close(closePrice float64) (trengin.Position, error) {
	if err := p.position.Close(time.Now(), closePrice); err != nil {
		return trengin.Position{}, err
	}

	position := *p.position
	p.closed <- position
	p.position, p.stopLossID, p.takeProfitID = nil, "", ""
	return position, nil
}
