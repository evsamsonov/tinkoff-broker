package tnkbroker

type Money interface {
	GetUnits() int64
	GetNano() int32
}

type MoneyValue struct {
	Money
}

func NewMoneyValue(v Money) *MoneyValue {
	return &MoneyValue{
		Money: v,
	}
}

func (v *MoneyValue) ToFloat() float64 {
	return float64(v.GetUnits()) + float64(v.GetNano())/10e8
}

func NewZeroMoneyValue() *MoneyValue {
	return &MoneyValue{&ZeroMoney{}}
}

type ZeroMoney struct{}

func (*ZeroMoney) GetUnits() int64 { return 0 }

func (*ZeroMoney) GetNano() int32 { return 0 }
