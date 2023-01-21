package tnkposition

import (
	"fmt"
	"testing"
	"time"

	"github.com/evsamsonov/trengin"
	"github.com/stretchr/testify/assert"
)

func TestStorage_Run(t *testing.T) {
	// todo сделать позицию закрытой
	position := &trengin.Position{}
	assert.NoError(t, position.Close(time.Now(), 123))

	fmt.Println(position)

	//storage := &Storage{
	//	list: map[trengin.PositionID]*Position{
	//		trengin.PositionID(uuid.New()): {
	//			position: &trengin.Position{
	//				ID:         trengin.PositionID{},
	//				FIGI:       "",
	//				Type:       0,
	//				Quantity:   0,
	//				OpenTime:   time.Time{},
	//				OpenPrice:  0,
	//				CloseTime:  time.Time{},
	//				ClosePrice: 0,
	//				StopLoss:   0,
	//				TakeProfit: 0,
	//				Commission: 0,
	//			},
	//			closed:       nil,
	//			stopLossID:   "",
	//			takeProfitID: "",
	//			orderTrades:  nil,
	//			instrument:   nil,
	//		},
	//	},
	//	deleteTimeout: 1 * time.Second,
	//	closedPositionLifetime:      1 * time.Second,
	//}
	//
	//err := storage.Run(context.Background())
	//assert.NoError(t, err)
}
