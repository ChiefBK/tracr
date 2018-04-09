package actions

import (
	"tracr-daemon/util"
	"github.com/inconshreveable/log15"
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
	SetOrderData(data OrderData)
}

type ExternalAction struct {
	Id       string
	Intent   ActionIntent
	Order    OrderData
}

func (self *ExternalAction) SetOrderData(data OrderData) {
	self.Order = data
}

func (self ExternalAction) Consumer() ActionConsumer {
	return EXTERNAL
}

type InternalAction struct {
	Id         string
	PropChange map[string]interface{}
}

func (self *InternalAction) SetOrderData(data OrderData) {
	log15.Error("can not set order data for internal action")
	return
}

func (self InternalAction) Consumer() ActionConsumer {
	return INTERNAL
}

func NewExternalAction(intent ActionIntent) *ExternalAction {
	id := util.RandString(20)
	return &ExternalAction{Id: id, Intent: intent}
}

func NewInternalAction(change map[string]interface{}) *InternalAction {
	id := util.RandString(20)
	return &InternalAction{id, change}

}
