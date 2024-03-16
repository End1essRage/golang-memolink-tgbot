package telegram

import (
	"errors"
	"fmt"
	"log"
	"memolink-bot/clients/telegram"
	"memolink-bot/storage"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd" // получение рандомной ссылки
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (ep EventProcessor) doCmd(text string, chatId int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("Got new command '%s' from %s", text, username)

	if isAddCmd(text) {
		return ep.savePage(chatId, text, username)
	}

	switch text {
	case RndCmd:
		return ep.sendRandom(chatId, username)
	case HelpCmd:
		return ep.sendHelp(chatId)
	case StartCmd:
		return ep.sendHello(chatId)
	default:
		return ep.tg.SendMessage(chatId, msgUnknownCommand)
	}
}

func (ep *EventProcessor) savePage(chatId int, pageUrl string, username string) error {
	page := &storage.Page{
		URL:      pageUrl,
		UserName: username,
	}

	send := newMessageSender(chatId, ep.tg)

	isExists, err := ep.storage.IsExists(page)
	if err != nil {
		return fmt.Errorf("cant save page: %w", err)
	}
	if isExists {
		return send(msgAlreadyExists)
	}

	if err := ep.storage.Save(page); err != nil {
		return fmt.Errorf("cant save page: %w", err)
	}

	return send(msgSaved)
}

func (ep *EventProcessor) sendRandom(chatId int, username string) error {
	page, err := ep.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return fmt.Errorf("cant send random")
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return ep.tg.SendMessage(chatId, msgNoSavedPages)
	}

	if err := ep.tg.SendMessage(chatId, page.URL); err != nil {
		return fmt.Errorf("cant send random")
	}

	return ep.storage.Remove(page)
}

func (ep *EventProcessor) sendHelp(chatId int) error {
	return ep.tg.SendMessage(chatId, msgHelp)
}

func (ep *EventProcessor) sendHello(chatId int) error {
	return ep.tg.SendMessage(chatId, msgHello)
}

func newMessageSender(chatId int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatId, msg)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
