package bots

import (
	"io/ioutil"
	"encoding/json"
	"tracr/bots/conditions"
	"tracr/bots/actions"
	"bytes"
	"encoding/gob"
	"encoding/base64"
	log "github.com/inconshreveable/log15"
	"errors"
)

func readBotFile(filePath string) (*Bot, error) {
	rawJson, _ := ioutil.ReadFile(filePath)
	var data map[string]interface{}

	err := json.Unmarshal(rawJson, &data)

	if err != nil {
		//log.Error("there was an error un-marshalling stratagies file", "module", "command", "file", filePath)
		return nil, errors.New("there was an error un-marshalling strategies file")
	}

	name := data["name"].(string)
	pair := data["pair"].(string)
	exchange := data["exchange"].(string)
	bot := NewBot(name, exchange, pair)
	strategies := data["strategies"].([]interface{})

	for _, strategy := range strategies {
		position := strategy.(map[string]interface{})["position"].(string)
		trees := strategy.(map[string]interface{})["trees"].([]interface{})
		strat := NewStategy(position)

		for _, tree := range trees {
			rootSignal := tree.(map[string]interface{})["root"]
			signal := buildRoot(rootSignal.(map[string]interface{}))
			decisionTree := newDecisionTree(signal)
			strat.AddTree(decisionTree)
		}

		bot.addStrategy(strat)
	}

	return bot, nil
}

func buildRoot(root map[string]interface{}) *Signal {
	signal := createSignalFromInterface(root)
	children := root["children"].([]interface{})
	for _, child := range children {
		signal.addChild(buildRoot(child.(map[string]interface{})))
	}
	return signal
}

func createSignalFromInterface(raw map[string]interface{}) *Signal {
	isRoot := raw["isRoot"].(bool)
	conditionFunctionName := raw["condition"].(string)
	actionFunctionName, actionNotNull := raw["action"].(string)

	_, condFuncExists := conditions.ConditionFunctions[conditionFunctionName]

	if !condFuncExists {
		log.Error("condition function does not exist", "module", "command", "function name", conditionFunctionName)
	}

	var action *actions.Action
	if actionNotNull {
		actionFunc, ok := actions.ActionFunctions[actionFunctionName]

		if !ok {
			// handle error
		}

		action = actionFunc()
	} else {
		action = nil
	}

	signal := NewSignal(conditionFunctionName, action, isRoot)
	return signal
}

// binary encoder
func toGOB64(bot *Bot) string {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(bot)
	if err != nil {
		log.Error("failed Bot encode", "module", "command", "botKey", bot.Key, "error", err)
	}
	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}

// binary decoder
func fromGOB64(str string) *Bot {
	bot := new(Bot)
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Error("failed Bot decode", "module", "command", "error", err)
	}
	buffer := bytes.Buffer{}
	buffer.Write(by)
	decoder := gob.NewDecoder(&buffer)
	err = decoder.Decode(bot)
	if err != nil {
		log.Error("failed Bot decode(2)", "module", "command", "error", err)
	}
	return bot
}
