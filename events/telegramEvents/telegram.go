package telegram_events

import (
	"errors"
	"telegram-bot/clients/telegram"
	"telegram-bot/events"
	"telegram-bot/lib/e"
	"telegram-bot/storage"
)

type ProcessorTelegram struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	UserName string
}

func New(client *telegram.Client, stor storage.Storage) *ProcessorTelegram {
	return &ProcessorTelegram{
		tg:      client,
		storage: stor,
	}
}

func (t *ProcessorTelegram) Fetch(limit int) ([]events.Event, error) {
	updates, err := t.tg.Update(t.offset, limit)
	if err != nil {
		return nil, e.Wrap("Fetching error", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))
	for _, upd := range updates {
		res = append(res, event(upd))
	}

	t.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (t *ProcessorTelegram) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return t.processMessage(event)
	default:
		return e.Wrap("can't process message", errors.New("unknown event type"))
	}
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			UserName: upd.Message.From.Username,
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

func (t *ProcessorTelegram) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("impossible to process message", err)
	}

	if err := t.doCmd(event.Text, meta.ChatID, meta.UserName); err != nil {
		return e.Wrap("cant process message", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("impossible to get meta", errors.New("unknown meta type"))
	}
	return res, nil
}
