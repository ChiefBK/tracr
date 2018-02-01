package bots

import (
	log "github.com/inconshreveable/log15"
	"tracr/bots/conditions"
	"tracr/bots/actions"
	"time"
	"encoding/gob"
	"tracr-cache"
)

func Init() {
	log.Info("Initializing bots module", "module", "command")

	runningBots = make(map[string]*Bot)

	// initialize condition function map
	conditions.ConditionFunctions["TrueFunction"] = conditions.TrueFunction

	// initialize action function map
	actions.ActionFunctions["ShortPositionAction"] = actions.ShortPositionAction

	//gob.Register(Bot{})
	//gob.Register(Strategy{})
	//gob.Register(BotData{})
	//gob.Register(DecisionTree{})
	//gob.Register(actions.Action{})
	gob.Register(actions.OrderType(0))

	tracr_cache.Init()

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

func Start(botName string, interval time.Duration, startImmediately bool) {
	//log.Info("Starting command module", "module", "command")
	//for _, bot := range bots {
	//	go bot.start()
	//}
	bot := fetchBot(botName)

	if bot == nil {
		// error - doesn't contain bot
		log.Error("bot key specified doesn't not exist", "module", "command")
		return
	}

	runningBots[botName] = bot
	bot.start(interval, startImmediately)
	log.Info("bot has started", "module", "command", "botKey", botName)

	log.Info("waiting a few minutes...")
	<-time.After(10 * time.Minute)
}

func Stop(botName string) {
	bot := runningBots[botName]
	bot.stop()
	delete(runningBots, botName)
	log.Info("bot has stopped", "module", "command", "botKey", botName)
}

func InitializeBot(templatePath string) {
	bot, error := readBotFile(templatePath)

	log.Debug("real bot", "module", "command", "'root signal condition'", bot.Strategies["closed"].DecisionTrees[0].Root.ConditionFunctionName)

	if error != nil {
		log.Error("There was an error initializing the bot", "module", "command", "templatePath", templatePath)
		return
	}

	addBot(bot)

	// test code
	botName := bot.Key
	botCopy := fetchBot(botName)

	if botCopy == nil {
		log.Debug("bot copy is nil", "module", "command")
		return
	}

	log.Debug("bot copy", "module", "command", "key", botCopy.Key)
	log.Debug("bot copy", "module", "command", "position", botCopy.Position)
	log.Debug("bot copy", "module", "command", "exchange", botCopy.Exchange)
	log.Debug("bot copy", "module", "command", "data", botCopy.Data)
	log.Debug("bot copy", "module", "command", "strategy length", len(botCopy.Strategies))
	log.Debug("bot copy", "module", "command", "tree length", len(botCopy.Strategies["closed"].DecisionTrees))
	log.Debug("bot copy", "module", "command", "root signal condition", botCopy.Strategies["closed"].DecisionTrees[0].Root.ConditionFunctionName)
	log.Debug("bot copy", "module", "command", "root signal children", len(botCopy.Strategies["closed"].DecisionTrees[0].Root.Children))
	log.Debug("bot copy", "module", "command", "root signal child's children", len(botCopy.Strategies["closed"].DecisionTrees[0].Root.Children[0].Children))
	log.Debug("bot copy", "module", "command", "root signal child's isRoot", botCopy.Strategies["closed"].DecisionTrees[0].Root.Children[0].IsRoot)
	log.Debug("bot copy", "module", "command", "root signal child's action", botCopy.Strategies["closed"].DecisionTrees[0].Root.Children[0].Action)
	log.Debug("bot copy", "module", "command", "root signal child's condition", botCopy.Strategies["closed"].DecisionTrees[0].Root.Children[0].ConditionFunctionName)
}
