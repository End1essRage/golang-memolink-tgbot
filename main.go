package main

import (
	"flag"
	"log"
	tgClient "memolink-bot/clients/telegram"
	event_consumer "memolink-bot/consumer/event-consumer"
	telegram "memolink-bot/events/telegram"
	"memolink-bot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped")
	}
}

func mustToken() string {
	token := flag.String("bot-token", "", "token for access to telegram bot")

	flag.Parse()

	if *token == "" {
		panic("token is empty")
	}

	return *token
}
