package bots

import (
	log "github.com/inconshreveable/log15"
	"tracr/bots/actions"
	"time"
	"strings"
	"tracr-cache"
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
	Key        string // must be unique amongst bots
	Exchange   string
	Pair       string
	Position   string
	Strategies map[string]*Strategy
	Data       BotData
	IsRunning  bool
	Interval   time.Duration // at the top of which interval to run bot
	Command    chan string
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

var runningBots map[string]*Bot // map botKey to bot

func NewBot(key, exchange, pair string) (bot *Bot) {
	bot = new(Bot)
	bot.Strategies = make(map[string]*Strategy)
	bot.Exchange = exchange
	bot.Pair = pair
	bot.Key = key
	bot.Position = CLOSED_POSITION
	bot.Data = buildDefaultBotData()
	bot.IsRunning = false
	bot.Interval = 5 * time.Minute // default duration of 5 minutes
	bot.Command = make(chan string)

	return
}

func addBot(bot *Bot) {
	// TODO - respond with error if bot already exists. Require --overwrite flag to override

	log.Info("Adding Bot", "module", "command", "botKey", bot.Key)
	botEncoding := toGOB64(bot)
	log.Debug("bot encoding", "module", "command", "encoding", botEncoding)
	tracr_cache.PutBotEncoding(bot.Key, botEncoding)
}

/**
	Returns nil if Bot doesn't exist
 */
func fetchBot(botName string) *Bot {
	log.Info("Fetching Bot", "module", "command", "botKey", botName)
	botEncoding, err := tracr_cache.GetBotEncoding(botName)

	if err != nil { // error getting bot from cache. maybe doesn't exist?
		log.Error("error getting bot from cache", "module", "command", "botKey", botName, "error", err)
		return nil
	}

	bot := fromGOB64(botEncoding)
	return bot
}

func (self *Bot) start(interval time.Duration, startImmediately bool) {
	log.Info("Starting bot", "botKey", self.Key, "module", "command", "interval", interval, "startImmediately", startImmediately)

	go func() {
		// if interval specified in cmd args
		if interval.Minutes() != 0 {
			self.Interval = interval
		}

		for {
			if !startImmediately {
				<-time.After(1 * time.Minute) // wait one minute
			} else {
				log.Debug("starting immediately", "module", "command", "botKey", self.Key)
				startImmediately = false
				go self.run()
				continue
			}

			now := time.Now()
			nowUnix := now.Unix() / 60
			botIntervalUnix := int64(self.Interval.Minutes())

			// TODO - keep record of the last running time - if time.Now minus interval is greater than last running time -> start
			// so if the exact minute the bot runs in is missed it'll run anyways
			if nowUnix%botIntervalUnix == 0 { // if it's time for bot to run
				log.Debug("starting at interval", "module", "command", "botKey", self.Key, "interval", self.Interval, "now", now)
				go self.run()
			}

			select {
			case command := <-self.Command:
				if strings.ToLower(command) == "stop" { // if stop command issued
					log.Info("stop command received", "module", "command")
					return
				}
			case <-time.After(time.Nanosecond * 1): // timeout if no command has been issued
			}
		}

	}()

}

func (self *Bot) stop() {
	log.Info("stop command issued", "module", "command")
	self.Command <- "stop"
}

func (self *Bot) run() {
	ready := self.preChecks()

	if !ready {
		// print warning bot could not run because of precheck
		log.Error("failed pre-check", "module", "command", "botKey", self.Key)
		return
	}

	self.IsRunning = true

	log.Debug("simulate bot running for a few seconds", "module", "command", "botKey", self.Key)
	<-time.After(time.Second * 10)
	log.Debug("done", "module", "command", "botKey", self.Key)

	//var signalActionChan = make(chan *actions.ActionQueue)
	//self.runStrategy(signalActionChan)
	//
	//signalActionQueue := <-signalActionChan
	//
	//botActionQueue := actions.NewActionQueue()
	//log.Debug("received actions from strategy", "botKey", self.Key, "module", "command", "actionLen", signalActionQueue.Length())
	//
	//action := signalActionQueue.Dequeue()
	//
	//for action != nil {
	//	log.Debug("processing action from strategy", "botKey", self.Key, "module", "command", "action", action)
	//
	//	//return
	//	if action.Consumer == actions.BOT {
	//		// handle internal action
	//	} else { // if actions.EXECUTOR
	//		action.SetVolume(self.data.volume())
	//		action.SetLeverage(self.data.leverage())
	//		action.SetMargin(self.data.margin())
	//		action.SetOrderType(self.data.orderType())
	//		action.SetPair(self.pair)
	//		action.SetExchange(self.exchange)
	//		action.SetBotKey(self.Key)
	//		botActionQueue.Push(action)
	//	}
	//
	//	action = signalActionQueue.Dequeue()
	//}
	//
	//responseChannel := broker.GetBotResponseChannel(self.Key)
	////send actions to action receiver
	//broker.SendToExecutor(self.Key, *botActionQueue)
	//executorResponse := <-responseChannel
	//
	//log.Debug("Received executor response", "botKey", self.Key, "module", "command", "response", executorResponse)

	self.IsRunning = false
}

func (self *Bot) addStrategy(strategy *Strategy) {
	pos := strategy.Position
	self.Strategies[pos] = strategy
}

func (self *Bot) runStrategy(actionChan chan<- *actions.ActionQueue) {
	go self.Strategies[self.Position].run(actionChan)
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

// Checks if bot has green light to run
// Returns true if bot is ok to run, false otherwise
func (self *Bot) preChecks() bool {
	if self.IsRunning {
		log.Error("Bot already running", "module", "command", "botKey", self.Key)
		return false
	}

	return true
}
