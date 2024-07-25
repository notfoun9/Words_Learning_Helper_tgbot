package main

import (
	"flag"
	"log"
	"telegram-bot/clients/telegram_client"
	event_consumer "telegram-bot/consumer/event-consumer"
	telegram_events "telegram-bot/events/telegramEvents"
	"telegram-bot/storage/sqlite"
)

const (
	tgBotHost        = "api.telegram.org"
	storagePath      = `storage/data/users`
	storageSql       = `storage/data/sqlite`
	updatesBatchSize = 100
)

func main() {
	tgClient := telegram_client.NewClient(tgBotHost, mustToken())

	storage, err := sqlite.NewSQLStorage(storageSql)
	if err != nil {
		log.Println(err.Error())
	}
	storage.Init()

	eventsProcessor := telegram_events.New(
		tgClient,
		storage,
	)

	log.Printf("Service launched")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, updatesBatchSize)

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
