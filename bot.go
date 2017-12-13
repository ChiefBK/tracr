package tracr

import (
	log "github.com/inconshreveable/log15"
	"tracr/actions"
	"tracr-daemon/broker"
)

const (
	CLOSED_POSITION = "CLOSED"
	LONG_POSITION   = "LONG"
	SHORT_POSITION  = "SHORT"
)

type BotData map[string]interface{}

func (self BotData) volume() float64 {
	return self["volume"].(float64)
}

func (self BotData) margin() float64 {
	return self["margin"].(float64)
}

func (self BotData) leverage() int {
	return self["leverage"].(int)
}

func (self BotData) orderType() actions.OrderType {
	return self["orderType"].(actions.OrderType)
}

type Bot struct {
	key        string // must be unique amongst bots
	exchange   string
	pair       string
	position   string
	strategies map[string]*Strategy
	data       BotData
}

//func addBot(botKey, exchange, pair string, data map[string]interface{}, trees ...*DecisionTree) {
//
//	bot1 := NewBot(botKey, exchange, pair)
//
//	var closedPositionTrees []*DecisionTree
//	var longPositionTrees []*DecisionTree
//	var shortPositionTrees []*DecisionTree
//
//	for _, tree := range trees {
//		switch tree.position {
//		case CLOSED_POSITION:
//			closedPositionTrees = append(closedPositionTrees, tree)
//		case LONG_POSITION:
//			longPositionTrees = append(longPositionTrees, tree)
//		case SHORT_POSITION:
//			shortPositionTrees = append(shortPositionTrees, tree)
//		}
//	}
//
//	closedStrat := NewStategy(closedPositionTrees)
//	bot1.addStrategy(CLOSED_POSITION, closedStrat)
//	longStrat := NewStategy(longPositionTrees)
//	bot1.addStrategy(LONG_POSITION, longStrat)
//	shortStrat := NewStategy(shortPositionTrees)
//	bot1.addStrategy(SHORT_POSITION, shortStrat)
//
//	broker.BotResponseChannels[botKey] = make(chan responses.ExecutorResponse) // open channel to receive response from executors module
//	broker.AddActionReceiverChannel(botKey)                                    // open channel to executors to receive requests from bot
//	bots = append(bots, bot1)                                                  // add bot to list of bots in strategy module
//}

func NewBot(key, exchange, pair string) (bot *Bot) {
	bot = new(Bot)
	bot.strategies = make(map[string]*Strategy)
	bot.exchange = exchange
	bot.pair = pair
	bot.key = key
	bot.position = CLOSED_POSITION
	bot.data = buildDefaultBotData()

	return
}

func (self *Bot) start() {
	log.Info("Starting bot", "botKey", self.key, "module", "command")
	var signalActionChan = make(chan *actions.ActionQueue)
	self.runStrategy(signalActionChan)

	signalActionQueue := <-signalActionChan

	botActionQueue := actions.NewActionQueue()
	log.Debug("received actions from strategy", "botKey", self.key, "module", "command", "actionLen", signalActionQueue.Length())

	action := signalActionQueue.Dequeue()

	for action != nil {
		log.Debug("processing action from strategy", "botKey", self.key, "module", "command", "action", action)

		//return
		if action.Consumer == actions.BOT {
			// handle internal action
		} else { // if actions.EXECUTOR
			action.SetVolume(self.data.volume())
			action.SetLeverage(self.data.leverage())
			action.SetMargin(self.data.margin())
			action.SetOrderType(self.data.orderType())
			action.SetPair(self.pair)
			action.SetExchange(self.exchange)
			action.SetBotKey(self.key)
			botActionQueue.Push(action)
		}

		action = signalActionQueue.Dequeue()
	}

	responseChannel := broker.GetBotResponseChannel(self.key)
	//send actions to action receiver
	broker.SendToExecutor(self.key, *botActionQueue)
	executorResponse := <-responseChannel

	log.Debug("Received executor response", "botKey", self.key, "module", "command", "response", executorResponse)
}

func (self *Bot) addStrategy(strategy *Strategy) {
	pos := strategy.position
	self.strategies[pos] = strategy
}

func (self *Bot) runStrategy(actionChan chan<- *actions.ActionQueue) {
	go self.strategies[self.position].run(actionChan)
}

func buildDefaultBotData() (data BotData) {
	data = make(BotData)
	data["volume"] = 1.0
	data["leverage"] = 2
	data["margin"] = 0.5

	var orderType actions.OrderType = actions.MARKET_ORDER
	data["orderType"] = orderType

	return
}