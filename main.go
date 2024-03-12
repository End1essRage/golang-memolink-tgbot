package main

import (
	"flag"
	"log/slog"
	"memolink-bot/clients/telegram"
	"os"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	tgClient := telegram.New(tgBotHost, mustToken())

	//fetcher := fetcher.New(tgClient)

	//processor := processor.New(tgClient)

	//consumer.Start(fetcher, processor)
}

func mustToken() string {
	token := flag.String("bot-token", "", "token for access to telegram bot")

	flag.Parse()

	if *token == "" {
		panic("token is empty")
	}

	return *token
}
