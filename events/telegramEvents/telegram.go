package telegram_events

import (
	"errors"
	"telegram-bot/clients/telegram_client"
	"telegram-bot/events"
	"telegram-bot/lib/e"
	"telegram-bot/storage"
)

type ProcessorState int

type ProcessorTelegram struct {
	tgClient    *telegram_client.Client
	offset      int
	wordStorage storage.WordStorage
}

type MetaOfUpdate struct {
	ChatID   int
	UserName string
}

func New(client *telegram_client.Client, wStor storage.WordStorage) *ProcessorTelegram {
	return &ProcessorTelegram{
		tgClient:    client,
		wordStorage: wStor,
	}
}

func (t *ProcessorTelegram) Fetch(limit int) ([]events.Event, error) {
	updatesPack, err := t.tgClient.Update(t.offset, limit)
	if err != nil {
		return nil, e.Wrap("Fetching error", err)
	}

	if len(updatesPack) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updatesPack))
	for _, upd := range updatesPack {
		res = append(res, event(upd))
	}

	t.offset = updatesPack[len(updatesPack)-1].ID + 1

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

func event(update telegram_client.Update) events.Event {
	updateType := fetchType(update)

	newEvent := events.Event{
		Type: updateType,
		Text: fetchText(update),
	}

	if updateType == events.Message {
		newEvent.Meta = MetaOfUpdate{
			ChatID:   update.ChatID(),
			UserName: update.Username(),
		}
	}
	return newEvent
}

func fetchType(update telegram_client.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(update telegram_client.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}

func (t *ProcessorTelegram) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("impossible to process message", err)
	}

	if len(event.Text) == 0 {
		return nil
	}

	if err := t.doCommand(event.Text, meta.ChatID, meta.UserName); err != nil {
		return e.Wrap("cant process message", err)
	}
	return nil
}

func meta(event events.Event) (MetaOfUpdate, error) {
	res, ok := event.Meta.(MetaOfUpdate)
	if !ok {
		return MetaOfUpdate{}, e.Wrap("impossible to get meta", errors.New("unknown meta type"))
	}
	return res, nil
}
