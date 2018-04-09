package bots

import (
	log "github.com/inconshreveable/log15"
	"tracr/bots/actions"
	"tracr/bots/conditions"
	"tracr-cache"
	"tracr-store"
)

type Signal struct {
	Condition      string
	Children       []*Signal
	ActionIntent   *actions.ActionIntent
	ActionConsumer *actions.ActionConsumer
	IsRoot         bool
}

func NewSignal(conditionFuncName string, action *actions.ExternalAction, isRoot bool) *Signal {
	//var children []*Signal
	//return &Signal{conditionFuncName, children, action, isRoot}
	return nil
}

func (self *Signal) addChild(signal *Signal) {
	self.Children = append(self.Children, signal)
}

func (self *Signal) run(actionChan chan<- actions.Action, bot *Bot, cacheClient *tracr_cache.CacheClient, storeClient tracr_store.Store) {
	log.Debug("running signal", "module", "command", "children", len(self.Children), "isRoot", self.IsRoot)
	condFunc := conditions.ConditionFunctions[self.Condition]
	botProps := bot.Props
	result := condFunc(botProps, cacheClient, storeClient)

	if result { // if signal is true
		if len(self.Children) == 0 && self.ActionIntent != nil && self.ActionConsumer != nil { // if leaf node
			log.Debug("sending action from signal", "module", "command", "action", *self.ActionIntent)

			switch *self.ActionConsumer {
			case actions.EXTERNAL:
				actionChan <- actions.NewExternalAction(*self.ActionIntent)
			case actions.INTERNAL:
				// TODO
			}
		}

		for _, child := range self.Children {
			child.run(actionChan, bot, cacheClient, storeClient)
		}
	}

	if self.IsRoot { // when all children of root have run then close action channel
		close(actionChan)
	}
}
