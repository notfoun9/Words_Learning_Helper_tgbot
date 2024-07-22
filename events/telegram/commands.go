package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"
	"telegram-bot/clients/telegram"
	"telegram-bot/lib/e"
	"telegram-bot/storage"
)

const (
	RndCmd    = "/end"
	HelpCmd   = "/help"
	StartHelp = "/start"
)

func (p *ProcessorTelegram) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new cmd %s from %s", text, username)

	if isAddCmd(text) {
		p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartHelp:
		return p.sayHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCmd)
	}
}

func (t *ProcessorTelegram) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.Wrap("can't do comand: save page", err) }()
	sendMessage := NewMessageSender(chatID, t.tg)

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	doesExist, err := t.storage.DoesExist(page)
	if err != nil {
		return err
	}

	if doesExist {
		return sendMessage(msgAlreadyExists)
	}

	if err := t.storage.Save(page); err != nil {
		return err
	}

	if err := sendMessage(msgSaved); err != nil {
		return err
	}

	return nil
}

func (t *ProcessorTelegram) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.Wrap("can't do comand: pick random", err) }()
	sendMessage := NewMessageSender(chatID, t.tg)

	page, err := t.storage.PickRandom(username)

	if err != nil {
		if errors.Is(err, storage.ErrNoPagesSaved) {
			sendMessage(msgNoSavedPages)
		} else {
			return err
		}
	}

	if err := sendMessage(page.URL); err != nil {
		return err
	}

	return t.storage.Remove(page)
}

func (p *ProcessorTelegram) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *ProcessorTelegram) sayHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func NewMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
