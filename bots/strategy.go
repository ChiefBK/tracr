package bots

import (
	log "github.com/inconshreveable/log15"
	"tracr/bots/actions"
	"tracr-cache"
	"tracr-store"
)

type Strategy struct {
	Trees    []*DecisionTree
	Position string
}

func NewStategy(position string) *Strategy {
	var trees []*DecisionTree
	return &Strategy{trees, position}
}

func (self *Strategy) AddTree(tree *DecisionTree) {
	self.Trees = append(self.Trees, tree)
}

func (self *Strategy) run(botActionChan chan<- *actions.ActionQueue, bot *Bot, cacheClient *tracr_cache.CacheClient, storeClient tracr_store.Store) {
	log.Debug("running strategy", "module", "command", "strategyPosition", self.Position)
	botActionQueue := actions.NewActionQueue() // the queue that will be sent back to the bot

	for _, tree := range self.Trees {
		treeActionChan := make(chan *actions.ActionQueue)
		go tree.run(treeActionChan, bot, cacheClient, storeClient)

		treeActionQueue := <-treeActionChan // gets actions from tree

		for _, action := range treeActionQueue.Queue { // add actions from tree action queue to bot action queue
			botActionQueue.Push(action)
		}
	}

	botActionChan <- botActionQueue

}
