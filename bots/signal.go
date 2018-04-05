package bots

import (
	log "github.com/inconshreveable/log15"
	"tracr/bots/actions"
)

type Signal struct {
	Condition    string
	Children     []*Signal
	ActionIntent *string
	IsRoot       bool
}

func NewSignal(conditionFuncName string, action *actions.ExternalAction, isRoot bool) *Signal {
	//var children []*Signal
	//return &Signal{conditionFuncName, children, action, isRoot}
	return nil
}

func (self *Signal) addChild(signal *Signal) {
	self.Children = append(self.Children, signal)
}

func (self *Signal) run(actionChan chan<- *actions.ExternalAction) {
	log.Debug("running signal", "module", "command", "children", len(self.Children), "isRoot", self.IsRoot)
	//condFunc := conditions.ConditionFunctions[self.Condition]
	//result := condFunc()

	//if result { // if signal is true
	//	if len(self.Children) == 0 && self.ActionIntent != nil { // if leaf node
	//		log.Debug("sending action from signal", "module", "command", "action", self.ActionIntent)
	//		actionChan <- self.ActionIntent // send action to tree
	//	}
	//
	//	for _, child := range self.Children {
	//		child.run(actionChan)
	//	}
	//}
	//
	//if self.IsRoot { // when all children of root have run then close action channel
	//	close(actionChan)
	//}
}
