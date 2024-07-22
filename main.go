package main

import (
	"flag"
	"log"
	"telegram-bot/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgCliemt := telegram.New(tgBotHost, mustToken())

	// fetcher = fetcher.New(tgCliemt)

	// processot = processor.New(tgCliemt)
	// consumer.Start(fetcher, processor)

}

func mustToken() string {
	token := flag.String(
		"telegram-bot-token",           // name
		"",                             // value
		"token to acsess telegram bot", //usage
	)
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
