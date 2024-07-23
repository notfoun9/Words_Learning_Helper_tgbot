package telegram_events

import (
	"errors"
	"log"
	"os"
	"strings"
	"telegram-bot/clients/telegram"
	"telegram-bot/lib/e"
	"telegram-bot/storage"
)

const (
	RndCmd    = "/rnd"
	HelpCmd   = "/help"
	StartHelp = "/start"
	RemoveCmd = "/rmv"
	AllCmd    = "/all"

	ReadyToRemove = "removeWord"
)

func (p *ProcessorTelegram) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new cmd %s from %s", text, username)

	state, err := p.wordStorage.GetState(username)
	if err != nil {
		return err
	}

	if state == "" {
		switch text {
		case RndCmd:
			return p.sendRandomWord(chatID, username)
		case HelpCmd:
			return p.sendHelp(chatID)
		case StartHelp:
			return p.sayHello(chatID)
		case RemoveCmd:
			return p.removeCmd(chatID, username)
		case AllCmd:
			return p.printAll(chatID, username)
		default:
			return p.saveWord(chatID, text, username)
		}
	} else if state == ReadyToRemove {
		return p.removeWord(chatID, username, text)
	} else {
		return p.giveDefinitionWord(chatID, username, text)
	}
}

func (p *ProcessorTelegram) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *ProcessorTelegram) sayHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func (t *ProcessorTelegram) saveWord(chatID int, word string, username string) (err error) {
	defer func() { err = e.Wrap("can't do comand: save page", err) }()
	sendMessage := NewMessageSender(chatID, t.tg)

	wordToAdd := &storage.Word{
		Word:     word,
		UserName: username,
	}

	doesExist, err := t.wordStorage.DoesExistWord(username, word)
	if err != nil {
		return err
	}
	if doesExist {
		return sendMessage(msgAlreadyExists)
	}

	if err := t.wordStorage.SaveWord(wordToAdd); err != nil {
		return err
	}

	return sendMessage(msgGiveDefinition)
}

func (t *ProcessorTelegram) sendRandomWord(chatID int, username string) (err error) {
	defer func() { err = e.Wrap("can't do comand: pick random", err) }()
	sendMessage := NewMessageSender(chatID, t.tg)

	word, err := t.wordStorage.PickRandomWord(username)

	if err != nil {
		if errors.Is(err, storage.ErrNoPagesSaved) {
			return sendMessage(msgNoSavedWords)
		} else {
			return err
		}
	}

	if err := sendMessage(word.Word); err != nil {
		return err
	}

	return t.tg.SendSpoilerMessage(chatID, word.Definition)
}

func (t *ProcessorTelegram) removeCmd(chatID int, username string) error {
	if err := t.tg.SendMessage(chatID, msgWordToDelete); err != nil {
		return err
	}

	return t.wordStorage.SetState(username, ReadyToRemove)
}

func (t *ProcessorTelegram) removeWord(chatID int, username string, word string) error {
	err := t.wordStorage.RemoveWord(username, word)
	if err == os.ErrNotExist {
		return t.tg.SendMessage(chatID, msgNoSuchWord)
	} else if err != nil {
		return err
	}

	return t.tg.SendMessage(chatID, word+msgWordRemoved)
}

func (t *ProcessorTelegram) giveDefinitionWord(chatID int, username string, def string) (err error) {
	err = t.wordStorage.GiveDefinition(username, def)
	if err != nil {
		return err
	}
	return t.tg.SendMessage(chatID, msgSaved)
}

func (t *ProcessorTelegram) printAll(chatID int, username string) (err error) {
	defer func() { e.Wrap("Unable to print the list ", err) }()

	words, err := t.wordStorage.AllWords(username)
	if err == storage.ErrNoPagesSaved {
		return t.tg.SendMessage(chatID, msgNoSavedWords)
	} else if err != nil {
		return err
	}

	for i := 0; i < len(words); i++ {
		err := t.tg.SendMessage(chatID, words[i].Word)
		if err != nil {
			return err
		}

		err = t.tg.SendSpoilerMessage(chatID, words[i].Definition)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}
