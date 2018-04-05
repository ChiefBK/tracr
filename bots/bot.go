package bots

import (
	log "github.com/inconshreveable/log15"
	"time"
	"strings"
	"tracr-cache"
	"tracr/bots/properties"
)

const (
	CLOSED_POSITION = "CLOSED"
	LONG_POSITION   = "LONG"
	SHORT_POSITION  = "SHORT"
)

type Bot struct {
	Name             string // must be unique amongst bots
	Pair             string
	Exchange         string
	Interval         int // at which interval to run bot (in seconds)
	Strategies       []*Strategy
	Props            properties.Props // persisted bot properties to be used over many runs
	OrderData        properties.Props // the argument data used to place orders
	Command          chan string
	Position         *string
	IsRunning        *bool
	SuccessfulRuns   *int
	UnsuccessfulRuns *int
	lastTimeRun      time.Time
	//conditionFunctions map[string]conditions.ConditionFunction
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

//func NewBot(name, exchange, pair string) (bot *Bot) {
//	bot = new(Bot)
//	bot.Strategies = make(map[string]*Strategy)
//	bot.Exchange = exchange
//	bot.Pair = pair
//	bot.Name = name
//	bot.Position = CLOSED_POSITION
//	bot.Properties = make(map[string]interface{})
//	bot.Data = buildDefaultBotData()
//	bot.IsRunning = false
//	bot.Interval = 5 * time.Minute // default duration of 5 minutes
//	bot.Command = make(chan string)
//
//	return
//}

func saveBot(bot *Bot, client *tracr_cache.CacheClient) {
	// TODO - respond with error if bot already exists. Require --overwrite flag to override

	log.Info("Adding Bot", "module", "command", "botKey", bot.Name)
	botEncoding := toGOB64(bot)
	log.Debug("bot encoding", "module", "command", "encoding", botEncoding)
	client.PutBotEncoding(bot.Name, botEncoding)
}

/**
	Returns nil if Bot doesn't exist or error
 */
func fetchBot(botName string, client *tracr_cache.CacheClient) *Bot {
	log.Info("Fetching Bot", "module", "command", "botKey", botName)
	botEncoding, err := client.GetBotEncoding(botName)

	if err != nil { // error getting bot from cache. maybe doesn't exist?
		log.Error("error getting bot from cache", "module", "command", "botKey", botName, "error", err)
		return nil
	}

	bot := fromGOB64(botEncoding)
	return bot
}

func (self *Bot) start(interval time.Duration, startImmediately bool) {
	log.Info("Starting bot", "botKey", self.Name, "module", "command", "interval", interval, "startImmediately", startImmediately)

	self.initializeRuntimeProperties()

	go func() {
		// if interval specified in cmd args
		if interval.Minutes() != 0 {
			self.Interval = int(interval.Minutes())
		}

		for {
			if !startImmediately {
				<-time.After(1 * time.Minute) // wait one minute
			} else {
				log.Debug("starting immediately", "module", "command", "botKey", self.Name)
				startImmediately = false
				go self.runStrategy()
				continue
			}

			now := time.Now()
			nowUnix := now.Unix() / 60
			botIntervalUnix := int64(self.Interval)

			// TODO - keep record of the last running time - if time.Now minus interval is greater than last running time -> start
			// so if the exact minute the bot runs in is missed it'll run anyways
			if nowUnix%botIntervalUnix == 0 { // if it's time for bot to run
				log.Debug("starting at interval", "module", "command", "botKey", self.Name, "interval", self.Interval, "now", now)
				go self.runStrategy()
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

func (self *Bot) runStrategy() {
	ready := self.preChecks()

	if !ready {
		// print warning bot could not run because of precheck
		log.Error("failed pre-check", "module", "command", "botKey", self.Name)
		return
	}

	*self.IsRunning = true
	self.lastTimeRun = time.Now()

	log.Debug("simulate bot running for a few seconds", "module", "command", "botKey", self.Name)
	<-time.After(time.Second * 10)
	log.Debug("done", "module", "command", "botKey", self.Name)

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

	*self.IsRunning = false
}

//func (self *Bot) addStrategy(strategy *Strategy) {
//	pos := strategy.Position
//	self.Strategies[pos] = strategy
//}
//
//func (self *Bot) runStrategy(actionChan chan<- *actions.ActionQueue) {
//	go self.Strategies[self.Position].run(actionChan)
//}

func (self *Bot) initializeRuntimeProperties() {
	if self.Props == nil {
		self.Props = make(properties.Props)
	}

	if self.OrderData == nil {
		self.OrderData = make(properties.Props)
		self.OrderData["volume"] = 0.1
		self.OrderData["leverage"] = 2
		self.OrderData["type"] = "market"
	}

	if self.Position == nil {
		self.Position = new(string)
		*self.Position = CLOSED_POSITION
	}

	if self.IsRunning == nil {
		self.IsRunning = new(bool)
		*self.IsRunning = false
	}

	if self.SuccessfulRuns == nil {
		self.SuccessfulRuns = new(int)
		*self.SuccessfulRuns = 0
	}

	if self.UnsuccessfulRuns == nil {
		self.UnsuccessfulRuns = new(int)
		*self.UnsuccessfulRuns = 0
	}
}

// Checks if bot has green light to run
// Returns true if bot is ok to run, false otherwise
func (self *Bot) preChecks() bool {
	if *self.IsRunning {
		log.Error("Bot already running", "module", "command", "botKey", self.Name)
		return false
	}

	return true
}
