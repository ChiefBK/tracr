package actions

import (
	"tracr-daemon/util"
)

const (
	OPEN_SHORT_POSITION ActionIntent = "openShortPosition"
	OPEN_LONG_POSITION  ActionIntent = "openLongPosition"
	CLOSE_POSITION      ActionIntent = "closePosition"
)

const (
	MARKET OrderType = "market"
	LIMIT  OrderType = "limit"
)

const (
	INTERNAL ActionConsumer = "internal"
	EXTERNAL ActionConsumer = "external"
)

type ActionIntent string   // OPEN_SHORT_POSITION or OPEN_LONG_POSITION or CLOSE_POSITION or ...
type OrderType string      // MARKET_ORDER or LIMIT_ORDER or ...
type ActionConsumer string // BOT or EXECUTOR or ...

type OrderData struct {
	Id       string
	Price    float64
	Volume   float64
	Leverage float64
	Type     OrderType
}

type Action interface {
	Consumer() ActionConsumer
}

type ExternalAction struct {
	Id       string
	Intent   ActionIntent
	Order    OrderData
}

func (self ExternalAction) Consumer() ActionConsumer {
	return EXTERNAL
}

type InternalAction struct {
	Id         string
	PropChange map[string]interface{}
}

func (self InternalAction) Consumer() ActionConsumer {
	return INTERNAL
}

func newExternalAction(intent ActionIntent, orderData OrderData) *ExternalAction {
	id := util.RandString(20)
	return &ExternalAction{id, intent, orderData}
}

func newInternalAction(change map[string]interface{}) *InternalAction {
	id := util.RandString(20)
	return &InternalAction{id, change}

}
