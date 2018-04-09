package bots

import (
	log "github.com/inconshreveable/log15"
	"tracr/bots/actions"
	"tracr-cache"
	"tracr-store"
)

type DecisionTree struct {
	Name string
	Root *Signal
}

func newDecisionTree(name string, rootSignal *Signal) *DecisionTree {
	return &DecisionTree{name, rootSignal}
}

func (self *DecisionTree) run(actionQueueChan chan<- *actions.ActionQueue, bot *Bot, cacheClient *tracr_cache.CacheClient, storeClient tracr_store.Store) {
	log.Debug("running tree", "module", "command")

	signalActionChan := make(chan actions.Action)
	actionQueue := actions.NewActionQueue()

	go self.Root.run(signalActionChan, bot, cacheClient, storeClient) // runs root signal of tree

	for action := range signalActionChan { // reads actions from signals thru channel
		actionQueue.Push(action)
	}

	actionQueueChan <- actionQueue // Sends queue of actions to Strategy
}

//func BuildDecisionChain(position string, signals ...*Signal) *DecisionTree {
//	var rootSignal *Signal
//	var refSignal *Signal
//
//	for _, signal := range signals {
//		if rootSignal == nil { // if root signal
//			rootSignal = signal
//			refSignal = signal
//			continue
//		}
//
//		refSignal.addChild(signal)
//		refSignal = signal
//	}
//
//	return newDecisionTree(rootSignal)
//}
