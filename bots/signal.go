package bots

import (
	log "github.com/inconshreveable/log15"
	"tracr/bots/actions"
	"tracr/bots/conditions"
)

type Signal struct {
	ConditionFunctionName string
	Children              []*Signal
	Action                *actions.Action
	IsRoot                bool
}

func NewSignal(conditionFuncName string, action *actions.Action, isRoot bool) *Signal {
	var children []*Signal
	return &Signal{conditionFuncName, children, action, isRoot}
}

func (self *Signal) addChild(signal *Signal) {
	self.Children = append(self.Children, signal)
}

func (self *Signal) run(actionChan chan<- *actions.Action) {
	log.Debug("running signal", "module", "command", "children", len(self.Children), "isRoot", self.IsRoot)
	condFunc := conditions.ConditionFunctions[self.ConditionFunctionName]
	result := condFunc()

	if result { // if signal is true
		if len(self.Children) == 0 && self.Action != nil { // if leaf node
			log.Debug("sending action from signal", "module", "command", "action", self.Action)
			actionChan <- self.Action // send action to tree
		}

		for _, child := range self.Children {
			child.run(actionChan)
		}
	}

	if self.IsRoot { // when all children of root have run then close action channel
		close(actionChan)
	}
}
