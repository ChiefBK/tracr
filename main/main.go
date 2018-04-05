package main

import (
	"os"
	"fmt"
	"tracr/bots"
	"flag"
	"time"
	log "github.com/inconshreveable/log15"
	"strings"
	"tracr-cache"
	"tracr-store"
)

/*

	Program Usage

	tracr create [--overwrite] <botTemplateFilePath>
	tracr start (-s | -i) (<botName>... | --all)
	tracr start (-s -i) (<botName>... | --all)
	tracr (stop|destroy|monitor) (<botName>... | --all)
	tracr list


	Options

	--help -h 										show help
	--all -a										selects all bots
	--start-now -s									starts bot(s) now
	--interval <interval>, -i <interval>			the interval for the bot in minutes
	--overwrite										overwrite bot if already exists


	see http://docopt.org/ for docs on program usage syntax

 */

func main() {
	args := os.Args[1:]
	fmt.Println(args)

	if len(args) == 0 {
		// print help
		return
	}

	if len(args) < 2 {
		log.Error("Must have at least two cmd line arguments")
		return
	}

	action := determineAction(args)
	startNow1 := flag.Bool("s", false, "Start the bot immediately")
	startNow2 := flag.Bool("start-now", false, "Start the bot immediately")
	interval1 := flag.Int("i", 0, "Interval of bot")
	interval2 := flag.Int("interval", 0, "Interval of bot")
	allBots := flag.Bool("all", false, "Selects all bots")
	flag.Parse()

	log.Debug("startNow", "module", "command", "sn1", *startNow1, "sn2", *startNow2)

	var botName string
	var interval time.Duration
	var startNow = false

	if *allBots {
		botName = ""
	} else {
		botName = args[len(args)-1]
	}

	if *interval1 != 0 {
		interval = time.Duration(*interval1) * time.Minute
	} else if *interval2 != 0 {
		interval = time.Duration(*interval2) * time.Minute
	} else {
		interval = 0 * time.Minute
	}

	if *startNow1 || *startNow2 {
		startNow = true
	}

	bots.Init()

	cacheClient, err := tracr_cache.NewCacheClient()

	if err != nil {
		log.Error("Can not continue with nil cache client")
		return
	}

	storeClient, err := tracr_store.NewStore()

	if err != nil {
		log.Error("Can not continue without valid store")
		return
	}

	switch action {
	case "start":
		bots.Start(botName, interval, startNow, cacheClient, storeClient)
	case "stop":
		bots.Stop(botName)
	case "destroy":
		destroy(botName)
	case "create":
		if len(args) != 2 {
			log.Error("Invalid number of arguments for initializing bot", "module", "command", "numOfArgsProvided", len(args))
			return
		}

		templatePath := args[1]

		bots.InitializeBot(templatePath, cacheClient)
	case "monitor":
		monitor(args[2:]...)
	case "list":
		list()
	default:
		log.Error("specified action not valid")
		return
	}
}

func determineAction(args []string) string {
	for _, arg := range args {
		switch strings.ToLower(arg) {
		case "start":
			return "start"
		case "stop":
			return "stop"
		case "create":
			return "create"
		}
	}

	return ""
}

func stop(botName string) {

}

func destroy(botName string) {

}

func monitor(botNames ... string) {

}

func list() {

}
