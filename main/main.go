package main

import (
	"os"
	"fmt"
	"tracr/bots"
)

/*
	tracr start <botName>
	tracr stop <botName>
	tracr init <botName> <botTemplateFilePath>
	tracr monitor <botName> ...
 */

func main() {
	args := os.Args[1:]
	fmt.Println(args)

	if len(args) == 0 {
		// print help
		return
	}

	if len(args) < 2 {
		// print error - must have at least 2
		return
	}

	bots.Init()

	action := args[0]
	botName := args[1]

	switch action {
	case "start":
		start(botName)
	case "stop":
		stop(botName)
	case "init":
		if len(args) < 3 {
			// print error
			return
		}

		templatePath := args[2]

		initialize(botName, templatePath)
	case "monitor":
		monitor(args[2:]...)
	default:
		// error - action not valid
		return
	}
}

func initialize(botName, templatePath string) {

}

func start(botName string) {

}

func stop(botName string) {

}

func monitor(botNames... string) {

}
