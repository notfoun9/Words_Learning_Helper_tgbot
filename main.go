package main

import (
	"flag"
	"log"
	"telegram-bot/clients/telegram"
	event_consumer "telegram-bot/consumer/event-consumer"
	telegram_events "telegram-bot/events/telegramEvents"
	"telegram-bot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {
	tgClient := telegram.New(tgBotHost, mustToken())

	eventsProcessor := telegram_events.New(
		tgClient,
		files.New(storagePath),
	)

	log.Printf("Service launched")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal()
	}

}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",                 // name
		"",                             // value
		"token to acsess telegram bot", //usage
	)
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
