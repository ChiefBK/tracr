package bots

import (
	log "github.com/inconshreveable/log15"
	"time"
	"tracr-cache"
	"tracr/bots/properties"
	"tracr/bots/actions"
	"tracr-store"
)

const (
	CLOSED_POSITION = "closed"
	LONG_POSITION   = "long"
	SHORT_POSITION  = "short"
)

type Bot struct {
	Name             string // must be unique amongst bots
	Pair             string
	Exchange         string
	Interval         int // at which interval to run bot (in seconds)
	Strategies       []*Strategy
	Props            properties.Props   // persisted bot properties to be used over many runs
	OrderData        *actions.OrderData // the argument data used to place orders
	Command          chan string
	Position         *string
	IsRunning        *bool
	SuccessfulRuns   *int
	UnsuccessfulRuns *int
	lastTimeRun      time.Time
	//conditionFunctions map[string]conditions.ConditionFunction
}

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

func (self *Bot) start(interval time.Duration, startImmediately bool, cacheClient *tracr_cache.CacheClient, storeClient tracr_store.Store) {
	log.Info("Starting bot", "botKey", self.Name, "module", "command", "interval", interval, "startImmediately", startImmediately)

	self.initializeRuntimeProperties()

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
			go self.runStrategy(cacheClient, storeClient)
			continue
		}

		now := time.Now()
		nowUnix := now.Unix() / 60
		botIntervalUnix := int64(self.Interval)

		// TODO - keep record of the last running time - if time.Now minus interval is greater than last running time -> start
		// so if the exact minute the bot runs in is missed it'll run anyways
		if nowUnix%botIntervalUnix == 0 { // if it's time for bot to run
			log.Debug("starting at interval", "module", "command", "botKey", self.Name, "interval", self.Interval, "now", now)
			self.runStrategy(cacheClient, storeClient)
		}
	}
}

func (self *Bot) stop() {
	log.Info("stop command issued", "module", "command")
	self.Command <- "stop"
}

func (self *Bot) runStrategy(cacheClient *tracr_cache.CacheClient, storeClient tracr_store.Store) {
	ready := self.preChecks()

	if !ready {
		// print warning bot could not run because of precheck
		log.Error("failed pre-check", "module", "command", "botKey", self.Name)
		return
	}

	*self.IsRunning = true
	self.lastTimeRun = time.Now()

	//log.Debug("simulate bot running for a few seconds", "module", "command", "botKey", self.Name)
	//<-time.After(time.Second * 10)
	//log.Debug("done", "module", "command", "botKey", self.Name)

	var actionQueueChan = make(chan *actions.ActionQueue)

	log.Debug("strategies", "num", len(self.Strategies), "botPos", *self.Position)

	stratsAttempted := 0
	for _, strategy := range self.Strategies {
		log.Debug("strategy", "pos", strategy.Position)
		if strategy.Position == *self.Position {
			go strategy.run(actionQueueChan, self, cacheClient, storeClient)
			break
		} else {
			stratsAttempted++
		}
	}

	// close channel if no strategies were run
	if stratsAttempted >= len(self.Strategies) {
		close(actionQueueChan)
	}

	actionQueue := <-actionQueueChan

	log.Debug("received actions from strategy", "botKey", self.Name, "module", "command", "actionLen", actionQueue.Length())

	for _, action := range actionQueue.Queue {
		log.Debug("processing action", "action", action)

		switch action.Consumer() {
		case actions.INTERNAL:
			// TODO
		case actions.EXTERNAL:
			action.SetOrderData(*self.OrderData)
		}

	}

	log.Debug("processed all actions")

	//for action != nil {
	//	log.Debug("processing action", "botKey", self.Name, "module", "command", "action", action)
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
	//	}
	//
	//	action = actionQueue.Dequeue()
	//}

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
		self.OrderData = &actions.OrderData{
			Type:     actions.MARKET,
			Volume:   0.1,
			Leverage: 2,
			Price:    -1,
			Id:       "",
		}
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
