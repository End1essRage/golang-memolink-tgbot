package telegram

import (
	"errors"
	"fmt"
	"memolink-bot/clients/telegram"
	"memolink-bot/events"
	"memolink-bot/storage"
)

type EventProcessor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatId   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *EventProcessor {
	return &EventProcessor{
		tg:      client,
		storage: storage,
	}
}

func (ep *EventProcessor) Fetch(limit int) ([]events.Event, error) {
	updates, err := ep.tg.Updates(ep.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("cant fetch updates: %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	ep.offset = updates[len(updates)-1].Id + 1

	return res, nil
}

func (ep *EventProcessor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return ep.processMessage(event)
	default:
		return ErrUnknownEventType
	}
}

func (ep *EventProcessor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("cant process message: %w", err)
	}

	if err := ep.doCmd(event.Text, meta.ChatId, meta.Username); err != nil {
		return fmt.Errorf("cant process message: %w", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, ErrUnknownMetaType
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {

	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatId:   upd.Message.Chat.Id,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}
