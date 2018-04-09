package bots

import (
	log "github.com/inconshreveable/log15"
	"tracr/bots/conditions"
	"tracr/bots/actions"
	"time"
	"encoding/gob"
	"tracr-cache"
	"tracr-store"
)

func Init() error {
	log.Info("Initializing bots module", "module", "command")

	conditions.Init()

	//actions.Init()

	//gob.Register(Bot{})
	//gob.Register(Strategy{})
	//gob.Register(BotData{})
	//gob.Register(DecisionTree{})
	//gob.Register(actions.Action{})
	gob.Register(actions.OrderType(0))

	return nil

	//path := filepath.Join("bot_templates", "bot1Template.json")
	//bot, err := readBotFile(path)

	//if err != nil {
	//	log.Error("there was an error reading the bot file", "module", "command", "file", path, "error", err)
	//}

	//log.Debug("Bot created", "module", "command", "bot", bot)
	//log.Debug("strategies created", "module", "command", "strats", bot.strategies)
	//log.Debug("trees created", "module", "command", "trees", bot.strategies["CLOSED"].decisionTrees)
	//log.Debug("root created", "module", "command", "signals", bot.strategies["CLOSED"].decisionTrees[0].root)
	//log.Debug("child created", "module", "command", "signals", bot.strategies["CLOSED"].decisionTrees[0].root.children[0])

	//rootSignal := NewSignal(conditions.TrueFunction, nil, true)
	//leafSignal := NewSignal(conditions.TrueFunction, actions.ShortPositionAction(), false)
	//tree := BuildDecisionChain(CLOSED_POSITION, rootSignal, leafSignal)
	//addBot("bot1", "poloniex", "USDT_BTC", nil, tree)

}

func Start(botName string, interval time.Duration, startImmediately bool, cacheClient *tracr_cache.CacheClient, storeClient tracr_store.Store) {
	//log.Info("Starting command module", "module", "command")
	//for _, bot := range bots {
	//	go bot.start()
	//}
	bot := fetchBot(botName, cacheClient)

	if bot == nil {
		// error - doesn't contain bot
		log.Error("bot key specified doesn't not exist", "module", "command")
		return
	}

	bot.start(interval, startImmediately, cacheClient, storeClient)
}

func Stop(botName string) {
	panic("stop not implemented")
	//log.Info("bot has stopped", "module", "command", "botKey", botName)
}

func InitializeBot(templatePath string, cacheClient *tracr_cache.CacheClient) {
	bot, err := readBotFile(templatePath)

	log.Debug("real bot", "module", "command", "num of strategies", bot.IsRunning)
	//log.Debug("real bot", "module", "command", "'root signal condition'", bot.Strategies[0].Trees[0].Root.Condition)

	if err != nil {
		log.Error("There was an error initializing the bot", "module", "command", "templatePath", templatePath)
		return
	}

	saveBot(bot, cacheClient)

	// test code
	//botName := bot.Name
	//botCopy := fetchBot(botName, cacheClient)
	//
	//if botCopy == nil {
	//	log.Debug("bot copy is nil", "module", "command")
	//	return
	//}
	//
	//log.Debug("bot copy", "module", "command", "key", botCopy.Name)
	//log.Debug("bot copy", "module", "command", "exchange", botCopy.Exchange)
	//log.Debug("bot copy", "module", "command", "data", botCopy.Data)
	//log.Debug("bot copy", "module", "command", "strategy length", len(botCopy.Strategies))
	//log.Debug("bot copy", "module", "command", "tree length", len(botCopy.Strategies["closed"].Trees))
	//log.Debug("bot copy", "module", "command", "root signal condition", botCopy.Strategies["closed"].Trees[0].Root.Condition)
	//log.Debug("bot copy", "module", "command", "root signal children", len(botCopy.Strategies["closed"].Trees[0].Root.Children))
	//log.Debug("bot copy", "module", "command", "root signal child's children", len(botCopy.Strategies["closed"].Trees[0].Root.Children[0].Children))
	//log.Debug("bot copy", "module", "command", "root signal child's isRoot", botCopy.Strategies["closed"].Trees[0].Root.Children[0].IsRoot)
	//log.Debug("bot copy", "module", "command", "root signal child's action", botCopy.Strategies["closed"].Trees[0].Root.Children[0].ActionIntent)
	//log.Debug("bot copy", "module", "command", "root signal child's condition", botCopy.Strategies["closed"].Trees[0].Root.Children[0].Condition)
}
